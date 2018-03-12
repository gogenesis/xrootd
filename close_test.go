package xrootd

import (
	"context"
	"testing"

	"github.com/EgorMatirov/xrootd/requests/open"
	"github.com/stretchr/testify/assert"
)

func TestClient_Close(t *testing.T) {
	client, err := New(context.Background(), *Addr)
	assert.NoError(t, err)

	_, err = client.Login(context.Background(), "gopher")
	assert.NoError(t, err)

	handle, err := client.Open(context.Background(), "/tmp/testFiles/close", open.ModeOwnerWrite, open.OptionsMkPath|open.OptionsDelete)
	assert.NoError(t, err)
	assert.NotNil(t, handle)

	err = client.Close(context.Background(), handle, 0)
	assert.NoError(t, err)
}
