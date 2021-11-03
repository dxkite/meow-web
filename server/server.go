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
	if s.cfg.EnableTls {
		// 开启TLS
		cert, err := tls.LoadX509KeyPair(s.cfg.TlsCert, s.cfg.TlsKey)
		if err != nil {
			return err
		}
		c := &tls.Config{
			Certificates: []tls.Certificate{cert},
		}

		// 是否校验客户端
		if s.cfg.TlsVerifyClient {
			pool := x509.NewCertPool()
			rootCa, err := ioutil.ReadFile(s.cfg.TlsCa)
			if err != nil {
				return err
			}
			pool.AppendCertsFromPEM(rootCa)
			c.ClientCAs = pool
			c.ClientAuth = tls.RequireAndVerifyClientCert
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

	// 配置CORS请求头
	w.Header().Set("Via", "gw")
	s.writeCorsConfig(w, r)
	if r.Method == http.MethodOptions {
		return
	}

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

	// 发起后端请求
	b := rr.Backend.Get()
	// 剔除多余请求头
	s.normalizeRequest(r)
	var processor proto.Processor
	switch b.BackendType() {
	case "http", "https":
		processor = proto.NewHttpProcessor()
	default:
		if builder, ok := b.(Builder); ok {
			processor = builder.Build()
		} else {
			http.Error(w, "unsupported backend", http.StatusBadRequest)
			return
		}
	}

	resp := NewResponse(s, rr.Config, w)
	err := processor.Do(&proto.BackendContext{
		Cfg:     s.cfg,
		Uin:     uin,
		Ticket:  tks,
		Route:   rr,
		Backend: b,
	}, resp, r)

	// 检查是否写入了请求头
	if !resp.WroteHeader() {
		resp.WriteHttpHeader()
	}

	if err != nil {
		log.Error(err)
		http.Error(w, "backend unavailable", http.StatusInternalServerError)
	}
}

func (s *Server) readTicket(r *http.Request) string {
	c, _ := r.Cookie(s.cfg.Session().GetName())
	var t string
	if c != nil {
		t = c.Value
	}
	if len(t) == 0 {
		t = r.Header.Get("Authorization")
	}
	return ticket.DecodeBase64WebString(t)
}

func (s *Server) readSession(r *http.Request) (string, *ticket.SessionData) {
	tk := s.readTicket(r)

	if len(tk) == 0 {
		return "", nil
	}

	t, err := s.tp.Decode(tk)
	if err != nil {
		log.Debug("decode ticket error", err)
		return tk, nil
	}

	log.Debug("uin", t.Uin, "create_time", t.CreateTime)
	expires := time.Unix(int64(t.CreateTime), 0).Add(s.cfg.Session().GetExpiresIn())
	if time.Now().After(expires) {
		log.Error("session expired", t.Uin)
		return tk, nil
	}

	// 非严格模式直接返回
	if !s.cfg.Session().Strict {
		return tk, t
	}

	// 检查是否在内存中
	sessionInMem := s.sm.CheckSession(t.Uin)
	strictCheck := sessionInMem || s.checkSlo(tk)
	if !strictCheck {
		log.Error("[strict] session not exist", t.Uin)
		return tk, nil
	}

	if !sessionInMem {
		sloExpires := time.Unix(int64(t.CreateTime), 0).Add(s.cfg.Session().GetSloExpiresIn())
		err := s.sm.CreateSession(t.Uin, sloExpires)
		if err != nil {
			log.Error(err)
		} else {
			log.Debug("[strict] sso session", t.Uin, "will expires", sloExpires)
		}
	}
	return tk, t
}

func (s *Server) checkSlo(tk string) bool {
	u := s.cfg.Session().SloUrl
	if len(u) == 0 {
		return false
	}

	// 准备请求
	client := &http.Client{}
	// 设置超时
	client.Timeout = s.cfg.Session().GetSloTimeout()
	// 设置强制检查
	client.Transport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		log.Error("create request", err)
		return false
	}

	// 将检测链接发送
	req.Header.Set("Authorization", tk)
	if resp, err := client.Do(req); err != nil {
		return false
	} else if resp.StatusCode == http.StatusOK {
		return true
	}

	return false
}

func (s *Server) SignIn(w http.ResponseWriter, uin uint64) {
	log.Info("signin", uin)
	t, err := s.tp.Encode(&ticket.SessionData{
		Uin: uin,
	})
	if err != nil {
		log.Error("create session", err)
		return
	}
	expires := time.Now().Add(s.cfg.Session().GetExpiresIn())
	http.SetCookie(w, &http.Cookie{
		Name:     s.cfg.Session().GetName(),
		Value:    ticket.EncodeBase64WebString(t),
		Domain:   s.cfg.Session().Domain,
		Expires:  expires,
		Secure:   s.cfg.Session().Secure,
		Path:     s.cfg.Session().GetPath(),
		HttpOnly: s.cfg.Session().HttpOnly,
	})
	if err := s.sm.CreateSession(uin, expires); err != nil {
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
	if len(origin) == 0 {
		log.Debug("empty origin")
		return
	}

	_, ok := s.corsOrigin[origin]
	if !(ok || s.corsOriginAny) {
		log.Error("not allow origin", origin)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", origin)
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
