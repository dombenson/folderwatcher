package mcastdiscover

import (
	"dgeb/interfaces"
	"encoding/json"
	"errors"
	"log"
	"net"
	"time"
)

type mcastAdvertiser struct {
	conf          interfaces.Config
	stopPollChan  chan (chan (struct{}))
	running       bool
	advertiseAddr string
}

func NewAdvertiser(conf interfaces.Config, advertiseAddr string) interfaces.Advertiser {
	a := &mcastAdvertiser{}
	a.conf = conf
	a.advertiseAddr = advertiseAddr
	a.stopPollChan = make(chan (chan (struct{})), 1)
	return a
}

func (a *mcastAdvertiser) Advertise() error {
	if a.running {
		return errors.New("Already running")
	}
	log.Println("Starting advertiser")
	pollAddr, err := net.ResolveUDPAddr("udp", a.advertiseAddr)
	if err != nil {
		return err
	}
	poller, err := net.DialUDP("udp", nil, pollAddr)
	if err != nil {
		return err
	}

	myInfo := pingMsg{}
	myInfo.InstanceID = a.conf.GetInstanceID()
	myInfo.Port = a.conf.GetHTTPPort()
	myInfoJSON, err := json.Marshal(myInfo)
	if err != nil {
		return err
	}

	a.running = true
	go a.poll(poller, myInfoJSON)
	return nil
}

func (a *mcastAdvertiser) poll(conn *net.UDPConn, msg []byte) {
	var rChan chan (struct{})
	stopLoop := false
	pollTicker := time.NewTicker(a.conf.GetDiscoverInterval())
	for {
		conn.Write(msg)
		select {
		case rChan = <-a.stopPollChan:
			a.stopPollChan <- rChan
			stopLoop = true
			break
		case <-pollTicker.C:
		}
		if stopLoop {
			break
		}
	}
	pollTicker.Stop()
	if rChan != nil {
		rChan <- struct{}{}
	}
}

func (a *mcastAdvertiser) Stop() {
	if !a.running {
		return
	}
	a.running = false
	log.Println("Stopping advertiser")
	resChan := make(chan (struct{}))
	a.stopPollChan <- resChan
	<-resChan
	log.Println("Stopped advertiser")
}
