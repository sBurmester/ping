package http

import (
	"net/http"
	"sync/atomic"
)

type Probe struct {
	up atomic.Bool
}

func (p *Probe) Up() {
	p.up.Store(true)
}

func (p *Probe) Down() {
	p.up.Store(false)
}

func (p *Probe) ServeHTTP(wr http.ResponseWriter, r *http.Request) {
	if p.up.Load() {
		wr.WriteHeader(http.StatusNoContent)
		return // prevent superfluous call of response.WriteHeader with StatusServiceUnavailable
	}
	wr.WriteHeader(http.StatusServiceUnavailable)
}

var (
	Ready = &Probe{}
	Live  = &Probe{}
)
