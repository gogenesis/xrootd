package main

import (
	"flag"
	"fmt"
	"github.com/EgorMatirov/xrootd"
	"log"
	"os"
)

var addr = flag.String("addr", "0.0.0.0:9001", "address of xrootd server")

func main() {
	flag.Parse()

	client, err := xrootd.New(*addr)
	checkError(err)

	response, securityInfo, err := client.Protocol()
	checkError(err)
	log.Printf("Protocol binary version is %d. Security level is %d.", response.BinaryProtocolVersion, securityInfo.SecurityLevel)

	loginResult, err := client.Login("gopher")
	checkError(err)
	log.Printf("Logged in! Security information length is %d. Value is \"%s\"\n", len(loginResult.SecurityInformation), loginResult.SecurityInformation)

	err = client.Ping()
	checkError(err)
	log.Print("Pong!")

	dirs, err := client.Dirlist("/tmp/")
	checkError(err)
	log.Printf("dirlist /tmp: %s", dirs)

	log.Println("Calling invalid function...")
	err = client.Invalid()
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
