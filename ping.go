package xrootd

import (
	"context"

	"github.com/EgorMatirov/xrootd/requests/ping"
)

// Ping determines if the server is still alive
func (client *Client) Ping(ctx context.Context) error {
	_, err := client.call(ctx, ping.RequestID, ping.NewRequest())
	return err
}
