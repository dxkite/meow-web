package agent

import (
	"net/http"
	"sync"
)

type Server struct {
	hd  http.Handler
	mtx *sync.Mutex
}

func New() *Server {
	return &Server{mtx: &sync.Mutex{}}
}

func (s *Server) Use(h http.Handler) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.hd = h
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.hd.ServeHTTP(w, req)
}

func (s *Server) Run(addr string) {
	http.ListenAndServe(addr, s)
}
