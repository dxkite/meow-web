package agent

import "net/http"

type forwardGroup struct {
	items []ForwardGroupItem
}

type ForwardGroupItem interface {
	RequestMatcher
	RequestForwardHandler
}

type forwardGroupItem struct {
	RequestMatcher
	RequestForwardHandler
}

func NewForwardGroupItem(matcher RequestMatcher, handler RequestForwardHandler) RequestForwardHandler {
	return &forwardGroupItem{RequestMatcher: matcher, RequestForwardHandler: handler}
}

func NewForwardGroup(items []ForwardGroupItem) RequestForwardHandler {
	return &forwardGroup{items: items}
}

func (fg *forwardGroup) HandleRequest(w http.ResponseWriter, req *http.Request) {
	for _, v := range fg.items {
		if v.MatchRequest(req) {
			v.HandleRequest(w, req)
			return
		}
	}
	fg.items[0].HandleRequest(w, req)
}
