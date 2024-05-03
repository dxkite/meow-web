package agent

import (
	"math/rand"
	"time"
)

type StaticForwardHandler struct {
	*BasicForwardHandler
	targets []*EndpointTarget
	timeout int
}

type EndpointTarget struct {
	Network string
	Address string
}

func NewStaticForwardHandler(targets []*EndpointTarget, timeout int) *StaticForwardHandler {
	h := new(StaticForwardHandler)
	h.BasicForwardHandler = NewBasicForwardHandler(h)
	h.timeout = timeout
	h.targets = targets
	return h
}

func (h *StaticForwardHandler) ForwardTarget() (network, address string, timeout time.Duration) {
	n := len(h.targets)
	i := intn(n)
	network = h.targets[i].Network
	address = h.targets[i].Address
	timeout = time.Duration(h.timeout) * time.Millisecond
	return
}

func intn(v int) int {
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	return rd.Intn(v)
}
