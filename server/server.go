package server

import "C"
import (
	"crypto/tls"
	"crypto/x509"
	"dxkite.cn/gateway/config"
	"dxkite.cn/gateway/proto"
	"dxkite.cn/gateway/route"
	"dxkite.cn/gateway/session"
	"dxkite.cn/gateway/session/memsm"
	"dxkite.cn/gateway/ticket"
	"dxkite.cn/log"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/textproto"
	"net/url"
	"strings"
	"time"
)

var DefaultAllowHeader = []string{
	"Content-Type", "Content-Length",
}

type Server struct {
	tp            ticket.TicketEnDecoder
	cfg           *config.Config
	r             *route.Route
	sm            session.SessionManager
	hf            map[string]bool
	corsOrigin    map[string]bool
	corsOriginAny bool
}

func NewServer(cfg *config.Config, r *route.Route) *Server {
	return &Server{
		tp:  ticket.NewAESTicket(cfg.Session().AesTicketKey()),
		cfg: cfg,
		r:   r,
	}
}

func (s *Server) InitTicketMode(mode string) error {
	mode = strings.ToLower(mode)
	switch mode {
	case "rsa":
		// 允许不设置私钥
		tp, err := ticket.NewRsaTicket(s.cfg.Session().RsaKey, s.cfg.Session().RsaCert)
		if err != nil {
			return err
		}
		s.tp = tp
		return nil
	case "aes":
		s.tp = ticket.NewAESTicket(s.cfg.Session().AesTicketKey())
		return nil
	default:
		return fmt.Errorf("unsupported key mode %s", mode)
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

	// 读取会话
	tks, sess := s.readSession(r)
	// 需要登陆才能访问的接口
	if sess == nil && rr.Config.Sign {
		s.procUnauthorized(w, r)
		return
	}

	if sess != nil {
		uin = sess.Uin
	}

	// 配置CORS请求头
	w.Header().Set("Via", "gw")
	s.writeCorsConfig(w, r)
	if r.Method == http.MethodOptions {
		return
	}

	// 发起后端请求
	b := rr.Backend.Get()
	// 剔除多余请求头
	s.normalizeRequest(r)
	var processor proto.Processor
	switch b.BackendType() {
	case "http", "https":
		processor = proto.NewHttpProcessor(&proto.BackendContext{
			Cfg:     s.cfg,
			Req:     r,
			Writer:  w,
			Route:   rr,
			Backend: b,
		})
	default:
		http.Error(w, "unsupported backend", http.StatusBadRequest)
		return
	}

	user, status, header, reader, err := processor.Do(uin, tks)
	if err != nil {
		http.Error(w, "backend unavailable", http.StatusInternalServerError)
	}

	if user > 0 {
		if rr.Config.SignIn {
			s.SignIn(w, user)
		}
		if rr.Config.SignOut {
			s.SignOut(w, user)
		}
	}

	s.createRespHeader(w, header)
	w.WriteHeader(status)

	if _, err := io.Copy(w, reader); err != nil {
		log.Error("copy error", err)
	}
}

func (s *Server) readTicket(r *http.Request) string {
	c, err := r.Cookie(s.cfg.Session().GetName())
	if err != nil {
		return ""
	}
	t := c.Value
	if len(c.Value) == 0 {
		t = r.Header.Get("Authorization")
	}
	return ticket.DecodeBase64WebString(t)
}

func (s *Server) readSession(r *http.Request) (string, *ticket.SessionData) {
	tk := s.readTicket(r)

	t, err := s.tp.Decode(tk)
	if err != nil {
		log.Debug("decode ticket error", err)
		return tk, nil
	}

	log.Debug("uin", t.Uin, "create_time", t.CreateTime)
	expiresAt := time.Unix(int64(t.CreateTime), 0).Add(s.cfg.Session().GetExpiresIn())
	if time.Now().After(expiresAt) {
		log.Error("session expired", t.Uin)
		return tk, nil
	}

	if !s.sm.CheckSession(t.Uin) {
		log.Error("session not exist", t.Uin)
		return tk, nil
	}
	return tk, t
}

func (s *Server) SignIn(w http.ResponseWriter, uin uint64) {
	log.Info("signin", uin)
	t, _ := s.tp.Encode(&ticket.SessionData{
		Uin: uin,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     s.cfg.Session().GetName(),
		Value:    ticket.EncodeBase64WebString(t),
		Domain:   s.cfg.Session().Domain,
		Expires:  time.Now().Add(s.cfg.Session().GetExpiresIn()),
		Secure:   s.cfg.Session().Secure,
		Path:     s.cfg.Session().Path,
		HttpOnly: s.cfg.Session().HttpOnly,
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

func (s *Server) createRespHeader(w http.ResponseWriter, header http.Header) {
	for k, v := range header {
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

func (s *Server) normalizeRequest(req *http.Request) {
	for k := range req.Header {
		_, ok := s.hf[textproto.CanonicalMIMEHeaderKey(k)]
		if !ok {
			req.Header.Del(k)
		}
	}
	// 设置UA
	req.Header.Set("User-Agent", "gateway")
	// 设置 ClientIP
	clientIp := req.Header.Get("Client-Ip")
	if len(clientIp) == 0 {
		clientIp, _, _ = net.SplitHostPort(req.RemoteAddr)
	}
	req.Header.Set("Client-Ip", clientIp)
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

	if s.cfg.Cors.AllowCredentials {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}

	if len(s.cfg.Cors.AllowHeader) > 0 {
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(s.cfg.Cors.AllowHeader, ","))
	}

	if len(s.cfg.Cors.AllowMethod) > 0 {
		methods := strings.ToUpper(strings.Join(s.cfg.Cors.AllowMethod, ","))
		w.Header().Set("Access-Control-Allow-Methods", methods)
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
