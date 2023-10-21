package hello

import (
	"fmt"

	"github.com/lspaccatrosi16/go-cli-tools/input"
)

func Hello() error {
	name := input.GetInput("Name")

	fmt.Printf("Hello, %s!\n", name)
	return nil
}
