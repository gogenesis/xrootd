package xrootd

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func checkPing(t *testing.B, client *Client, done chan<- bool) {
	err := client.Ping()
	assert.NoError(t, err)
	done <- true
}

func BenchmarkHundredPings(b *testing.B) {
	client, err := New(*Addr)
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

func ExampleClient_Ping() {
	client, _ := New(*Addr)
	client.Ping()
	fmt.Print("Pong!")
	// Output: Pong!
}
