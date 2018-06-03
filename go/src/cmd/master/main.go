package main

import (
	"context"
	"dgeb/config/iniconfig"
	"dgeb/httpmessenger"
	"dgeb/httpreporter"
	"dgeb/mcastdiscover"
	"dgeb/memorystorer"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"goji.io/pat"

	"goji.io"
)

func main() {
	var confPath string
	flag.StringVar(&confPath, "conf", "", "The path to the config file to load")
	flag.Parse()
	conf := iniconfig.Get(confPath)

	log.Println("Starting master with ID: ", conf.GetInstanceID())
	defer log.Println("Shutdown complete")

	discoverer := mcastdiscover.NewDiscoverer(conf, conf.GetServerAddr())
	advertiser := mcastdiscover.NewAdvertiser(conf, conf.GetClientAddr())
	storer := memorystorer.NewStorer()
	receiver := httpmessenger.NewReceiver(conf, storer)

	discoverer.AddRemoveCb(storer.RemovePeer)

	err := discoverer.Discover()
	if err != nil {
		panic(err)
	}
	defer discoverer.Stop()
	err = advertiser.Advertise()
	if err != nil {
		panic(err)
	}
	defer advertiser.Stop()

	reporter := httpreporter.NewReporter(storer)

	httpMux := goji.NewMux()
	receiver.AddMux(httpMux)
	httpMux.Handle(pat.Get("/"), reporter)

	listenAddr := fmt.Sprintf(":%d", conf.GetHTTPPort())
	httpServer := &http.Server{Addr: listenAddr, Handler: httpMux}
	log.Println("Starting webserver on " + listenAddr)
	go httpServer.ListenAndServe()

	defer func() {
		log.Println("Stopping webserver")
		ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
		httpServer.Shutdown(ctx)
		ctxCancel()
	}()

	stopChan := make(chan (os.Signal))
	signal.Notify(stopChan, os.Interrupt)
	<-stopChan

	log.Println("Shutting down")
}
