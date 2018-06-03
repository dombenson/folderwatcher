package httpmessenger

import "dgeb/fsevt"

const mPath = "w"

type fullMsg struct {
	MyID     string   `json:"peerid"`
	FileList []string `json:"files"`
}

type partialMsg struct {
	MyID      string        `json:"peerid"`
	EventList []fsevt.FsEvt `json:"events"`
}
