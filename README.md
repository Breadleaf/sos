# Simple Object Store

Simple Object Store (or `sos` for short) is an extremely lightweight object
storage tool written in golang. It provides a cli tool to interact with the
server as well as a server binary. It also exposes its internal client library
to allow for interacting with the store from your own projects.

More documentation to come as I inevitably need to trouble shoot the code's
ability to import and integrate into other code bases.

## Install

For now, until I write a better system in `Makefile`, run the following:

```bash
go mod tidy
go build -o sos-cli ./cli
go build -o sos-server ./server
cd ./dashboard && go build -o sos-dashboard # should be ran inside of ./dashboard/
```

I would not recommend installing this version of the code until I create
automation to manage the install. You *can* install it to `~/.local/bin/` if you
really want to run it from anywhere in your file system.

## Usage

Near all functionality is self documented from the cli tool. To get started run
`sos` and navigate the menus until I write better docs.

## API

[sos example project](https://github.com/Breadleaf/sos)
