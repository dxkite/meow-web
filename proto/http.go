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
}

func NewHttpProcessor() Processor {
	return &httpProcessor{}
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

func (s *httpProcessor) Do(ctx *BackendContext, w http.ResponseWriter, r *http.Request) (err error) {
	rt, ok := ctx.Backend.(*route.UrlBackend)
	if !ok {
		return fmt.Errorf("unsupported endpoint %T", ctx.Backend)
	}
	u := buildHttpBackend(rt.Url, r)

	log.Println("do req", r.Method, u)

	req, err := http.NewRequest(r.Method, u, r.Body)
	if err != nil {
		return err
	}

	req.Header = r.Header.Clone()
	req.Header.Set(ctx.Cfg.UinHeaderName, strconv.Itoa(int(ctx.Uin)))
	req.Header.Set("Authorization", ctx.Ticket)

	client, err := createClient(ctx.Cfg, rt.Url)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	user := s.getUin(ctx.Cfg.UinHeaderName, resp)
	w.Header().Set(ctx.Cfg.UinHeaderName, strconv.Itoa(user))
	w.WriteHeader(resp.StatusCode)

	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Error("copy error", err)
	}
	return nil
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

func (s *httpProcessor) getUin(name string, resp *http.Response) int {
	if resp.StatusCode != http.StatusOK {
		return 0
	}
	uin, _ := strconv.Atoi(resp.Header.Get(name))
	return uin
}
