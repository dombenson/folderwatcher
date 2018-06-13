package httpmessenger

import (
	"dgeb/interfaces"
	"encoding/json"
	"log"
	"net/http"

	"goji.io"
	"goji.io/pat"
)

type receiver struct {
	c interfaces.Config
	s interfaces.Storer
}

// Receiver extends the base receiver with a registration method for a mux
type Receiver interface {
	interfaces.Receiver
	AddMux(*goji.Mux)
}

// NewReceiver makes a new Receiver using HTTP
func NewReceiver(c interfaces.Config, s interfaces.Storer) Receiver {
	r := &receiver{}
	r.c = c
	r.s = s
	return r
}

func (r *receiver) AddMux(m *goji.Mux) {
	mux := goji.SubMux()
	pattern := pat.New("/" + mPath + "/*")
	m.Handle(pattern, mux)

	mux.HandleFunc(pat.Post("/"+fullPath), r.handleFull)
	mux.HandleFunc(pat.Post("/"+partialPath), r.handlePartial)
}

func (r *receiver) handleFull(w http.ResponseWriter, req *http.Request) {
	msg := fullMsg{}
	bodyDecoder := json.NewDecoder(req.Body)
	err := bodyDecoder.Decode(&msg)
	if err != nil {
		w.WriteHeader(400)
		log.Println(err)
	}
	err = r.s.SetFull(msg.MyID, msg.FileList)
	if err != nil {
		w.WriteHeader(500)
		log.Println(err)
	}
}

func (r *receiver) handlePartial(w http.ResponseWriter, req *http.Request) {
	msg := partialMsg{}
	bodyDecoder := json.NewDecoder(req.Body)
	err := bodyDecoder.Decode(&msg)
	if err != nil {
		w.WriteHeader(400)
		log.Println(err)
	}
	err = r.s.AddEvents(msg.MyID, msg.EventList)
	if err != nil {
		w.WriteHeader(500)
		log.Println(err)
	}
}
