package agent

import "net/http"

type forwardGroup struct {
	items []MatchForwardHandler
}

type MatchForwardHandler interface {
	RequestMatcher
	RequestForwardHandler
}

type matchForwardHandler struct {
	RequestMatcher
	RequestForwardHandler
}

func NewMatchForwardHandler(matcher RequestMatcher, handler RequestForwardHandler) MatchForwardHandler {
	return &matchForwardHandler{RequestMatcher: matcher, RequestForwardHandler: handler}
}

func NewForwardGroup(items []MatchForwardHandler) RequestForwardHandler {
	return &forwardGroup{items: items}
}

func (fg *forwardGroup) HandleRequest(w http.ResponseWriter, req *http.Request) {
	// 使用自定义处理请求转发
	for _, v := range fg.items {
		if v.MatchRequest(req) {
			v.HandleRequest(w, req)
			return
		}
	}

	// 无匹配，走第一条
	fg.items[0].HandleRequest(w, req)
}
