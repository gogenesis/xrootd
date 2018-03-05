## What is it?
A client for [xrootd](http://xrootd.org/).

## Usage example
A simple example which only connects to the server (see tests/main.go).
```go
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
	_, err := xrootd.New(os.Args[1])

	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}
```

Expected output is:
~~~
xrootd: 2018/03/05 08:50:00 Connected! Protocol version is 784. Server type is DataServer.
~~~

# Architecture
## Request flow
1. register a `channel` in `Client.responseWaiters`.
2. send a request to the server by single call to the `net.TCPConn.Write()` - this is thread-safe so no locking is required.
3. await for the response from the `channel`.
4. parse response data and return it to the caller.

## Response flow
Responses can came in any order so we use a map of channels to provide response back to the sender.
1. retrieve `streamID` - we'll use it to find the  `channel`.
2. retrieve `status` - check if error occurred or any follow-up request \ response reading is needed.
3. read `rlen` - the length of response data.
4. read response data and pass it to the sender via specific `channel` from `Client.responseWaiters`.

# Supported requests:
## Protocol
```go
response, securityInfo, _ := client.Protocol()
log.Printf("Protocol binary version is %d. Security level is %d.", response.BinaryProtocolVersion, securityInfo.SecurityLevel)

```

## Login
```go
loginResult, _ := client.Login("gopher")
log.Printf("Logged in! Security information length is %d. Value is \"%s\"\n", len(loginResult.SecurityInformation), loginResult.SecurityInformation)
```

## Ping
```go
err = client.Ping()
if err == nil {
    log.Print("Pong!")
}
```

## Dirlist
```go
dirs, _ := client.Dirlist("/tmp/")
log.Printf("dirlist /tmp: %s", dirs)
```