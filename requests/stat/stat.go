package stat

import (
	"bytes"
	"strconv"
	"github.com/pkg/errors"
)

const RequestID uint16 = 3017

type Request struct {
	Options byte
	Reserved [11]byte
	FileHandle [4]byte
	PathLength int32
	Path       []byte
}

type Response struct {
	ID int64
	Size int64
	Flags int32
	ModificationTime int64
}

func NewRequest(path string) Request {
	var pathBytes = make([]byte, len(path))
	copy(pathBytes, path)

	return Request{0, [11]byte{}, [4]byte{}, int32(len(path)), pathBytes}
}

func ParseReponsee(data []byte) (*Response, error){
	dataParts := bytes.Split(data, []byte(" "))
	if len(dataParts) != 4 {
		return nil, errors.Errorf("Not enough fields in stat response: %s", data)
	}
	id, err := strconv.ParseInt(string(dataParts[0]), 10, 64)
	if err != nil {
		return nil, err
	}

	size, err := strconv.ParseInt(string(dataParts[1]), 10, 64)
	if err != nil {
		return nil, err
	}

	flags64, err := strconv.ParseInt(string(dataParts[2]), 10, 32)
	if err != nil {
		return nil, err
	}
	flags := int32(flags64)

	modificationTime, err := strconv.ParseInt(string(dataParts[3][:len(dataParts[3])-1]), 10, 64)
	if err != nil {
		return nil, err
	}

	return &Response{id, size, flags, modificationTime}, nil
}

