package xrootd

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"os"

	"github.com/EgorMatirov/xrootd/requests/handshake"
)

var logger = log.New(os.Stderr, "xrootd: ", log.LstdFlags)

// A Client to xrootd server
type Client struct {
	connection      *net.TCPConn
	responseWaiters map[[2]byte]chan<- serverResponse
}

type serverResponse struct {
	Data []byte
}

func (client *Client) consume() {
	for {
		var streamID [2]byte
		err := binary.Read(client.connection, binary.BigEndian, &streamID)
		if err != nil {
			logger.Panic(err)
		}

		var status uint16
		err = binary.Read(client.connection, binary.BigEndian, &status)
		if err != nil {
			logger.Panic(err)
		}

		var dataLength int32
		err = binary.Read(client.connection, binary.BigEndian, &dataLength)
		if err != nil {
			logger.Panic(err)
		}

		data := make([]byte, dataLength)
		err = binary.Read(client.connection, binary.BigEndian, &data)
		if err != nil {
			logger.Panic(err)
		}

		client.responseWaiters[streamID] <- serverResponse{data}
		delete(client.responseWaiters, streamID)
	}
}

func (client *Client) createResponseChannelWithID(streamID [2]byte) <-chan serverResponse {
	channel := make(chan serverResponse, 1)
	client.responseWaiters[streamID] = channel
	return channel
}

func (client *Client) handshake() error {
	responseChannel := client.createResponseChannelWithID([2]byte{0, 0})

	_, err := client.connection.Write(handshake.GetHandshakeBytes())
	if err != nil {
		return err
	}

	handshakeResponse := <-responseChannel

	handshakeResponseBuffer := new(bytes.Buffer)
	_, err = handshakeResponseBuffer.Write(handshakeResponse.Data)
	if err != nil {
		return err
	}

	handshakeResult, err := handshake.ReadHandshake(handshakeResponseBuffer)
	if err != nil {
		return err
	}

	logger.Printf("Connected! Protocol version is %d. Server type is %s.", handshakeResult.ProtocolVersion, handshakeResult.ServerType)

	return nil
}

// Connect creates a client to xrootd server at address
func Connect(address string) (*Client, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}

	client := &Client{conn, make(map[[2]byte]chan<- serverResponse)}

	go client.consume()

	err = client.handshake()
	if err != nil {
		return nil, err
	}

	return client, nil
}
