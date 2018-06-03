package interfaces

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
