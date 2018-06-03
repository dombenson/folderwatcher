package defaults

import (
	"time"
)

// GlobalConfigFile is the default global configuration ini path
const GlobalConfigFile = "/etc/folderwatcher.ini"

// UserConfigFile is the default path for a per-user config file
const UserConfigFile = ".folderwatcher.ini"

// ServerAdvertiseAddress is the multicast IP listened to by masters and probed by watchers
const ServerAdvertiseAddress = "224.0.240.1"

// WatcherAdvertiseAddress is the multicast IP listened to by watchers and probed by masters
const WatcherAdvertiseAddress = "224.0.240.2"

// AdvertisePort is the default port for multicast advertisement and discovery
const AdvertisePort = 9100

// HTTPPort is the default port for the master HTTP interface
const HTTPPort = 8080

// ReadBufSize is the default buffer size for receiving UDP datagrams; 1500 more than sufficient for 1500-MTU
const ReadBufSize = 1500

// NotifyInterval is the dead time over which events will be batched
// It is also the time to notify new masters
const NotifyInterval = 150 * time.Millisecond

// NotifyBatch is the max number of events to send in one batch
// If this number is reached, the batch will be sent regardless of
// notify interval
const NotifyBatch = 20

// DiscoverInterval is the period to send discover announces
const DiscoverInterval = 5 * time.Second

// DiscoverHeartbeat is the age at which to consider a peer dead
const DiscoverHeartbeat = 15 * time.Second
