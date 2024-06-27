package api_v2

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/goccy/go-json"
	"github.com/rocket-pool/node-manager-core/beacon"
)

const (
	// Route for interacting with the list of validators
	validatorsPath string = "validators"

	// The requester doesn't own the provided validator
	invalidValidatorOwnerKey string = "invalid_validator_owner"

	// The exit message provided was invalid
	invalidExitMessage string = "invalid_exit_message"
)

var (
	// The requester doesn't own the provided validator
	ErrInvalidValidatorOwner error = fmt.Errorf("this node doesn't own one of the provided validators")

	// The exit message provided was invalid
	ErrInvalidExitMessage error = fmt.Errorf("the provided exit message was invalid")
)

// Details of an exit message
type ExitMessageDetails struct {
	Epoch          string `json:"epoch"`
	ValidatorIndex string `json:"validator_index"`
}

// Voluntary exit message
type ExitMessage struct {
	Message   ExitMessageDetails `json:"message"`
	Signature string             `json:"signature"`
}

// Data for a pubkey's voluntary exit message
type ExitData struct {
	Pubkey      string      `json:"pubkey"`
	ExitMessage ExitMessage `json:"exit_message"`
}

// Validator status info
type ValidatorStatus struct {
	Pubkey              beacon.ValidatorPubkey `json:"pubkey"`
	Status              string                 `json:"status"`
	ExitMessageUploaded bool                   `json:"exitMessage"`
}

// Response to a validators request
type ValidatorsData struct {
	Validators []ValidatorStatus `json:"validators"`
}

// Get a list of all of the pubkeys that have already been registered with NodeSet for this node
func (c *NodeSetClient) GetRegisteredValidators(ctx context.Context) (ValidatorsData, error) {
	queryParams := map[string]string{
		"network": c.networkName,
	}
	code, response, err := SubmitRequest[ValidatorsData](c, ctx, true, http.MethodGet, nil, queryParams, validatorsPath)
	if err != nil {
		return ValidatorsData{}, fmt.Errorf("error getting registered validators: %w", err)
	}

	// Handle response based on return code
	switch code {
	case http.StatusOK:
		return response.Data, nil

	case http.StatusBadRequest:
		switch response.Error {
		case invalidNetworkKey:
			// Network not known
			return ValidatorsData{}, ErrInvalidNetwork
		}

	case http.StatusUnauthorized:
		switch response.Error {
		case invalidSessionKey:
			// Invalid or expird session
			return ValidatorsData{}, ErrInvalidSession
		}
	}
	return ValidatorsData{}, fmt.Errorf("nodeset server responded to request with code %d: [%s]", code, response.Message)
}

// Submit signed exit data to Nodeset
func (c *NodeSetClient) UploadSignedExitData(ctx context.Context, exitData []ExitData) error {
	// Serialize the exit data into JSON
	jsonData, err := json.Marshal(exitData)
	if err != nil {
		return fmt.Errorf("error marshalling exit data to JSON: %w", err)
	}
	params := map[string]string{
		"network": c.networkName,
	}

	// Submit the PATCH request with the serialized JSON
	code, response, err := SubmitRequest[struct{}](c, ctx, true, http.MethodPatch, bytes.NewBuffer(jsonData), params, validatorsPath)
	if err != nil {
		return fmt.Errorf("error submitting exit data: %w", err)
	}

	// Handle response based on return code
	switch code {
	case http.StatusOK:
		return nil

	case http.StatusBadRequest:
		switch response.Error {
		case invalidNetworkKey:
			// Network not known
			return ErrInvalidNetwork

		case malformedInputKey:
			// Invalid input
			return ErrMalformedInput

		case invalidValidatorOwnerKey:
			// Invalid validator owner
			return ErrInvalidValidatorOwner

		case invalidExitMessage:
			// Invalid exit message
			return ErrInvalidExitMessage
		}

	case http.StatusUnauthorized:
		switch response.Error {
		case invalidSessionKey:
			// Invalid or expird session
			return ErrInvalidSession
		}
	}
	return fmt.Errorf("nodeset server responded to request with code %d: [%s]", code, response.Message)
}
