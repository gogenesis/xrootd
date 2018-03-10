package xrootd

import (
	"context"

	"github.com/EgorMatirov/xrootd/encoder"
	"github.com/EgorMatirov/xrootd/requests/handshake"
)

func (client *Client) handshake(ctx context.Context) error {
	responseChannel, err := client.chm.ClaimWithID([2]byte{0, 0})
	if err != nil {
		return err
	}

	requestBytes, err := encoder.Marshal(handshake.NewRequest())
	if err != nil {
		return err
	}

	responseBytes, err := client.callWithBytesAndResponseChannel(ctx, responseChannel, requestBytes)
	if err != nil {
		return err
	}

	var handshakeResult handshake.Response
	if encoder.Unmarshal(responseBytes, &handshakeResult) != nil {
		return err
	}

	client.protocolVersion = handshakeResult.ProtocolVersion
	logger.Printf("Connected! Protocol version is %d. Server type is %s.", handshakeResult.ProtocolVersion, handshakeResult.ServerType)

	return nil
}
