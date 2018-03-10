package encoder

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

type request struct {
	X int32
	A uint8
	C uint16
	D [2]byte
	Y []byte
}

type benchmarkRequest struct {
	X   int32
	A   uint8
	A1  uint8
	A2  uint8
	A3  uint8
	A4  int32
	A5  int32
	A6  int32
	A7  int32
	A8  int32
	A9  int32
	A10 int32
	A11 int32
	A12 int32
	A13 int32
	A14 int32
	A16 int32
	C   uint16
	D   [10]byte
	Z   [10]byte
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

func BenchmarkMarshal(b *testing.B) {
	br := benchmarkRequest{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := Marshal(br); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkUnMarshal(b *testing.B) {
	br := benchmarkRequest{}
	var buffer = &bytes.Buffer{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buffer.Write(make([]byte, 78))
		if err := UnmarshalFromReader(buffer, &br); err != nil {
			b.Error(err)
		}
	}
}
