package main

import (
	"fmt"
	"os"

	"github.com/shaunmolloy/bugbox/cmd/setup"
)

func main() {
	if err := setup.Setup(); err != nil {
		fmt.Fprintf(os.Stderr, "Setup failed: %v\n", err)
		os.Exit(1)
	}
}
