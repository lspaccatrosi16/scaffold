package main

import (
	"example.com/my/project/lib/error"
	"example.com/my/project/lib/foo"
	"example.com/my/project/lib/hello"
	"example.com/my/project/lib/types"
)

func main() {
	manager := types.NewManager()

	manager.Register("hello", "Says hello to you", hello.Hello)
	manager.Register("foo", "FooBarBiz", foo.Foo)
	manager.Register("error", "Simulates an error being raised", error.Error)

	manager.Gui()
}
