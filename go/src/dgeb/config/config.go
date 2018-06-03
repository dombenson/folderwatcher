package config

import (
	"dgeb/defaults"
	"dgeb/interfaces"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
)

type Config struct {
	ServerDiscover     string
	ClientDiscover     string
	ServerDiscoverPort int
	ClientDiscoverPort int
	InstanceID         string
	HTTPPort           int
	MonitorPath        string
	BatchInterval      time.Duration
	BatchSize          int
	DiscoverInterval   time.Duration
	DiscoverHeartbeat  time.Duration
}

func (c *Config) GetServerAddr() string {
	return fmt.Sprintf("%s:%d", c.ServerDiscover, c.ServerDiscoverPort)
}

func (c *Config) GetClientAddr() string {
	return fmt.Sprintf("%s:%d", c.ClientDiscover, c.ClientDiscoverPort)
}

func (c *Config) GetInstanceID() string {
	return c.InstanceID
}

func (c *Config) GetHTTPPort() int {
	return c.HTTPPort
}

func (c *Config) GetPath() string {
	return c.MonitorPath
}

func (c *Config) GetBatchInterval() time.Duration {
	return c.BatchInterval
}

func (c *Config) GetBatchSize() int {
	return c.BatchSize
}

func (c *Config) GetDiscoverInterval() time.Duration {
	return c.DiscoverInterval
}

func (c *Config) GetDiscoverHeartbeat() time.Duration {
	return c.DiscoverHeartbeat
}

// Create makes a default-populated config object
// This is for internal use by individual config providers, and should
// not normally be called outside
func Create() *Config {
	c := &Config{}
	applyDefaultConfig(c)
	return c
}

// GetDefault gets a basic configuration provider
// This uses only builtin default settings
func GetDefault() interfaces.Config {
	return Create()
}

func applyDefaultConfig(c *Config) {
	c.ClientDiscover = defaults.WatcherAdvertiseAddress
	c.ServerDiscover = defaults.ServerAdvertiseAddress
	c.ClientDiscoverPort = defaults.AdvertisePort
	c.ServerDiscoverPort = defaults.AdvertisePort
	c.InstanceID = uuid.New().String()
	c.HTTPPort = defaults.HTTPPort
	// Default to current working directory.
	// No error response is helpful here
	c.MonitorPath, _ = os.Getwd()
	c.BatchInterval = defaults.NotifyInterval
	c.BatchSize = defaults.NotifyBatch
	c.DiscoverHeartbeat = defaults.DiscoverHeartbeat
	c.DiscoverInterval = defaults.DiscoverInterval
}
