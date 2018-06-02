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

	err := discoverer.Discover(conf.GetClientAddr(), conf.GetServerAddr())
	if err != nil {
		panic(err)
	}
	time.Sleep(10 * time.Second)
	peers := discoverer.GetPeers()

	discoverer.Stop()

	for _, v := range peers {
		log.Println(v.GetAddr())
	}
}
