package xrootd

import (
	"context"

	"github.com/EgorMatirov/xrootd/requests/sync"
)

// Sync commits all pending writes to an open file
func (client *Client) Sync(ctx context.Context, fileHandle [4]byte) error {
	_, err := client.call(ctx, sync.RequestID, sync.NewRequest(fileHandle))
	return err
}
