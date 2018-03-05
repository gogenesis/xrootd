package xrootd

import (
	"github.com/EgorMatirov/xrootd/encoder"
	"github.com/EgorMatirov/xrootd/requests/handshake"
)

func (client *Client) handshake() error {
	responseChannel, err := client.chm.ClaimWithID([2]byte{0, 0})
	if err != nil {
		return err
	}

	b, err := encoder.Marshal(handshake.NewRequest())
	if err != nil {
		return err
	}
	_, err = client.connection.Write(b)
	if err != nil {
		return err
	}

	var handshakeResult handshake.Response
	err = encoder.UnmarshalFromReader((<-responseChannel).(*serverResponse).Data, &handshakeResult)
	if err != nil {
		return err
	}

	logger.Printf("Connected! Protocol version is %d. Server type is %s.", handshakeResult.ProtocolVersion, handshakeResult.ServerType)

	return nil
}
