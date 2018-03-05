package xrootd

import (
	"bytes"
	"github.com/EgorMatirov/xrootd/encoder"
	"github.com/EgorMatirov/xrootd/requests/dirlist"
)

// Dirlist returns contents of a directory
func (client *Client) Dirlist(path string) ([]string, error) {
	serverResponse, err := client.call(dirlist.RequestID, dirlist.NewRequest(path))
	if err != nil {
		return nil, err
	}

	var result = &dirlist.Response{}
	err = encoder.UnmarshalFromReader(serverResponse, result)
	if err != nil {
		return nil, err
	}

	if len(result.Data) == 0 {
		return []string{}, nil
	}

	strings := bytes.Split(result.Data, []byte{'\n'})

	resultStrings := make([]string, len(strings))

	for i := 0; i < len(strings); i++ {
		resultStrings[i] = string(strings[i])
	}

	return resultStrings, nil
}
