package xrootd

import (
	"github.com/EgorMatirov/xrootd/encoder"
	"github.com/EgorMatirov/xrootd/requests/login"
)

// Login initializes a server connection using username
func (client *Client) Login(username string) (*login.Response, error) {
	serverResponse, err := client.call(login.RequestID, login.NewRequest(username))
	if err != nil {
		return nil, err
	}

	var response = &login.Response{}
	err = encoder.UnmarshalFromReader(serverResponse, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
