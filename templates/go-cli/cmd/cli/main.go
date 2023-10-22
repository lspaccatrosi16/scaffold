package main

import (
	"example.com/my/project/lib/error"
	"example.com/my/project/lib/foo"
	"example.com/my/project/lib/hello"
	"github.com/lspaccatrosi16/go-cli-tools/command"
)

func main() {
	manager := command.NewManager(command.ManagerConfig{Searchable: true})

	manager.Register("hello", "Says hello to you", hello.Hello)
	manager.Register("foo", "FooBarBiz", foo.Foo)
	manager.Register("error", "Simulates an error being raised", error.Error)

	manager.Tui()
}
