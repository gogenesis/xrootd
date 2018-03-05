package main

import (
	"github.com/EgorMatirov/xrootd"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandshake(t *testing.T) {
	_, err := xrootd.New(*addr)
	assert.NoError(t, err)
}
