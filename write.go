package xrootd

import (
	"context"

	"github.com/EgorMatirov/xrootd/requests/write"
)

// Write data to an open file
func (client *Client) Write(ctx context.Context, fileHandle [4]byte, offset int64, pathID byte, data []byte) error {
	_, err := client.call(ctx, write.RequestID, write.NewRequest(fileHandle, offset, pathID, data))
	return err
}
