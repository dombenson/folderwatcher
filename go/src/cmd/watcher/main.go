package main

import (
	"dgeb/iniconfig"
	"dgeb/mcastdiscover"
	"flag"
	"log"
	"time"
)

func main() {
	var confPath string
	flag.StringVar(&confPath, "conf", "", "The path to the config file to load")
	flag.Parse()
	conf := iniconfig.Get(confPath)

	log.Println("Starting master with ID: ", conf.GetInstanceID())

	discoverer := mcastdiscover.NewDiscoverer(conf)

	err := discoverer.Discover(conf.GetServerAddr(), conf.GetClientAddr())
	if err != nil {
		panic(err)
	}
	time.Sleep(10 * time.Second)
	discoverer.Stop()
}
