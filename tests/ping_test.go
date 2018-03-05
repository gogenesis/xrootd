package main

import (
	"github.com/EgorMatirov/xrootd"
	"github.com/stretchr/testify/assert"
	"testing"
)

func checkPing(t *testing.B, client *xrootd.Client, done chan<- bool) {
	err := client.Ping()
	assert.NoError(t, err)
	done <- true
}

func BenchmarkHundredPings(b *testing.B) {
	client, err := xrootd.New(*addr)
	assert.NoError(b, err)

	count := 100
	done := make(chan bool, count)
	b.ResetTimer()

	for i := 0; i < count; i++ {
		go checkPing(b, client, done)
	}

	for i := 0; i < count; i++ {
		<-done
	}
}
