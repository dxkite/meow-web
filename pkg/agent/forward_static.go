package agent

import (
	"math/rand"
	"net/http"
	"time"
)

type StaticForwardHandler struct {
	Timeout int
	Target  []*EndpointTarget
}

func NewStaticForwardHandler() *StaticForwardHandler {
	return new(StaticForwardHandler)
}

type EndpointTarget struct {
	Network string
	Address string
}

func (h *StaticForwardHandler) HandleRequest(w http.ResponseWriter, req *http.Request) {
	n := len(h.Target)
	i := intn(n)
	target := h.Target[i]
	f := NewForward(target.Network, target.Address, time.Duration(h.Timeout)*time.Millisecond)
	f.HandleRequest(w, req)
}

func intn(v int) int {
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	return rd.Intn(v)
}
