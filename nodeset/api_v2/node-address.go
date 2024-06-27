package api_v2

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rocket-pool/node-manager-core/log"
	"github.com/rocket-pool/node-manager-core/utils"
)

const (
	// Route for registering a node address with the NodeSet server
	nodeAddressPath string = "node-address"

	// Format for signing node address messages
	nodeAddressMessageFormat string = `{"email":"%s","node_address":"%s"}`

	// The node address has already been confirmed on a NodeSet account
	addressAlreadyAuthorizedKey string = "address_already_authorized"

	// The node address hasn't been whitelisted on the provided NodeSet account
	addressMissingWhitelistKey string = "address_missing_whitelist"
)

var (
	// The node address has already been confirmed on a NodeSet account
	ErrAlreadyRegistered error = errors.New("node has already been registered with the NodeSet server")

	// The node address hasn't been whitelisted on the provided NodeSet account
	ErrNotWhitelisted error = errors.New("node address hasn't been whitelisted on the provided NodeSet account")
)

// Request to register a node with the NodeSet server
type NodeAddressRequest struct {
	// The email address of the NodeSet account
	Email string `json:"email"`

	// The node's wallet address
	NodeAddress string `json:"node_address"`

	// Signature of the request
	Signature string `json:"signature"` // Must be 0x-prefixed hex encoded
}

// Registers the node with the NodeSet server. Assumes wallet validation has already been done and the actual wallet address
// is provided here; if it's not, the signature won't come from the node being registered so it will fail validation.
func (c *NodeSetClient) NodeAddress(ctx context.Context, email string, nodeWallet common.Address) error {
	// Get the logger
	logger, exists := log.FromContext(ctx)
	if !exists {
		panic("context didn't have a logger!")
	}

	// Sign the message
	message := fmt.Sprintf(nodeAddressMessageFormat, email, nodeWallet.Hex())
	sigBytes, err := c.wallet.SignMessage([]byte(message))
	if err != nil {
		return fmt.Errorf("error signing registration message: %w", err)
	}

	// Create the request
	signature := utils.EncodeHexWithPrefix(sigBytes)
	request := NodeAddressRequest{
		Email:       email,
		NodeAddress: nodeWallet.Hex(),
		Signature:   signature,
	}
	jsonData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("error marshalling registration request: %w", err)
	}

	logger.Debug("Sending NodeSet register node request",
		slog.String(log.BodyKey, string(jsonData)),
	)

	// Submit the request
	code, response, err := SubmitRequest[struct{}](c, ctx, false, http.MethodPost, bytes.NewBuffer(jsonData), nil, nodeAddressPath)
	if err != nil {
		return fmt.Errorf("error registering node: %w", err)
	}

	// Handle response based on return code
	switch code {
	case http.StatusOK:
		// Node successfully registered
		return nil

	case http.StatusBadRequest:
		switch response.Error {
		case addressAlreadyAuthorizedKey:
			// Already registered
			return ErrAlreadyRegistered

		case addressMissingWhitelistKey:
			// Not whitelisted in the user account
			return ErrNotWhitelisted

		case invalidSignatureKey:
			// Invalid signature
			return ErrInvalidSignature

		case malformedInputKey:
			// Malformed input
			return ErrMalformedInput
		}
	}
	return fmt.Errorf("nodeset server responded to node-address request with code %d: [%s]", code, response.Message)
}
