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
	defer log.Println("Shutdown")

	discoverer := mcastdiscover.NewDiscoverer(conf)

	messenger := httpmessenger.NewMessenger(conf)

	watcher := fsnotifywatcher.NewWatcher(conf, messenger, discoverer)

	err := discoverer.Discover(conf.GetServerAddr(), conf.GetClientAddr())
	if err != nil {
		panic(err)
	}
	defer discoverer.Stop()

	err = watcher.Watch(conf.GetPath())
	if err != nil {
		panic(err)
	}
	defer watcher.Stop()

	stopChan := make(chan (os.Signal))
	signal.Notify(stopChan, os.Interrupt)
	<-stopChan
	log.Println(watcher.Files())

}
