package mcastdiscover

import (
	"fmt"
	"time"
)

type peer struct {
	ipAddr     string
	port       int
	instanceID string
	lastSeen   time.Time
	stale      bool
}

func (p *peer) GetAddr() string {
	return fmt.Sprintf("%s:%d", p.ipAddr, p.port)
}

func (p *peer) StaleTime() time.Duration {
	return time.Since(p.lastSeen)
}

func (p *peer) GetID() string {
	return p.instanceID
}
func (p *peer) setStale(stale bool) bool {
	ret := stale && !p.stale
	p.stale = stale
	return ret
}
func (p *peer) IsStale() bool {
	return p.stale
}
