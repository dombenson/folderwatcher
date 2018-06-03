package interfaces

import "dgeb/fsevt"

// Messenger can send info to a Peer
type Messenger interface {
	SendFull(Peer, []string) error
	SendPartial(Peer, []fsevt.FsEvt) error
}
