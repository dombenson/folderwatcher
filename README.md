# Building

This repo uses git submodules for external dependencies, so git clone --recursive to load them

Set GOPATH to the go folder in this directory.
Run `go install ./...`
Binaries `go/bin/master` and `go/bin/watcher` are created

# Running

The default configuration is sufficient to start and discover processes. Launch one or more master processes and one or more watcher processes.
The default watcher behaviour is to monitor the current working directory in which it is launched

# Configuration

Both components use the same configuration system. They load internal defaults, then overlay settings from `/etc/folderwatcher.ini`, then per-user settings from `~/.folderwatcher.ini`, and finally an optional file specified using the CLI flag `'-conf'`

There is an example configuration file at `example.ini`, see notes therein

# Code notes

See [`defaults/constants.go`](https://godoc.org/github.com/dombenson/folderwatcher/go/src/dgeb/defaults) for default configuration options

The layout centres around the interfaces package; the cmd/ main packages just glue together concrete implementations of them. This is intended to allow for alternatives:
* Storage (e.g. for very large datasets it might be better to use a store backed by something like Redis)
* Configuration tools (e.g. yaml/toml, to fit with other tooling)
* Communication mechanisms (particularly secure transports, rather than the simple-minded HTTP transport currently implemented)
* Reporting format (e.g. to support live websocket client updates)
* Discovery (e.g. DNS, DHCP option, broadcast, or static unicast)
* Mock implementations for tests

See [`interfaces`](https://godoc.org/github.com/dombenson/folderwatcher/go/src/dgeb/interfaces) for info about these interfaces.

It is a design feature that Config is an interface, so that a custom extension can use extra properties to set up e.g. a keypair for secure comms, while satisfying existing usage

# Features

The major feature that is present here that was not in the spec is that multiple masters all get the full list, thanks to the discovery protocol. Starting up an additional master is sufficient.

# TODO

The elephant here is that there are no tests yet; I wanted to focus on building modularity and implementing interesting feature behaviour (multiple masters, auto discovery)
