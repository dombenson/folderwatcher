package httpmessenger

import "dgeb/fsevt"

const mPath = "w"

const fullPath = "full"
const partialPath = "partial"

type fullMsg struct {
	MyID     string   `json:"peerid"`
	FileList []string `json:"files"`
}

type partialMsg struct {
	MyID      string        `json:"peerid"`
	EventList []fsevt.FsEvt `json:"events"`
}
