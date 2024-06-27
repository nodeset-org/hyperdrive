package api_v2

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/goccy/go-json"
	"github.com/nodeset-org/hyperdrive-daemon/shared/types"
	"github.com/rocket-pool/node-manager-core/utils"
)

const (
	// Route for getting the latest deposit data set from the NodeSet server
	depositDataPath string = "deposit-data"

	// Subroute for getting the version of the latest deposit data
	metaPath string = "meta"

	// Deposit data has withdrawal creds that don't match a StakeWise vault
	vaultNotFoundKey string = "vault_not_found"

	// Deposit data can't be uploaded to Mainnet because the user isn't allowed to use Mainnet yet
	invalidPermissionsKey string = "invalid_permissions"
)

var (
	// The requested StakeWise vault didn't exist
	ErrVaultNotFound error = errors.New("deposit data has withdrawal creds that don't match a StakeWise vault")

	// The user isn't allowed to use Mainnet yet
	ErrInvalidPermissions error = errors.New("deposit data can't be uploaded to Mainnet because you aren't permitted to use Mainnet yet")
)

// Response to a deposit data request
type DepositDataData struct {
	Version     int                         `json:"version"`
	DepositData []types.ExtendedDepositData `json:"depositData"`
}

// Response to a deposit data meta request
type DepositDataMetaData struct {
	Version int `json:"version"`
}

// Get the aggregated deposit data from the server
func (c *NodeSetClient) DepositData_Get(ctx context.Context, vault common.Address, network string) (DepositDataData, error) {
	// Create the request params
	vaultString := utils.RemovePrefix(strings.ToLower(vault.Hex()))
	params := map[string]string{
		"vault":   vaultString,
		"network": network,
	}

	// Send it
	code, response, err := SubmitRequest[DepositDataData](c, ctx, true, http.MethodGet, nil, params, depositDataPath)
	if err != nil {
		return DepositDataData{}, fmt.Errorf("error getting deposit data: %w", err)
	}

	// Handle response based on return code
	switch code {
	case http.StatusOK:
		return response.Data, nil

	case http.StatusBadRequest:
		switch response.Error {
		case invalidNetworkKey:
			// Network not known
			return DepositDataData{}, ErrInvalidNetwork
		}

	case http.StatusUnauthorized:
		switch response.Error {
		case invalidSessionKey:
			// Invalid or expird session
			return DepositDataData{}, ErrInvalidSession
		}
	}
	return DepositDataData{}, fmt.Errorf("nodeset server responded to deposit-data-get request with code %d: [%s]", code, response.Message)
}

// Get the current version of the aggregated deposit data on the server
func (c *NodeSetClient) DepositDataMeta(ctx context.Context, vault common.Address, network string) (DepositDataMetaData, error) {
	// Create the request params
	vaultString := utils.RemovePrefix(strings.ToLower(vault.Hex()))
	params := map[string]string{
		"vault":   vaultString,
		"network": network,
	}

	// Send it
	code, response, err := SubmitRequest[DepositDataMetaData](c, ctx, true, http.MethodGet, nil, params, depositDataPath, metaPath)
	if err != nil {
		return DepositDataMetaData{}, fmt.Errorf("error getting deposit data version: %w", err)
	}

	// Handle response based on return code
	switch code {
	case http.StatusOK:
		return response.Data, nil

	case http.StatusBadRequest:
		switch response.Error {
		case invalidNetworkKey:
			// Network not known
			return DepositDataMetaData{}, ErrInvalidNetwork
		}

	case http.StatusUnauthorized:
		switch response.Error {
		case invalidSessionKey:
			// Invalid or expird session
			return DepositDataMetaData{}, ErrInvalidSession
		}
	}
	return DepositDataMetaData{}, fmt.Errorf("nodeset server responded to deposit-data-meta request with code %d: [%s]", code, response.Message)
}

// Uploads deposit data to Nodeset
func (c *NodeSetClient) DepositData_Post(ctx context.Context, depositData []*types.ExtendedDepositData) error {
	// Create the request body
	serializedData, err := json.Marshal(depositData)
	if err != nil {
		return fmt.Errorf("error serializing deposit data: %w", err)
	}

	// Send it
	code, response, err := SubmitRequest[struct{}](c, ctx, true, http.MethodPost, bytes.NewBuffer(serializedData), nil, depositDataPath)
	if err != nil {
		return fmt.Errorf("error uploading deposit data: %w", err)
	}

	// Handle response based on return code
	switch code {
	case http.StatusOK:
		return nil

	case http.StatusBadRequest:
		switch response.Error {
		case vaultNotFoundKey:
			// The requested StakeWise vault didn't exist
			return ErrVaultNotFound
		}

	case http.StatusUnauthorized:
		switch response.Error {
		case invalidSessionKey:
			// Invalid or expird session
			return ErrInvalidSession
		}

	case http.StatusForbidden:
		switch response.Error {
		case invalidPermissionsKey:
			// The user isn't allowed to use Mainnet yet
			return ErrInvalidPermissions
		}
	}
	return fmt.Errorf("nodeset server responded to deposit-data-post request with code %d: [%s]", code, response.Message)
}
