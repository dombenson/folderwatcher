package httpmessenger

import (
	"dgeb/fsevt"
	"dgeb/interfaces"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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
	msg := fullMsg{MyID: m.c.GetInstanceID(), FileList: fileList}
	return m.doSend(peer, msg, fullPath)
}
func (m *messenger) SendPartial(peer interfaces.Peer, eventList []fsevt.FsEvt) error {
	log.Println("Send partial:", peer, eventList)
	msg := partialMsg{MyID: m.c.GetInstanceID(), EventList: eventList}
	return m.doSend(peer, msg, partialPath)
}

func (m *messenger) doSend(peer interfaces.Peer, msg interface{}, endpoint string) error {
	pr, pw := io.Pipe()
	msgEncoder := json.NewEncoder(pw)
	go func() {
		msgEncoder.Encode(msg)
		pw.Close()
	}()
	client := http.DefaultClient
	postURL := fmt.Sprintf("http://%s/%s/%s", peer.GetAddr(), mPath, endpoint)
	res, err := client.Post(postURL, "application/json", pr)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}
