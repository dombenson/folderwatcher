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
	pr, pw := io.Pipe()
	msg := fullMsg{MyID: m.c.GetInstanceID(), FileList: fileList}
	msgEncoder := json.NewEncoder(pw)
	go func() {
		msgEncoder.Encode(msg)
		pw.Close()
	}()
	client := http.DefaultClient
	postURL := fmt.Sprintf("http://%s/%s/full", peer.GetAddr(), mPath)
	res, err := client.Post(postURL, "application/json", pr)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}
func (m *messenger) SendPartial(peer interfaces.Peer, eventList []fsevt.FsEvt) error {
	return nil
}
