package xrootd

import "context"

// Invalid sends invalid request. For testing purposes only.
func (client *Client) Invalid(ctx context.Context) error {
	_, err := client.call(ctx, 0, struct {
		B [16]byte
		S int32
	}{})
	return err
}
