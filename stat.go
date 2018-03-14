package xrootd

import (
	"context"

	"github.com/EgorMatirov/xrootd/requests/stat"
)

// Stat obtains status information for a path
func (client *Client) Stat(ctx context.Context, path string) (*stat.Response, error) {
	serverResponse, err := client.call(ctx, stat.RequestID, stat.NewRequest(path))
	if err != nil {
		return nil, err
	}

	return stat.ParseReponsee(serverResponse)
}
