package interfaces

import "dgeb/fsevt"

// Storer receives and stores lists from Receivers
type Storer interface {
	SetFull(peerID string, fileList []string) error
	AddEvents(peerID string, eventList []fsevt.FsEvt) error
	GetList() []string
	RemovePeer(peer Peer)
}
