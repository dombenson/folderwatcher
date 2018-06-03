package main

import (
	"dgeb/config/iniconfig"
	"dgeb/fsnotifywatcher"
	"dgeb/httpmessenger"
	"dgeb/mcastdiscover"
	"flag"
	"log"
	"os"
	"os/signal"
)

func main() {
	var confPath string
	flag.StringVar(&confPath, "conf", "", "The path to the config file to load")
	flag.Parse()
	conf := iniconfig.Get(confPath)

	log.Println("Starting watcher with ID: ", conf.GetInstanceID())
	defer log.Println("Shutdown complete")

	discoverer := mcastdiscover.NewDiscoverer(conf, conf.GetClientAddr())
	advertiser := mcastdiscover.NewAdvertiser(conf, conf.GetServerAddr())

	messenger := httpmessenger.NewMessenger(conf)

	watcher := fsnotifywatcher.NewWatcher(conf, messenger, discoverer)

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

	err = watcher.Watch(conf.GetPath())
	if err != nil {
		panic(err)
	}
	defer watcher.Stop()

	stopChan := make(chan (os.Signal))
	signal.Notify(stopChan, os.Interrupt)
	<-stopChan
	log.Println("Shutting down")

}
