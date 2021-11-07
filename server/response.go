package server

import (
	"dxkite.cn/gateway/route"
	"dxkite.cn/log"
	"net/http"
	"net/textproto"
	"strconv"
)

type Response struct {
	w    http.ResponseWriter
	s    *Server
	rCfg *route.RouteConfig
	req  *http.Request
	// Uin
	uin uint64
	// 头部
	h http.Header
	// 是否处理了头部
	wh bool
	// HTTP
	status int
}

func NewResponse(s *Server, uin uint64, r *route.RouteConfig, req *http.Request, w http.ResponseWriter) *Response {
	return &Response{
		status: http.StatusOK,
		h:      http.Header{},
		w:      w,
		req:    req,
		s:      s,
		rCfg:   r,
		wh:     false,
		uin:    uin,
	}
}

func (r *Response) WroteHeader() bool {
	return r.wh
}

func (r *Response) WriteHttpHeader() {
	if r.wh == false {
		log.Debug("resp", r.h)
		// 获取响应的uin
		uin, ok := r.getUin()
		// 过滤不需要的请求头
		r.filterRespHeader()
		// 处理登录状态
		if ok && uin > 0 && r.rCfg.SignIn {
			r.s.SignIn(r.req, r.w, uin)
			log.Debug("signin", uin)
		}
		if ok && r.rCfg.SignOut {
			r.s.SignOut(r.req, r.w, r.uin)
			log.Debug("signout", r.uin)
		}
		r.wh = true
	}
}

func (r *Response) Write(p []byte) (int, error) {
	r.WriteHttpHeader()
	return r.w.Write(p)
}

func (r *Response) WriteHeader(statusCode int) {
	r.status = statusCode
	r.WriteHttpHeader()
	r.w.WriteHeader(statusCode)
}

func (r *Response) Header() http.Header {
	return r.h
}

func (r *Response) getUin() (uint64, bool) {
	if r.status != http.StatusOK {
		return 0, false
	}
	u := r.h.Get(r.s.cfg.UinHeaderName)
	uin, _ := strconv.Atoi(u)
	return uint64(uin), len(u) > 0
}

func (r *Response) filterRespHeader() {
	for k, v := range r.h {
		_, ok := r.s.hf[textproto.CanonicalMIMEHeaderKey(k)]
		if !ok {
			continue
		}
		for _, vv := range v {
			r.w.Header().Add(k, vv)
		}
	}
}
