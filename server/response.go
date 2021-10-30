package server

import (
	"dxkite.cn/gateway/route"
	"dxkite.cn/log"
	"net/http"
	"net/textproto"
	"strconv"
)

type AutoResponse struct {
	w http.ResponseWriter
	s *Server
	r *route.RouteConfig
	// 头部
	h http.Header
	// 是否处理了头部
	wh bool
	// HTTP
	status int
}

func NewAutoResponse(s *Server, r *route.RouteConfig, w http.ResponseWriter) http.ResponseWriter {
	return &AutoResponse{
		status: http.StatusOK,
		h:      http.Header{},
		w:      w,
		s:      s,
		r:      r,
		wh:     false,
	}
}

func (br *AutoResponse) Write(p []byte) (int, error) {
	if br.wh == false {
		br.normalizeResp()
		// 自动写入UIN
		uin := br.getUin()
		if uin > 0 {
			if br.r.SignIn {
				br.s.SignIn(br.w, uin)
			}
			if br.r.SignOut {
				br.s.SignOut(br.w, uin)
			}
			log.Debug("uin write", uin)
		}
		br.wh = true
	}
	return br.w.Write(p)
}

func (br *AutoResponse) WriteHeader(statusCode int) {
	br.status = statusCode
}

func (br *AutoResponse) Header() http.Header {
	return br.h
}

func (br *AutoResponse) getUin() uint64 {
	if br.status != http.StatusOK {
		return 0
	}
	uin, _ := strconv.Atoi(br.h.Get(br.s.cfg.UinHeaderName))
	return uint64(uin)
}

func (br *AutoResponse) normalizeResp() {
	for k, v := range br.h {
		_, ok := br.s.hf[textproto.CanonicalMIMEHeaderKey(k)]
		if !ok {
			continue
		}
		for _, vv := range v {
			br.w.Header().Set(k, vv)
		}
	}
}
