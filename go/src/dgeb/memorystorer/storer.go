package memorystorer

import (
	"dgeb/fsevt"
	"dgeb/interfaces"
	"errors"
	"sort"
	"sync"
)

type peerMap struct {
	m      map[string]bool
	l      sync.RWMutex
	inited bool
}

type storer struct {
	peerListMap map[string]*peerMap
	mapLock     sync.RWMutex
}

// NewStorer makes a new memory-based storer
func NewStorer() interfaces.Storer {
	s := &storer{}
	s.peerListMap = make(map[string]*peerMap)
	return s
}

func (s *storer) getPeerEntry(peerID string) *peerMap {
	s.mapLock.RLock()
	p, ok := s.peerListMap[peerID]
	s.mapLock.RUnlock()
	if !ok {
		s.mapLock.Lock()
		p, ok = s.peerListMap[peerID]
		if !ok {
			p = &peerMap{}
			p.m = make(map[string]bool)
			s.peerListMap[peerID] = p
		}
		s.mapLock.Unlock()
	}
	return p
}

func (s *storer) SetFull(peerID string, fileList []string) error {
	p := s.getPeerEntry(peerID)
	p.l.Lock()
	defer p.l.Unlock()
	p.m = make(map[string]bool)
	for _, v := range fileList {
		p.m[v] = true
	}
	p.inited = true
	return nil
}

func (s *storer) AddEvents(peerID string, eventList []fsevt.FsEvt) error {
	p := s.getPeerEntry(peerID)
	p.l.Lock()
	defer p.l.Unlock()
	if !p.inited {
		return errors.New("Cannot add events to uninitialised list")
	}
	for _, evt := range eventList {
		switch evt.Type {
		case fsevt.FsEvtAdd:
			p.m[evt.Name] = true
		case fsevt.FsEvtDel:
			delete(p.m, evt.Name)
		default:
			return errors.New("Unhandled event type")
		}
	}
	return nil
}

func (s *storer) doGetCurPeerList() []*peerMap {
	s.mapLock.RLock()
	defer s.mapLock.RUnlock()
	maps := make([]*peerMap, 0, len(s.peerListMap))
	for _, v := range s.peerListMap {
		maps = append(maps, v)
	}
	return maps
}
func (s *storer) addOnePeerMap(fullList []string, v *peerMap) []string {
	v.l.RLock()
	defer v.l.RUnlock()
	for f := range v.m {
		fullList = append(fullList, f)
	}
	return fullList
}

func (s *storer) GetList() []string {
	fullList := make([]string, 0, 100)
	maps := s.doGetCurPeerList()
	for _, v := range maps {
		fullList = s.addOnePeerMap(fullList, v)
	}
	sort.Strings(fullList)
	return fullList
}

func (s *storer) RemovePeer(peer interfaces.Peer) {
	s.mapLock.Lock()
	defer s.mapLock.Unlock()
	delete(s.peerListMap, peer.GetID())
}
