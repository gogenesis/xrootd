package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/EgorMatirov/xrootd"
)

var addr = flag.String("addr", "0.0.0.0:9001", "address of xrootd server")

func main() {
	flag.Parse()

	client, err := xrootd.New(context.Background(), *addr)
	checkError(err)

	response, securityInfo, err := client.Protocol(context.Background())
	checkError(err)
	log.Printf("Protocol binary version is %d. Security level is %d.", response.BinaryProtocolVersion, securityInfo.SecurityLevel)

	loginResult, err := client.Login(context.Background(), "gopher")
	checkError(err)
	log.Printf("Logged in! Security information length is %d. Value is \"%s\"\n", len(loginResult.SecurityInformation), loginResult.SecurityInformation)

	err = client.Ping(context.Background())
	checkError(err)
	log.Print("Pong!")

	dirs, err := client.Dirlist(context.Background(), "/tmp/")
	checkError(err)
	log.Printf("dirlist /tmp: %s", dirs)

	log.Println("Calling invalid function...")
	err = client.Invalid(context.Background())
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
