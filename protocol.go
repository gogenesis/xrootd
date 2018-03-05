package xrootd

import (
	"github.com/EgorMatirov/xrootd/encoder"
	"github.com/EgorMatirov/xrootd/requests/protocol"
)

// Protocol obtains the protocol version number, type of server and security information
func (client *Client) Protocol() (response *protocol.Response, securityInfo *protocol.SecurityInfo, err error) {
	serverResponse, err := client.call(protocol.RequestID, protocol.NewRequest())
	if err != nil {
		return
	}

	response = &protocol.Response{}
	securityInfo = &protocol.SecurityInfo{}

	err = encoder.UnmarshalFromReader(serverResponse, response)
	if err != nil {
		return
	}

	if serverResponse.Len() > 8 {
		err = encoder.UnmarshalFromReader(serverResponse, securityInfo)
	}
	return
}
