package fsnotifywatcher

import (
	"dgeb/fsevt"
	"dgeb/interfaces"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"time"

	"github.com/fsnotify/fsnotify"
)

type watcher struct {
	fileInfo     map[string]os.FileInfo
	peerStates   map[string]bool
	stopPollChan chan (chan (string))
	running      bool
	c            interfaces.Config
	m            interfaces.Messenger
	d            interfaces.Discoverer
}

// NewWatcher makes a new fsnotify-based watcher
func NewWatcher(c interfaces.Config, m interfaces.Messenger, d interfaces.Discoverer) interfaces.Watcher {
	w := &watcher{}
	w.fileInfo = make(map[string]os.FileInfo)
	w.peerStates = make(map[string]bool)
	w.stopPollChan = make(chan (chan (string)), 1)
	w.c = c
	w.m = m
	w.d = d
	return w
}

func (w *watcher) Files() []string {
	ret := make([]string, 0, len(w.fileInfo))
	for k := range w.fileInfo {
		ret = append(ret, k)
	}
	sort.Strings(ret)
	return ret
}

func (w *watcher) Stop() {
	if !w.running {
		return
	}
	log.Println("Stopping watcher")
	resChan := make(chan (string))
	w.stopPollChan <- resChan
	n := <-resChan
	log.Println("Stopped watcher subprocess:", n)
	n = <-resChan
	log.Println("Stopped watcher subprocess:", n)
	log.Println("Stopped watcher")
}

func (w *watcher) Watch(directory string) error {
	if w.running {
		return errors.New("Already running")
	}
	initList, err := ioutil.ReadDir(directory)
	if err != nil {
		return err
	}
	for _, v := range initList {
		w.fileInfo[v.Name()] = v
	}
	fsn, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	filteredEvtChan := make(chan (fsevt.FsEvt), 100)
	fsn.Add(directory)
	go w.processEvents(filteredEvtChan)
	go w.doWatch(fsn, filteredEvtChan)
	return nil
}

func (w *watcher) sendEvents(evts []fsevt.FsEvt, d interfaces.Discoverer) {
	for _, evt := range evts {
		switch evt.Type {
		case fsevt.FsEvtAdd:
			stat, err := os.Stat(evt.Name)
			if err != nil {
				break
			}
			w.fileInfo[stat.Name()] = stat
		case fsevt.FsEvtDel:
			delete(w.fileInfo, evt.Name)
		}
	}
	peers := d.GetPeers()
	for _, peer := range peers {
		lastOk, found := w.peerStates[peer.GetID()]
		if found && lastOk {
			err := w.m.SendPartial(peer, evts)
			if err != nil {
				w.peerStates[peer.GetID()] = false
				log.Println("Partial update error: ", err)
			}
		} else {
			err := w.m.SendFull(peer, w.Files())
			if err != nil {
				log.Println("Full update error: ", err)
			} else {
				w.peerStates[peer.GetID()] = true
			}
		}
	}
}

func (w *watcher) processEvents(filteredEvtChan chan (fsevt.FsEvt)) {
	batchSize := w.c.GetBatchSize()
	sendTimer := time.NewTicker(w.c.GetBatchInterval())
	evtQueue := make([]fsevt.FsEvt, 0, batchSize)
	var stop bool
	for {
		select {
		case <-sendTimer.C:
			w.sendEvents(evtQueue, w.d)
			evtQueue = make([]fsevt.FsEvt, 0, batchSize)
		case evt := <-filteredEvtChan:
			evtQueue = append(evtQueue, evt)
			if len(evtQueue) >= batchSize {
				w.sendEvents(evtQueue, w.d)
				evtQueue = make([]fsevt.FsEvt, 0, batchSize)
			}
		case resChan := <-w.stopPollChan:
			w.sendEvents(evtQueue, w.d)
			resChan <- "ProcessEvents"
			w.stopPollChan <- resChan
			stop = true
		}
		if stop {
			break
		}
	}
}

func (w *watcher) doWatch(fsn *fsnotify.Watcher, filteredEvtChan chan (fsevt.FsEvt)) {
	w.running = true
	for {
		stop := false
		select {
		case fsEvt := <-fsn.Events:
			switch fsEvt.Op {
			case fsnotify.Create:
				filteredEvtChan <- fsevt.FsEvt{Type: fsevt.FsEvtAdd, Name: fsEvt.Name}
			case fsnotify.Remove, fsnotify.Rename:
				filteredEvtChan <- fsevt.FsEvt{Type: fsevt.FsEvtDel, Name: fsEvt.Name}
			}
		case fsErr := <-fsn.Errors:
			log.Println("Nofify error:", fsErr)
		case resChan := <-w.stopPollChan:
			resChan <- "WatchFolder"
			w.stopPollChan <- resChan
			stop = true
		}
		if stop {
			break
		}
	}
}
