package xrootd

// Invalid sends invalid request. For testing purposes only.
func (client *Client) Invalid() error {
	_, err := client.call(0, struct {
		B [16]byte
		S int32
	}{})
	return err
}
