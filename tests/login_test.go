package main

import (
	"github.com/EgorMatirov/xrootd"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLogin(t *testing.T) {
	client, err := xrootd.New(*addr)
	assert.NoError(t, err)

	_, err = client.Login("gopher")
	assert.NoError(t, err)
}
