package xrootd

import (
	"context"
	"testing"

	"github.com/EgorMatirov/xrootd/requests/open"
	"github.com/stretchr/testify/assert"
)

func TestClient_Write(t *testing.T) {
	client, err := New(context.Background(), *Addr)
	assert.NoError(t, err)

	_, err = client.Login(context.Background(), "gopher")
	assert.NoError(t, err)

	handle, err := client.Open(context.Background(), "/tmp/testFiles/write", open.ModeOwnerWrite, open.OptionsMkPath|open.OptionsDelete)
	assert.NoError(t, err)
	assert.NotNil(t, handle)

	message := []byte("Hello")
	err = client.Write(context.Background(), handle, 0, 0, message)
	assert.NoError(t, err)

	err = client.Sync(context.Background(), handle)
	assert.NoError(t, err)

	err = client.Close(context.Background(), handle, int64(len(message)))
	assert.NoError(t, err)
}
