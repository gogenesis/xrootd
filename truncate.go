package xrootd

import (
	"context"

	"github.com/EgorMatirov/xrootd/requests/truncate"
)

// Truncate a file to a particular size.
func (client *Client) Truncate(ctx context.Context, path string, size int64) error {
	_, err := client.call(ctx, truncate.RequestID, truncate.NewRequestWithPath(path, size))
	return err
}
