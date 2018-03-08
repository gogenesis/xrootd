package xrootd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandshake(t *testing.T) {
	_, err := New(*Addr)
	assert.NoError(t, err)
}
