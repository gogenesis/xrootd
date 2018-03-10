package xrootd

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func checkPing(t *testing.B, client *Client, done chan<- bool) {
	err := client.Ping(context.Background())
	assert.NoError(t, err)
	done <- true
}

func BenchmarkHundredPings(b *testing.B) {
	client, err := New(context.Background(), *Addr)
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
	client, _ := New(context.Background(), *Addr)
	client.Ping(context.Background())
	fmt.Print("Pong!")
	// Output: Pong!
}
