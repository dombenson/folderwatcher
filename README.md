# Building

This repo uses git submodules for external dependencies, so git clone --recursive to load them

Set GOPATH to the go folder in this directory.
Run `go install ./...`
Binaries go/bin/master and go/bin/watcher are created

# Running

The default configuration is sufficient to start and discover processes. Launch one or more master processes and one or more watcher processes.
The default watcher behaviour is to monitor the current working directory in which it is launched

# Configuration

Both components use the same configuration system. They load internal defaults, then overlay settings from /etc/folderwatcher.ini, then per-user settings from ~/.folderwatcher.ini, and finally an optional file specified using the CLI flag '-conf'

There is an example configuration file at example.ini, see notes therein

# Code notes

See `defaults/constants.go` for default configuration options

The layout centres around the interfaces package; the cmd/ main packages just glue together concrete implementations of them. This is intended to allow for alternatives:
* Storage (e.g. for very large datasets it might be better to use a store backed by something like Redis)
* Configuration tools (e.g. yaml/toml, to fit with other tooling)
* Communication mechanisms (e.g. socket based rather than individual HTTP connections)
* Reporting format (e.g. to support live websocket client updates)
* Discovery (e.g. DNS, DHCP option, broadcast, or static unicast)
* Mock implementations for tests

# TODO

The elephant here is that there are no tests yet; I wanted to focus on building modularity and implementing interesting feature behaviour (multiple masters, auto discovery)