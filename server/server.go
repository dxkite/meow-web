package server

import "C"
import (
	"crypto/tls"
	"crypto/x509"
	"dxkite.cn/gateway/config"
	"dxkite.cn/gateway/route"
	"dxkite.cn/log"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"
)

type Server struct {
	tp  TicketProvider
	cfg *config.Config
	r   *route.Route
	sm  SessionManager
}

func NewServer(cfg *config.Config, r *route.Route) *Server {
	return &Server{
		tp:  NewAESTicketProvider(),
		cfg: cfg,
		r:   r,
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
	}
	sm, err := NewLevelDbSM(s.cfg)
	if err != nil {
		log.Error("session manager create error")
		return err
	}
	s.sm = sm
	return http.Serve(l, s)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Info("request", r.Host, r.URL, r.RemoteAddr, r.Referer())
	var uin uint64
	if m, rr := s.r.Match(r.RequestURI); rr != nil {
		log.Debug("match", m)
		// 读取登录票据
		t := s.ReadTicket(r)
		// 需要登陆才能访问的接口
		if t == nil && rr.Config.Sign {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		if t != nil {
			uin = t.Uin
		}
		b := rr.Group.Get()
		switch b.Type {
		case "http", "https":
			s.procHttp(uin, rr, b, w, r)
		default:
			http.Error(w, "unsupported backend", http.StatusBadRequest)
		}
	} else {
		log.Debug("request", r.Host, r.URL, "not found")
		http.Error(w, "not found", http.StatusNotFound)
	}
}

func (s *Server) ReadTicket(r *http.Request) *Ticket {
	if c, err := r.Cookie(s.cfg.CookieName); err != nil {
		log.Debug("read cookie failed", s.cfg.CookieName)
		return nil
	} else {
		if t, err := s.tp.DecodeTicket(c.Value); err != nil {
			log.Debug("decode ticket error", err)
			return nil
		} else {
			log.Debug("request uin", t.Uin, t.CreateTime)
			expireAt := time.Unix(int64(t.CreateTime), 0).
				Add(time.Second * time.Duration(s.cfg.SessionExpireIn))
			if time.Now().After(expireAt) {
				log.Error("session expired", t.Uin)
				return nil
			}
			if !s.sm.CheckSession(t.Uin) {
				log.Error("session not in db", t.Uin)
				return nil
			}
			return t
		}
	}
}

func (s *Server) procHttp(uin uint64, info *route.RouteInfo, b *route.Backend, w http.ResponseWriter, r *http.Request) {
	url := b.Type + "://" + net.JoinHostPort(b.Host, b.Port) + r.RequestURI
	log.Println("proxy request", url)
	client, err := createClient(s.cfg, b)
	if err != nil {
		http.Error(w, "unavailable", http.StatusServiceUnavailable)
		log.Error(err)
		return
	}
	req, err := http.NewRequest(r.Method, url, r.Body)
	if err != nil {
		http.Error(w, "unavailable", http.StatusServiceUnavailable)
		log.Error(err)
		return
	}
	req.Header = r.Header.Clone()
	req.Header.Set(s.cfg.UinHeaderName, strconv.Itoa(int(uin)))
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "unavailable", http.StatusServiceUnavailable)
		log.Error(err)
		return
	}
	s.procHttpSession(info, w, resp)
	s.procHttpHeader(w, resp)
	w.WriteHeader(resp.StatusCode)
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Error("copy error", err)
	}
}

func createClient(cfg *config.Config, b *route.Backend) (*http.Client, error) {
	c := &http.Client{}
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
	if resp.StatusCode == http.StatusOK {
		respUin, _ := strconv.Atoi(resp.Header.Get(s.cfg.UinHeaderName))
		if info.Config.SignIn && respUin > 0 {
			log.Info("signin", respUin)
			ticket, _ := s.tp.EncodeTicket(uint64(respUin))
			http.SetCookie(w, &http.Cookie{
				Name:    s.cfg.CookieName,
				Value:   ticket,
				Expires: time.Now().Add(time.Second * time.Duration(s.cfg.SessionExpireIn)),
				Secure:  true,
			})
			if err := s.sm.CreateSession(uint64(respUin)); err != nil {
				log.Println("save session error", err)
			}
			return
		}
		if info.Config.SignOut && respUin > 0 {
			log.Info("signout", respUin)
			http.SetCookie(w, &http.Cookie{
				Name:   s.cfg.CookieName,
				Value:  "",
				MaxAge: 0,
				Secure: true,
			})
			if err := s.sm.RemoveSession(uint64(respUin)); err != nil {
				log.Println("remove session error", err)
			}
		}
	}
}

func (s *Server) procHttpHeader(w http.ResponseWriter, resp *http.Response) {
	for k, v := range resp.Header {
		for _, vv := range v {
			w.Header().Set(k, vv)
		}
	}
}
