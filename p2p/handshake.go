package p2p

import "errors"

// Returned if handshake between the local and remote code
// cannot be established
var ErrInvalidHandshake = errors.New("invalid handshake")

type HandshakeFunc func(any) error

func NOPHandshakeFunc(any) error {
	return nil
}
