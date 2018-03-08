package xrootd

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLogin(t *testing.T) {
	client, err := New(*Addr)
	assert.NoError(t, err)

	_, err = client.Login("gopher")
	assert.NoError(t, err)
}

func ExampleClient_Login() {
	client, _ := New(*Addr)
	loginResult, _ := client.Login("gopher")
	fmt.Printf("Logged in! Security information length is %d. Value is \"%s\"\n", len(loginResult.SecurityInformation), loginResult.SecurityInformation)
	// Output: Logged in! Security information length is 0. Value is ""
}
