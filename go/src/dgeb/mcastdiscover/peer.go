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
