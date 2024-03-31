package gateway

import (
	"net/http"
	"net/url"
	"sort"
	"strings"

	"dxkite.cn/log"
	"dxkite.cn/meownest/src/utils"
)

const DefaultHostname = ""

type HttpServer struct {
	authHandler AuthorizationHandler
	authHeader  string
	router      map[string]*Router
	hosts       []string
}

type HttpRouterGroupEntry struct {
	Name          string
	Hostname      []string
	Authorization bool
	Endpoints     []string
	Rewrite       *RewriteConfig
	Paths         []string
}

func NewHttpServer() *HttpServer {
	return &HttpServer{router: map[string]*Router{}}
}

func (r *HttpServer) RegisterAuthorizationHandler(header string, handler AuthorizationHandler) {
	r.authHeader = header
	r.authHandler = handler
}

func (r *HttpServer) RegisterRouterGroup(group *HttpRouterGroupEntry) {
	if len(group.Hostname) == 0 {
		group.Hostname = []string{DefaultHostname}
	}
	for _, host := range group.Hostname {
		if len(group.Paths) == 0 {
			group.Paths = []string{"/"}
		}
		for _, path := range group.Paths {
			router := r.createOrGetRouterByHost(host)
			router.Add(path, &HttpForwardHandler{
				Name:           group.Name,
				Rewrite:        group.Rewrite,
				AuthCheck:      group.Authorization,
				AuthHandler:    r.authHandler,
				IdAssignHeader: r.authHeader,
				Endpoints:      group.Endpoints,
			})
		}
	}
}

func (r *HttpServer) createOrGetRouterByHost(host string) *Router {
	if v, ok := r.router[host]; ok {
		return v
	}

	r.router[host] = NewRouter()
	r.hosts = append(r.hosts, host)
	return r.router[host]
}

func (r *HttpServer) sortHost() {
	sort.Slice(r.hosts, func(a, b int) bool {
		if len(r.hosts[a]) > len(r.hosts[b]) {
			return true
		}
		return r.hosts[a] > r.hosts[b]
	})
}

func (r *HttpServer) Serve(addr string) error {
	uri, err := url.Parse(addr)
	if err != nil {
		return err
	}

	r.sortHost()

	l, err := utils.Listen(uri)
	if err != nil {
		return err
	}

	log.Info("listen", uri.String())
	if err := http.Serve(l, r); err != nil {
		return err
	}
	return nil
}

func (r *HttpServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	host := req.Host
	if v, ok := r.router[host]; ok {
		v.ServeHTTP(w, req)
		return
	}

	for _, v := range r.hosts {
		if strings.HasPrefix(host, v) && r.router[host] != nil {
			r.router[host].ServeHTTP(w, req)
			return
		}
	}

	if r.router[DefaultHostname] != nil {
		r.router[DefaultHostname].ServeHTTP(w, req)
		return
	}

	log.Debug("miss match host name", host)
	http.NotFound(w, req)
}
