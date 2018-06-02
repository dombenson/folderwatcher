package iniconfig

import (
	"dgeb/defaults"
	"dgeb/interfaces"
	"fmt"

	"github.com/google/uuid"
)

type config struct {
	serverDiscover string
	clientDiscover string
	discoverPort   int
	instanceID     string
	httpPort       int
}

func (c *config) GetServerAddr() string {
	return fmt.Sprintf("%s:%d", c.serverDiscover, c.discoverPort)
}

func (c *config) GetClientAddr() string {
	return fmt.Sprintf("%s:%d", c.clientDiscover, c.discoverPort)
}

func (c *config) GetInstanceID() string {
	return c.instanceID
}

func (c *config) GetHttpPort() int {
	return c.httpPort
}

// Get configuration from ini file
func Get(preferredFile string) interfaces.Config {
	c := &config{}
	c.clientDiscover = defaults.WatcherAdvertiseAddress
	c.serverDiscover = defaults.ServerAdvertiseAddress
	c.discoverPort = defaults.AdvertisePort
	c.instanceID = uuid.New().String()
	c.httpPort = defaults.HTTPPort
	return c
}
