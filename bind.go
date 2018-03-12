package xrootd

import (
	"context"

	"github.com/EgorMatirov/xrootd/encoder"
	"github.com/EgorMatirov/xrootd/requests/bind"
)

// Bind a socket to pre-existing session
func (client *Client) Bind(ctx context.Context, sessionID [16]byte) (byte, error) {
	serverResponse, err := client.call(ctx, bind.RequestID, bind.NewRequest(sessionID))
	if err != nil {
		return 0, err
	}

	var result = &bind.Response{}
	err = encoder.Unmarshal(serverResponse, result)
	if err != nil {
		return 0, err
	}

	return result.PathID, nil
}
