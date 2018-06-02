package defaults

// GlobalConfigFile is the default global configuration ini path
const GlobalConfigFile = "/etc/folderwatcher.ini"

// UserConfigFile is the default path for a per-user config file
const UserConfigFile = "~/.folderwatcher.ini"

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
