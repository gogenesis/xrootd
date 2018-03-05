package xrootd

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"os"

	"fmt"
	"github.com/EgorMatirov/xrootd/chanmanager"
	"github.com/EgorMatirov/xrootd/encoder"
)

var logger = log.New(os.Stderr, "xrootd: ", log.LstdFlags)

// A Client to xrootd server
type Client struct {
	connection *net.TCPConn
	chm        *chanmanager.Chanmanager
}

type serverResponse struct {
	Data  *bytes.Buffer
	Error error
}

type serverError struct {
	Code    int32
	Message string
}

func (err serverError) Error() string {
	return fmt.Sprintf("Server error %d: %s", err.Code, err.Message)
}

type responseHeader struct {
	StreamID   [2]byte
	Status     uint16
	DataLength int32
}

// New creates a client to xrootd server at address
func New(address string) (*Client, error) {
	conn, err := createTCPConnection(address)
	if err != nil {
		return nil, err
	}

	client := &Client{conn, chanmanager.New()}

	go client.consume()

	err = client.handshake()
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

		err := encoder.UnmarshalFromReader(client.connection, header)
		if err != nil {
			logger.Panic(err)
		}

		data := make([]byte, header.DataLength)
		err = binary.Read(client.connection, binary.BigEndian, &data)
		if err != nil {
			logger.Panic(err)
		}

		dataBuffer := &bytes.Buffer{}
		_, err = dataBuffer.Write(data)
		if err != nil {
			logger.Panic(err)
		}

		if err == nil && header.Status != 0 {
			err = extractError(header, data)
		}

		err = client.chm.SendData(header.StreamID, &serverResponse{dataBuffer, err})
		if err != nil {
			logger.Panic(err)
		}
		client.chm.Unclaim(header.StreamID)
	}
}

func extractError(header *responseHeader, data []byte) error {
	if header.Status == 4003 {
		errorBuffer := new(bytes.Buffer)
		errorBuffer.Write(data)

		var code int32
		err := binary.Read(errorBuffer, binary.BigEndian, &code)
		if err != nil {
			return err
		}

		message, err := errorBuffer.ReadString(0)
		if err != nil {
			return err
		}
		message = message[:len(message)-1]

		return serverError{code, message}
	}
	return nil
}

func (client *Client) call(requestID uint16, request interface{}) (*bytes.Buffer, error) {
	streamID, responseChannel, err := client.chm.Claim()
	if err != nil {
		return nil, err
	}
	requestData, err := encoder.MarshalRequest(requestID, streamID, request)
	if err != nil {
		return nil, err
	}

	_, err = client.connection.Write(requestData)
	if err != nil {
		return nil, err
	}

	response := (<-responseChannel).(*serverResponse)
	if response.Error != nil {
		return nil, response.Error
	}

	return response.Data, nil
}
