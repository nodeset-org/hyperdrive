package nodeset

import "github.com/ethereum/go-ethereum/common"

// Represents a session with the NodeSet service
type Session struct {
	// The nonce used to represent the session during establishment
	Nonce string

	// The session's token
	Token string

	// The address the session was established for
	Address common.Address
}
