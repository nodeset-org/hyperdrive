package api_v2

import (
	"context"
	"fmt"
	"net/http"
)

const (
	// Route for getting a nonce from the NodeSet server
	noncePath = "nonce"
)

// Data used returned from nonce requests
type NonceData struct {
	// The nonce for the session request
	Nonce string `json:"nonce"`

	// The auth token for the session if approved
	Token string `json:"token"`
}

// Get a nonce from the NodeSet server for a new session
func (c *NodeSetClient) Nonce(ctx context.Context) (NonceData, error) {
	// Get the nonce
	code, nonceResponse, err := SubmitRequest[NonceData](c, ctx, false, http.MethodGet, nil, nil, noncePath)
	if err != nil {
		return NonceData{}, fmt.Errorf("error getting nonce: %w", err)
	}

	// Handle response based on return code
	switch code {
	case http.StatusOK:
		return nonceResponse.Data, nil
	}
	return NonceData{}, fmt.Errorf("nodeset server responded to nonce request with code %d: [%s]", code, nonceResponse.Message)
}
