package xrootd

import (
	"context"

	"github.com/EgorMatirov/xrootd/encoder"
	"github.com/EgorMatirov/xrootd/requests/protocol"
)

// Protocol obtains the protocol version number, type of server and security information
func (client *Client) Protocol(ctx context.Context) (response *protocol.Response, securityInfo *protocol.SecurityInfo, err error) {
	serverResponse, err := client.call(ctx, protocol.RequestID, protocol.NewRequest())
	if err != nil {
		return
	}

	response = &protocol.Response{}
	securityInfo = &protocol.SecurityInfo{}

	err = encoder.Unmarshal(serverResponse, response)
	if err != nil {
		return
	}

	if len(serverResponse) > 8 {
		err = encoder.Unmarshal(serverResponse, securityInfo)
	}
	return
}
