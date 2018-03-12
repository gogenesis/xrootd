package xrootd

import (
	"context"

	"github.com/EgorMatirov/xrootd/requests/close"
)

// Close a previously opened file by handle
func (client *Client) Close(ctx context.Context, fileHandle [4]byte, fileSize int64) error {
	_, err := client.call(ctx, close.RequestID, close.NewRequest(fileHandle, fileSize))
	return err
}
