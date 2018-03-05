package xrootd

import (
	"github.com/EgorMatirov/xrootd/requests/ping"
)

// Ping determines if the server is still alive
func (client *Client) Ping() error {
	_, err := client.call(ping.RequestID, ping.NewRequest())
	return err
}
