package proto

import (
	"crypto/tls"
	"crypto/x509"
	"dxkite.cn/gateway/config"
	"dxkite.cn/gateway/route"
	"dxkite.cn/log"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type httpProcessor struct {
	ctx *BackendContext
}

func NewHttpProcessor(ctx *BackendContext) Processor {
	return &httpProcessor{
		ctx: ctx,
	}
}

func buildHttpBackend(u *url.URL, r *http.Request) string {
	baseUrl := u.Scheme + "://" + u.Host
	uri := r.RequestURI
	prefix := u.Query().Get("trim_prefix")
	if len(prefix) > 0 {
		uri = strings.TrimPrefix(uri, prefix)
	}
	if len(u.Path) > 0 {
		uri = u.Path + uri
	}
	return baseUrl + uri
}

func (s *httpProcessor) Do(uin uint64, ticket string) (user uint64, status int, header http.Header, body io.ReadCloser, err error) {
	rt, ok := s.ctx.Backend.(*route.UrlBackend)
	if !ok {
		return 0, 0, nil, nil, fmt.Errorf("unsupported endpoint %T", s.ctx.Backend)
	}
	u := buildHttpBackend(rt.Url, s.ctx.Req)

	log.Println("do req", s.ctx.Req.Method, u)

	req, err := http.NewRequest(s.ctx.Req.Method, u, s.ctx.Req.Body)
	if err != nil {
		return uin, 0, nil, nil, err
	}

	req.Header = s.ctx.Req.Header.Clone()
	req.Header.Set(s.ctx.Cfg.UinHeaderName, strconv.Itoa(int(uin)))
	req.Header.Set("Authorization", ticket)

	client, err := createClient(s.ctx.Cfg, rt.Url)
	if err != nil {
		return uin, 0, nil, nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return uin, 0, nil, nil, err
	}

	user = s.getUin(resp)
	return user, resp.StatusCode, resp.Header, resp.Body, nil
}

func createClient(cfg *config.Config, u *url.URL) (*http.Client, error) {
	c := &http.Client{}
	c.Timeout = 10 * time.Second
	if u.Scheme == "https" {
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
		sn := u.Query().Get("server_name")
		if len(sn) != 0 {
			cfg.ServerName = sn
		}
		c.Transport = &http.Transport{TLSClientConfig: cfg}
	}
	return c, nil
}

func (s *httpProcessor) getUin(resp *http.Response) uint64 {
	if resp.StatusCode != http.StatusOK {
		return 0
	}
	uin, _ := strconv.Atoi(resp.Header.Get(s.ctx.Cfg.UinHeaderName))
	return uint64(uin)
}