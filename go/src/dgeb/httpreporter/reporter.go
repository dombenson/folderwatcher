package httpreporter

import (
	"dgeb/interfaces"
	"encoding/json"
	"net/http"
)

// HTTPReporter is a reporter that can be attached as an HTTP hander
type HTTPReporter interface {
	interfaces.Reporter
	http.Handler
}

type httpReporter struct {
	s interfaces.Storer
}

// NewReporter makes an HTTPReporter based on a supplied Storer
func NewReporter(s interfaces.Storer) HTTPReporter {
	h := &httpReporter{}
	h.s = s
	return h
}

func (h *httpReporter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fileList := h.s.GetList()
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(fileList)
}
