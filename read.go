package xrootd

import (
	"context"

	"github.com/EgorMatirov/xrootd/encoder"
	"github.com/EgorMatirov/xrootd/requests/read"
)

// Read data from an open file
func (client *Client) Read(ctx context.Context, fileHandle [4]byte, offset int64, length int32) ([]byte, error) {
	serverResponse, err := client.call(ctx, read.RequestID, read.NewRequest(fileHandle, offset, length))
	if err != nil {
		return nil, err
	}

	var result = &read.Response{}
	if err = encoder.Unmarshal(serverResponse, result); err != nil {
		return nil, err
	}

	return result.Data, nil
}
