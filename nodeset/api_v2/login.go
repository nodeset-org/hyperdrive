package api_v2

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rocket-pool/node-manager-core/utils"
)

const (
	// Format for signing login messages
	LoginMessageFormat string = `{"nonce":"%s","address":"%s"}`

	// Route for logging into the NodeSet server
	loginPath string = "login"

	// The provided nonce didn't match an expected one
	invalidNonceKey string = "invalid_nonce"

	// Value of the auth response header if the node hasn't registered yet
	unregisteredAddressKey string = "unregistered_address"
)

var (
	// The provided nonce didn't match an expected one
	ErrInvalidNonce error = errors.New("invalid nonce provided for login")

	// The node hasn't been registered with the NodeSet server yet
	ErrUnregisteredNode error = errors.New("node hasn't been registered with the NodeSet server yet")
)

// Request to log into the NodeSet server
type LoginRequest struct {
	// The nonce for the session request
	Nonce string `json:"nonce"`

	// The node's wallet address
	Address string `json:"address"`

	// Signature of the login request
	Signature string `json:"signature"` // Must be 0x-prefixed hex encoded
}

// Response to a login request
type LoginData struct {
	// The auth token for the session if approved
	Token string `json:"token"`
}

// Logs into the NodeSet server, starting a new session
func (c *NodeSetClient) Login(ctx context.Context, nonce string, address common.Address, signature []byte) (LoginData, error) {
	// Create the request body
	addressString := address.Hex()
	signatureString := utils.EncodeHexWithPrefix(signature)
	request := LoginRequest{
		Nonce:     nonce,
		Address:   addressString,
		Signature: signatureString,
	}
	jsonData, err := json.Marshal(request)
	if err != nil {
		return LoginData{}, fmt.Errorf("error marshalling login request: %w", err)
	}

	// Submit the request
	code, response, err := SubmitRequest[LoginData](c, ctx, true, http.MethodPost, bytes.NewBuffer(jsonData), nil, loginPath)
	if err != nil {
		return LoginData{}, fmt.Errorf("error submitting login request: %w", err)
	}

	// Handle response based on return code
	switch code {
	case http.StatusOK:
		// Login successful, session established
		return response.Data, nil

	case http.StatusBadRequest:
		switch response.Error {
		case invalidSignatureKey:
			// Invalid signature
			return LoginData{}, ErrInvalidSignature

		case malformedInputKey:
			// Malformed input
			return LoginData{}, ErrMalformedInput

		case invalidNonceKey:
			// Invalid nonce
			return LoginData{}, ErrInvalidNonce
		}

	case http.StatusUnauthorized:
		switch response.Error {
		case unregisteredAddressKey:
			// Node hasn't been registered yet
			return LoginData{}, ErrUnregisteredNode

		case invalidSessionKey:
			// The nonce wasn't expected?
			return LoginData{}, ErrInvalidSession
		}
	}
	return LoginData{}, fmt.Errorf("nodeset server responded to login request with code %d: [%s]", code, response.Message)
}
