package main

import (
	"github.com/lspaccatrosi16/spac/lib/error"
	"github.com/lspaccatrosi16/spac/lib/foo"
	"github.com/lspaccatrosi16/spac/lib/hello"
	"github.com/lspaccatrosi16/spac/lib/types"
)

func main() {
	manager := types.NewManager()

	manager.Register("hello", "Says hello to you", hello.Hello)
	manager.Register("foo", "FooBarBiz", foo.Foo)
	manager.Register("error", "Simulates an error being raised", error.Error)

	manager.Gui()
}

// what do I want this to do?

// setup environment

// setup environment:
// install aup
// add scaffold to aup
// install useful things (support different package managers e.g. dnf, apt)
