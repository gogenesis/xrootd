package xrootd

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	client, err := New(context.Background(), *Addr)
	assert.NoError(t, err)

	_, err = client.Login(context.Background(), "gopher")
	assert.NoError(t, err)
}

func ExampleClient_Login() {
	client, _ := New(context.Background(), *Addr)
	loginResult, _ := client.Login(context.Background(), "gopher")
	fmt.Printf("Logged in! Security information length is %d. Value is \"%s\"\n", len(loginResult.SecurityInformation), loginResult.SecurityInformation)
	// Output: Logged in! Security information length is 0. Value is ""
}
