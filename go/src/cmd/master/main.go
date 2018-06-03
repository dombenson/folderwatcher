package main

import (
	"context"
	"dgeb/config/iniconfig"
	"dgeb/httpmessenger"
	"dgeb/mcastdiscover"
	"dgeb/memorystorer"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"goji.io"
)

func main() {
	var confPath string
	flag.StringVar(&confPath, "conf", "", "The path to the config file to load")
	flag.Parse()
	conf := iniconfig.Get(confPath)

	log.Println("Starting master with ID: ", conf.GetInstanceID())

	discoverer := mcastdiscover.NewDiscoverer(conf)
	storer := memorystorer.NewStorer()
	receiver := httpmessenger.NewReceiver(conf, storer)

	discoverer.AddRemoveCb(storer.RemovePeer)

	err := discoverer.Discover(conf.GetClientAddr(), conf.GetServerAddr())
	if err != nil {
		panic(err)
	}

	httpMux := goji.NewMux()
	receiver.AddMux(httpMux)

	listenAddr := fmt.Sprintf(":%d", conf.GetHTTPPort())
	httpServer := &http.Server{Addr: listenAddr, Handler: httpMux}
	go httpServer.ListenAndServe()

	stopChan := make(chan (os.Signal))
	signal.Notify(stopChan, os.Interrupt)
	<-stopChan

	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	httpServer.Shutdown(ctx)
	ctxCancel()

	peers := discoverer.GetPeers()

	discoverer.Stop()

	for _, v := range peers {
		log.Println(v.GetAddr())
	}

	log.Println(storer.GetList())
	log.Println("Shutdown")
}
