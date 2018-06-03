package interfaces

// Discoverer is used to locate watchers/masters
// Initial implementation is multicast-based
type Discoverer interface {
	// Get the current list of peers
	GetPeers() []Peer
	// Launches a goroutine to runs the discovery process on an ongoing basis, or returns an error
	Discover() error
	// Stop a running discover process
	Stop()
	// Add a callback function to run when a peer is removed
	AddRemoveCb(func(Peer))
}
