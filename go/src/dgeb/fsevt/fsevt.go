package fsevt

type fsEvtType int

const (
	// FsEvtAdd indicates a file has been added
	FsEvtAdd fsEvtType = iota
	// FsEvtDel indicates a file has been removed or renamed
	FsEvtDel
)

// FsEvt Interface for filesystem events to send
type FsEvt struct {
	Type fsEvtType `json:"type"`
	Name string    `json:"name"`
}
