package handshake

import (
	"bytes"
	"encoding/binary"
)

//go:generate stringer -type=ServerType

type ServerType int

const (
	LoadBalancingServer ServerType = iota
	DataServer          ServerType = iota
	UnknownServer       ServerType = 1000
)

type Handshake struct {
	ProtocolVersion int32
	ServerType      ServerType
}

func GetHandshakeBytes() []byte {
	var handshake = []interface{}{
		int32(0),
		int32(0),
		int32(0),
		int32(4),
		int32(2012),
	}
	buf := new(bytes.Buffer)
	for _, v := range handshake {
		binary.Write(buf, binary.BigEndian, v)
	}
	return buf.Bytes()
}

func ReadHandshake(buffer *bytes.Buffer) (*Handshake, error) {
	var protocolVersion int32
	err := binary.Read(buffer, binary.BigEndian, &protocolVersion)
	if err != nil {
		return nil, err
	}

	var flag int32
	err = binary.Read(buffer, binary.BigEndian, &flag)
	if err != nil {
		return nil, err
	}

	serverType := UnknownServer
	if flag == 0 || flag == 1 {
		serverType = ServerType(flag)
	}

	return &Handshake{protocolVersion, serverType}, nil
}
