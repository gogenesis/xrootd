package xrootd

import (
	"context"

	"github.com/EgorMatirov/xrootd/encoder"
	"github.com/EgorMatirov/xrootd/requests/login"
)

// Login initializes a server connection using username
func (client *Client) Login(ctx context.Context, username string) (*login.Response, error) {
	serverResponse, err := client.call(ctx, login.RequestID, login.NewRequest(username))
	if err != nil {
		return nil, err
	}

	var response = &login.Response{}
	err = encoder.Unmarshal(serverResponse, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
