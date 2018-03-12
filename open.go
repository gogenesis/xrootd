package xrootd

import (
	"context"

	"github.com/EgorMatirov/xrootd/encoder"
	"github.com/EgorMatirov/xrootd/requests/open"
)

// Open returns file handle for a file.
func (client *Client) Open(ctx context.Context, path string, mode open.Mode, options open.Options) ([4]byte, error) {
	serverResponse, err := client.call(ctx, open.RequestID, open.NewRequest(path, mode, options))
	if err != nil {
		return [4]byte{}, err
	}

	var result = &open.Response{}
	err = encoder.Unmarshal(serverResponse, result)
	if err != nil {
		return [4]byte{}, err
	}

	return result.FileHandle, nil
}
