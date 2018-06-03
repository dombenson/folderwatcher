package interfaces

import (
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
