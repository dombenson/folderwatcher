package interfaces

import "time"

// The Config interface is required to be satisfied by configuration providers
// It exposes semantic configuration properties
type Config interface {
	// server:port, usable directly as bind addr on the master end, and to send data from the watcher
	GetServerAddr() string
	// sever:port, usable directly as bind addr on the watcher end, and to send notification from the master
	GetClientAddr() string
	// An identifier, unique to this instance. Default use a UUID
	GetInstanceID() string
	// HTTP listener port
	GetHttpPort() int
}

// Discoverer is used to locate watchers/masters
// Initial implementation is multicast-based
type Discoverer interface {
	// Get the current list of peers
	GetPeers() []Peer
	// Launches a goroutine to runs the discovery process on an ongoing basis, or returns an error
	Discover(discoverAddr, advertiseAddr string) error
	// Stop a running discover process
	Stop()
}

// Peer represents a watcher/master, and how to communicate with it
type Peer interface {
	// The address to contact this peer
	GetAddr() string
	// How long since the peer has reported in
	StaleTime() time.Duration
}
