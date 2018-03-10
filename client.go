package xrootd

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/EgorMatirov/xrootd/chanmanager"
	"github.com/EgorMatirov/xrootd/encoder"
)

var logger = log.New(os.Stderr, "xrootd: ", log.LstdFlags)

// A Client to xrootd server
type Client struct {
	connection      *net.TCPConn
	chm             *chanmanager.Chanmanager
	protocolVersion int32
}

type serverResponse struct {
	Data  []byte
	Error error
}

type serverError struct {
	Code    int32
	Message string
}

func (err serverError) Error() string {
	return fmt.Sprintf("Server error %d: %s", err.Code, err.Message)
}

const responseHeaderSize = 2 + 2 + 4

type responseHeader struct {
	StreamID   [2]byte
	Status     uint16
	DataLength int32
}

// New creates a client to xrootd server at address
func New(ctx context.Context, address string) (*Client, error) {
	conn, err := createTCPConnection(address)
	if err != nil {
		return nil, err
	}

	client := &Client{conn, chanmanager.New(), 0}

	go client.consume()

	err = client.handshake(ctx)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func createTCPConnection(address string) (connection *net.TCPConn, err error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil {
		return
	}

	connection, err = net.DialTCP("tcp", nil, tcpAddr)
	return
}

func (client *Client) consume() {
	for {
		var header = &responseHeader{}

		var headerBytes = make([]byte, responseHeaderSize)
		if _, err := io.ReadFull(client.connection, headerBytes); err != nil {
			logger.Panic(err)
		}

		if err := encoder.Unmarshal(headerBytes, header); err != nil {
			logger.Panic(err)
		}

		data := make([]byte, header.DataLength)
		if _, err := io.ReadFull(client.connection, data); err != nil {
			logger.Panic(err)
		}

		response := &serverResponse{data, nil}
		if header.Status != 0 {
			response.Error = extractError(header, data)
		}

		if err := client.chm.SendData(header.StreamID, response); err != nil {
			logger.Panic(err)
		}

		client.chm.Unclaim(header.StreamID)
	}
}

func extractError(header *responseHeader, data []byte) error {
	if header.Status == 4003 {
		code := int32(binary.BigEndian.Uint32(data[0:4]))
		message := string(data[4 : len(data)-1]) // Skip \0 character at the end

		return serverError{code, message}
	}
	return nil
}

func (client *Client) callWithBytesAndResponseChannel(ctx context.Context, responseChannel <-chan interface{}, requestData []byte) (responseBytes []byte, err error) {
	if _, err = client.connection.Write(requestData); err != nil {
		return nil, err
	}

	select {
	case response := <-responseChannel:
		serverResponse := response.(*serverResponse)
		responseBytes = serverResponse.Data
		err = serverResponse.Error
	case <-ctx.Done():
		err = ctx.Err()
	}

	return
}

func (client *Client) call(ctx context.Context, requestID uint16, request interface{}) (responseBytes []byte, err error) {
	streamID, responseChannel, err := client.chm.Claim()
	if err != nil {
		return nil, err
	}

	requestData, err := encoder.MarshalRequest(requestID, streamID, request)
	if err != nil {
		return nil, err
	}

	return client.callWithBytesAndResponseChannel(ctx, responseChannel, requestData)
}
