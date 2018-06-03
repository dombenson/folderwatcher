package interfaces

// Advertiser is used to permit discovery by a Discoverer
type Advertiser interface {
	// Launches a goroutine to runs the discovery process on an ongoing basis, or returns an error
	Advertise() error
	// Stop a running discover process
	Stop()
}
