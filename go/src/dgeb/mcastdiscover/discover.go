package mcastdiscover

import (
	"dgeb/defaults"
	"dgeb/interfaces"
	"encoding/json"
	"errors"
	"log"
	"net"
	"time"
)

type mcastDiscoverer struct {
	conf         interfaces.Config
	stopPollChan chan (chan (struct{}))
	running      bool
	listener     *net.UDPConn
	peers        *map[string]interfaces.Peer
}

// NewDiscoverer makes a new multicast discoverer
func NewDiscoverer(conf interfaces.Config) interfaces.Discoverer {
	peerList := make(map[string]interfaces.Peer)
	d := &mcastDiscoverer{
		conf:         conf,
		stopPollChan: make(chan (chan (struct{}))),
		peers:        &peerList,
	}

	return d
}

func (d *mcastDiscoverer) GetPeers() []interfaces.Peer {
	curMap := *d.peers
	peerList := make([]interfaces.Peer, 0, len(curMap))
	for _, v := range curMap {
		peerList = append(peerList, v)
	}
	return peerList
}

func (d *mcastDiscoverer) Discover(discoverAddr, advertiseAddr string) error {
	if d.running {
		return errors.New("Already running")
	}
	log.Println("Starting discoverer")
	listenAddr, err := net.ResolveUDPAddr("udp", advertiseAddr)
	if err != nil {
		return err
	}
	pollAddr, err := net.ResolveUDPAddr("udp", discoverAddr)
	if err != nil {
		return err
	}
	listener, err := net.ListenMulticastUDP("udp", nil, listenAddr)
	if err != nil {
		return err
	}
	poller, err := net.DialUDP("udp", nil, pollAddr)
	if err != nil {
		return err
	}

	myInfo := pingMsg{}
	myInfo.InstanceID = d.conf.GetInstanceID()
	myInfo.Port = d.conf.GetHTTPPort()
	myInfoJSON, err := json.Marshal(myInfo)
	if err != nil {
		return err
	}

	listener.SetReadBuffer(defaults.ReadBufSize)
	d.listener = listener
	go d.listen()
	go d.poll(poller, myInfoJSON)
	return nil
}

func (d *mcastDiscoverer) Stop() {
	if !d.running {
		return
	}
	log.Println("Stopping discoverer")
	resChan := make(chan (struct{}))
	d.stopPollChan <- resChan
	<-resChan
	log.Println("Stopped discoverer")
}

func (d *mcastDiscoverer) poll(conn *net.UDPConn, msg []byte) {
	var rChan chan (struct{})
	d.running = true
	stopLoop := false
	for {
		conn.Write(msg)
		select {
		case rChan = <-d.stopPollChan:
			stopLoop = true
			break
		case <-time.After(5 * time.Second):
		}
		if stopLoop {
			break
		}
	}
	d.listener.Close()
	d.running = false
	if rChan != nil {
		rChan <- struct{}{}
	}
}

func (d *mcastDiscoverer) listen() {
	var buf = make([]byte, defaults.ReadBufSize)
	for {
		numBytes, addr, err := d.listener.ReadFromUDP(buf)
		if err != nil {
			log.Println(err)
			break
		}
		peerInfo := pingMsg{}
		err = json.Unmarshal(buf[:numBytes], &peerInfo)
		if err != nil {
			log.Println(err)
			continue
		}
		curPeers := *d.peers
		curInfo, found := curPeers[peerInfo.InstanceID]
		if !found {
			thisPeer := &peer{}
			thisPeer.ipAddr = addr.IP.String()
			thisPeer.port = peerInfo.Port
			thisPeer.lastSeen = time.Now()
			thisPeer.instanceID = peerInfo.InstanceID
			newPeers := make(map[string]interfaces.Peer, len(curPeers)+1)
			for k, v := range curPeers {
				newPeers[k] = v
				log.Println(v)
			}
			newPeers[peerInfo.InstanceID] = thisPeer
			log.Println(thisPeer)
			d.peers = &newPeers
		} else {
			curInfo.(*peer).lastSeen = time.Now()
		}
	}
}
