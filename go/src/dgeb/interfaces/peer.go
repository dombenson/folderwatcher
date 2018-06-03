package interfaces

import "time"

// Peer represents a watcher/master, and how to communicate with it
type Peer interface {
	// The address to contact this peer
	GetAddr() string
	// How long since the peer has reported in
	StaleTime() time.Duration
	// The ID of this peer
	GetID() string
	// Is currently marked as stale?
	IsStale() bool
}
