package server

import "C"
import (
	"crypto/tls"
	"crypto/x509"
	"dxkite.cn/gateway/config"
	"dxkite.cn/gateway/route"
	"dxkite.cn/gateway/session"
	"dxkite.cn/gateway/session/memsm"
	"dxkite.cn/log"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/textproto"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var DefaultAllowHeader = []string{
	"Content-Type", "Content-Length",
}

type Server struct {
	tp            session.TicketProvider
	cfg           *config.Config
	r             *route.Route
	sm            session.SessionManager
	hf            map[string]bool
	corsOrigin    map[string]bool
	corsOriginAny bool
}

func NewServer(cfg *config.Config, r *route.Route) *Server {
	return &Server{
		tp:  session.NewAESTicketProvider(),
		cfg: cfg,
		r:   r,
	}
}

func (s *Server) ApplyHeaderFilter(allows []string) {
	s.hf = map[string]bool{}
	for _, h := range DefaultAllowHeader {
		h = textproto.CanonicalMIMEHeaderKey(h)
		s.hf[h] = true
		log.Println("allow header", h)
	}
	for _, h := range allows {
		h = textproto.CanonicalMIMEHeaderKey(h)
		s.hf[h] = true
		log.Println("allow header", h)
	}
}

func (s *Server) ApplyCorsConfig(cfg *config.CORSConfig) {
	s.corsOrigin = map[string]bool{}
	for _, h := range cfg.AllowOrigin {
		if h == "*" {
			s.corsOriginAny = true
			break
		}
		s.corsOrigin[h] = true
	}
}

func (s *Server) Serve(l net.Listener) error {
	if s.cfg.EnableVerify {
		pool := x509.NewCertPool()
		rootCa, err := ioutil.ReadFile(s.cfg.CAPath)
		if err != nil {
			return err
		}
		pool.AppendCertsFromPEM(rootCa)
		cert, err := tls.LoadX509KeyPair(s.cfg.ModuleCertPath, s.cfg.ModuleKeyPath)
		if err != nil {
			return err
		}
		c := &tls.Config{
			Certificates: []tls.Certificate{cert},
			ClientAuth:   tls.RequireAndVerifyClientCert,
			ClientCAs:    pool,
		}
		l = tls.NewListener(l, c)
		log.Info("enable ssl config")
	}

	// 当前情况下为了安全session在程序down掉之后就失效了，所以不需要持久化
	sm, err := memsm.NewMemSM(s.cfg)
	if err != nil {
		log.Error("session manager create error")
		return err
	}
	s.sm = sm
	return http.Serve(l, s)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Info("do", r.Method, r.Host, r.URL, r.RemoteAddr, r.Referer())
	var uin uint64

	m, rr := s.r.Match(r.RequestURI)
	if rr == nil {
		log.Debug(r.Method, r.Host, r.URL, "not found")
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	log.Debug("match", m)
	// 读取登录票据
	t := s.ReadTicket(r)
	// 需要登陆才能访问的接口
	if t == nil && rr.Config.Sign {
		s.procUnauthorized(w, r)
		return
	}
	if t != nil {
		uin = t.Uin
	}
	// 配置CORS请求头
	s.writeCorsConfig(w, r)
	w.Header().Set("Via", "gw")
	// 发起后端请求
	b := rr.Group.Get()
	switch b.Type {
	case "http", "https":
		s.procHttp(uin, rr, b, w, r)
	default:
		http.Error(w, "unsupported backend", http.StatusBadRequest)
	}
}

func (s *Server) ReadTicket(r *http.Request) *session.Ticket {
	c, err := r.Cookie(s.cfg.Session().GetName())
	if err != nil || len(c.Value) == 0 {
		return nil
	}
	t, err := s.tp.DecodeTicket(c.Value)
	if err != nil {
		log.Debug("decode ticket error", err)
		return nil
	}
	log.Debug("uin", t.Uin, "create_time", t.CreateTime)
	expiresAt := time.Unix(int64(t.CreateTime), 0).Add(s.cfg.Session().GetExpiresIn())
	if time.Now().After(expiresAt) {
		log.Error("session expired", t.Uin)
		return nil
	}
	if !s.sm.CheckSession(t.Uin) {
		log.Error("session not exist", t.Uin)
		return nil
	}
	return t
}

func (s *Server) procHttp(uin uint64, info *route.RouteInfo, b *route.Backend, w http.ResponseWriter, r *http.Request) {
	u := b.Type + "://" + net.JoinHostPort(b.Host, b.Port) + r.RequestURI
	log.Println("do request", r.Method, u)
	client, err := createClient(s.cfg, b)
	if err != nil {
		http.Error(w, "unavailable", http.StatusServiceUnavailable)
		log.Error(err)
		return
	}
	req, err := http.NewRequest(r.Method, u, r.Body)
	if err != nil {
		http.Error(w, "unavailable", http.StatusServiceUnavailable)
		log.Error(err)
		return
	}
	s.createReqHeader(req, r)
	req.Header.Set(s.cfg.UinHeaderName, strconv.Itoa(int(uin)))
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "unavailable", http.StatusServiceUnavailable)
		log.Error(err)
		return
	}
	s.procHttpSession(info, w, resp)
	s.createRespHeader(w, resp)
	w.WriteHeader(resp.StatusCode)
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Error("copy error", err)
	}
}

func createClient(cfg *config.Config, b *route.Backend) (*http.Client, error) {
	c := &http.Client{}
	c.Timeout = 10 * time.Second
	if b.Type == "https" {
		cert, err := tls.LoadX509KeyPair(cfg.ModuleCertPath, cfg.ModuleKeyPath)
		if err != nil {
			return nil, err
		}
		pool := x509.NewCertPool()
		rootCa, err := ioutil.ReadFile(cfg.CAPath)
		if err != nil {
			return nil, err
		}
		pool.AppendCertsFromPEM(rootCa)
		cfg := &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      pool,
		}
		if len(b.ServerName) != 0 {
			cfg.ServerName = b.ServerName
		}
		c.Transport = &http.Transport{TLSClientConfig: cfg}
	}
	return c, nil
}

func (s *Server) procHttpSession(info *route.RouteInfo, w http.ResponseWriter, resp *http.Response) {
	if resp.StatusCode != http.StatusOK {
		return
	}
	uin, _ := strconv.Atoi(resp.Header.Get(s.cfg.UinHeaderName))
	if info.Config.SignIn && uin > 0 {
		s.SignIn(w, uint64(uin))
		return
	}
	if info.Config.SignOut && uin > 0 {
		s.SignOut(w, uint64(uin))
	}
}

func (s *Server) SignIn(w http.ResponseWriter, uin uint64) {
	log.Info("signin", uin)
	ticket, _ := s.tp.EncodeTicket(uin)
	http.SetCookie(w, &http.Cookie{
		Name:    s.cfg.Session().GetName(),
		Value:   ticket,
		Domain:  s.cfg.Session().Domain,
		Expires: time.Now().Add(s.cfg.Session().GetExpiresIn()),
		Secure:  s.cfg.Session().Secure,
	})
	if err := s.sm.CreateSession(uin); err != nil {
		log.Println("save session error", err)
	}
}

func (s *Server) SignOut(w http.ResponseWriter, uin uint64) {
	log.Info("signout", uin)
	http.SetCookie(w, &http.Cookie{
		Name:   s.cfg.Session().GetName(),
		Value:  "",
		Domain: s.cfg.Session().Domain,
	})
	if err := s.sm.RemoveSession(uin); err != nil {
		log.Println("remove session error", err)
	}
}

func (s *Server) createRespHeader(w http.ResponseWriter, resp *http.Response) {
	for k, v := range resp.Header {
		_, ok := s.hf[textproto.CanonicalMIMEHeaderKey(k)]
		if !ok {
			continue
		}
		for _, vv := range v {
			w.Header().Set(k, vv)
		}
	}
	// 删除UIN返回
	w.Header().Del(s.cfg.UinHeaderName)
}

func (s *Server) createReqHeader(dst, src *http.Request) {
	for k, v := range src.Header {
		_, ok := s.hf[textproto.CanonicalMIMEHeaderKey(k)]
		if !ok {
			continue
		}
		for _, vv := range v {
			dst.Header.Set(k, vv)
		}
	}
	// 设置UA
	dst.Header.Set("User-Agent", "gateway")
	// 设置 ClientIP
	clientIp := src.Header.Get("Client-Ip")
	if len(clientIp) == 0 {
		clientIp, _, _ = net.SplitHostPort(src.RemoteAddr)
	}
	dst.Header.Set("Client-Ip", clientIp)
}

func (s *Server) writeCorsConfig(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	_, ok := s.corsOrigin[origin]
	if len(origin) == 0 {
		origin = "*"
	}
	if ok || s.corsOriginAny {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	if r.Method == http.MethodOptions {
		if len(s.cfg.Cors.AllowMethod) > 0 {
			methods := strings.ToUpper(strings.Join(s.cfg.Cors.AllowMethod, ","))
			w.Header().Set("Access-Control-Allow-Methods", methods)
		}
		if len(s.cfg.Cors.AllowHeader) > 0 {
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(s.cfg.Cors.AllowHeader, ","))
		}
	}
}

func (s *Server) procUnauthorized(w http.ResponseWriter, r *http.Request) {
	if s.cfg.Sign == nil || len(s.cfg.Sign.RedirectUrl) == 0 {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	u := s.cfg.Sign.RedirectUrl
	name := s.cfg.Sign.RedirectName
	if len(name) == 0 {
		name = "redirect_url"
	}
	// 自动判断协议
	uu := fmt.Sprintf("//%s%s", r.Host, r.RequestURI)
	t := url.QueryEscape(uu)
	if strings.Contains(u, "?") {
		u += "&" + name + "=" + t
	} else {
		u += "?" + name + "=" + t
	}
	w.Header().Set("Location", u)
	w.WriteHeader(http.StatusFound)
}
