package interfaces

import (
	"dgeb/fsevt"
	"time"
)

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
	GetHTTPPort() int
	// Monitor path
	GetPath() string
	// Send interval
	GetBatchInterval() time.Duration
	// Batch size
	GetBatchSize() int
	// Interval to send announces
	GetDiscoverInterval() time.Duration
	// Time after which to consider peer dead
	GetDiscoverHeartbeat() time.Duration
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
	// The ID of this peer
	GetID() string
}

// Watcher monitors a directory and sends listings to the peers
// identified by a discoverer
type Watcher interface {
	// The current files
	Files() []string
	// Launch a watcher, uses a Messenger to send events to peers identified by a Discoverer
	Watch(directory string) error
	// Stop a running watcher
	Stop()
}

// Messenger can send info to a Peer
type Messenger interface {
	SendFull(Peer, []string) error
	SendPartial(Peer, []fsevt.FsEvt) error
}
