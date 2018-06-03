package httpmessenger

import (
	"dgeb/fsevt"
	"dgeb/interfaces"
	"log"
)

type messenger struct {
	c interfaces.Config
}

// NewMessenger makes a new Messenger using HTTP
func NewMessenger(c interfaces.Config) interfaces.Messenger {
	m := &messenger{}
	m.c = c
	return m
}

func (m *messenger) SendFull(peer interfaces.Peer, fileList []string) error {
	log.Println("Send full:", peer, fileList)
	return nil
}
func (m *messenger) SendPartial(peer interfaces.Peer, eventList []fsevt.FsEvt) error {
	return nil
}
