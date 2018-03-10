package xrootd

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandshake(t *testing.T) {
	_, err := New(context.Background(), *Addr)
	assert.NoError(t, err)
}
