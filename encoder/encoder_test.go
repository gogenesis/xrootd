package encoder

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

type request struct {
	X int32
	A uint8
	C uint16
	D [2]byte
	Y []byte
}

type undecodable struct {
	A float64
}

func TestEncode(t *testing.T) {
	var requestID uint16 = 1337
	var streamID = [2]byte{42, 37}
	expected := []byte{42, 37, 5, 57, 0, 0, 0, 1, 2, 0, 3, 6, 7, 11, 13}

	actual, err := MarshalRequest(requestID, streamID, request{1, 2, 3, [2]byte{6, 7}, []byte{11, 13}})

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestDecode(t *testing.T) {
	var expected = request{1, 2, 3, [2]byte{6, 7}, []byte{11, 13}}

	var buffer = &bytes.Buffer{}
	buffer.Write([]byte{0, 0, 0, 1, 2, 0, 3, 6, 7, 11, 13})

	var actual = &request{}
	err := UnmarshalFromReader(buffer, actual)

	assert.Equal(t, expected, *actual)
	assert.NoError(t, err)
}

func TestEncodeUndecodable(t *testing.T) {
	var requestID uint16 = 1337
	var streamID = [2]byte{42, 37}

	_, err := MarshalRequest(requestID, streamID, undecodable{1})

	assert.Error(t, err)
}

func TestDecodeUndecodable(t *testing.T) {
	var buffer = &bytes.Buffer{}
	buffer.Write([]byte{0, 0, 0, 1, 2, 0, 3, 6, 7, 11, 13})

	var actual = &undecodable{}
	err := UnmarshalFromReader(buffer, actual)

	assert.Error(t, err)
}
