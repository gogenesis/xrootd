package encoder

import (
	"bytes"
	"encoding/binary"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"reflect"
)

// MarshalRequest marshals request body together with request and stream ids
func MarshalRequest(requestID uint16, streamID [2]byte, requestBody interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	_, err := buf.Write(streamID[:])
	if err != nil {
		return nil, err
	}

	err = binary.Write(buf, binary.BigEndian, requestID)
	if err != nil {
		return nil, err
	}

	b, err := Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	return append(buf.Bytes(), b...), nil
}

// Marshal marshals structure to the bytes
func Marshal(x interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	v := reflect.ValueOf(x)
	var err error
	for i := 0; i < v.NumField() && err == nil; i++ {
		switch v.Field(i).Kind() {
		case reflect.Uint8:
			err = buf.WriteByte(v.Field(i).Interface().(uint8))
		case reflect.Uint16:
			err = binary.Write(buf, binary.BigEndian, v.Field(i).Interface().(uint16))
		case reflect.Int32:
			err = binary.Write(buf, binary.BigEndian, v.Field(i).Interface().(int32))
		case reflect.Slice:
			_, err = buf.Write(v.Field(i).Bytes())
		case reflect.Array:
			for j := 0; j < v.Field(i).Len() && err == nil; j++ {
				err = buf.WriteByte(v.Field(i).Index(j).Interface().(uint8))
			}
		default:
			err = errors.Errorf("Cannot encode kind %s", v.Field(i).Kind())
		}
	}
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// UnmarshalFromReader unmarshals data from reader
func UnmarshalFromReader(reader io.Reader, x interface{}) (err error) {
	v := reflect.ValueOf(x)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	for i := 0; i < v.NumField() && err == nil; i++ {
		field := v.Field(i)
		switch field.Kind() {
		case reflect.Uint8:
			var value uint8
			err = binary.Read(reader, binary.BigEndian, &value)
			field.Set(reflect.ValueOf(value))
		case reflect.Uint16:
			var value uint16
			err = binary.Read(reader, binary.BigEndian, &value)
			field.Set(reflect.ValueOf(value))
		case reflect.Int32:
			var value int32
			err = binary.Read(reader, binary.BigEndian, &value)
			field.Set(reflect.ValueOf(value).Convert(field.Type()))
		case reflect.Slice:
			var b []byte
			b, err = ioutil.ReadAll(reader)
			field.SetBytes(b)
		case reflect.Array:
			size := field.Len()
			b := make([]byte, size)
			_, err = io.ReadFull(reader, b)
			for j := 0; j < size; j++ {
				field.Index(j).Set(reflect.ValueOf(b[j]))
			}
		default:
			err = errors.Errorf("Cannot decode kind %s", field.Kind())
		}
	}
	return
}
