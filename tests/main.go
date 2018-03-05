package main

import (
	"fmt"
	"os"

	"github.com/EgorMatirov/xrootd"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s host:port ", os.Args[0])
		os.Exit(1)
	}
	_, err := xrootd.Connect(os.Args[1])

	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}
