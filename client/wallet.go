package client

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/nodeset-org/hyperdrive/shared/types/api"
)

type WalletRequester struct {
	context *RequesterContext
}

func NewWalletRequester(context *RequesterContext) *WalletRequester {
	return &WalletRequester{
		context: context,
	}
}

func (r *WalletRequester) GetName() string {
	return "Wallet"
}
func (r *WalletRequester) GetRoute() string {
	return "wallet"
}
func (r *WalletRequester) GetContext() *RequesterContext {
	return r.context
}

// Delete the wallet keystore's password from disk
func (r *WalletRequester) DeletePassword() (*api.ApiResponse[api.SuccessData], error) {
	return SendGetRequest[api.SuccessData](r, "delete-password", "DeletePassword", nil)
}

// Export wallet
func (r *WalletRequester) Export() (*api.ApiResponse[api.WalletExportData], error) {
	return SendGetRequest[api.WalletExportData](r, "export", "Export", nil)
}

// Export the wallet in encrypted ETH key format
func (r *WalletRequester) ExportEthKey() (*api.ApiResponse[api.WalletExportEthKeyData], error) {
	return SendGetRequest[api.WalletExportEthKeyData](r, "export-eth-key", "ExportEthKey", nil)
}

// Generate a validator key derived from the node wallet's seed
func (r *WalletRequester) GenerateValidatorKey(path string) (*api.ApiResponse[api.WalletGenerateValidatorKeyData], error) {
	args := map[string]string{
		"path": path,
	}
	return SendGetRequest[api.WalletGenerateValidatorKeyData](r, "generate-validator-key", "GenerateValidatorKey", args)
}

// Initialize the wallet with a new key
func (r *WalletRequester) Initialize(derivationPath *string, index *uint64, password string, save bool) (*api.ApiResponse[api.WalletInitializeData], error) {
	args := map[string]string{
		"password":      password,
		"save-password": fmt.Sprint(save),
	}
	if derivationPath != nil {
		args["derivation-path"] = *derivationPath
	}
	if index != nil {
		args["index"] = fmt.Sprint(*index)
	}
	return SendGetRequest[api.WalletInitializeData](r, "initialize", "Initialize", args)
}

// Set the node address to an arbitrary address
func (r *WalletRequester) Masquerade(address common.Address) (*api.ApiResponse[api.SuccessData], error) {
	args := map[string]string{
		"address": address.Hex(),
	}
	return SendGetRequest[api.SuccessData](r, "masquerade", "Masquerade", args)
}

// Rebuild the validator keys associated with the wallet
func (r *WalletRequester) Rebuild() (*api.ApiResponse[api.WalletRebuildData], error) {
	return SendGetRequest[api.WalletRebuildData](r, "rebuild", "Rebuild", nil)
}

// Recover wallet
func (r *WalletRequester) Recover(derivationPath *string, mnemonic *string, index *uint64, password string, save bool) (*api.ApiResponse[api.WalletRecoverData], error) {
	args := map[string]string{
		"password":      password,
		"save-password": fmt.Sprint(save),
	}
	if derivationPath != nil {
		args["derivation-path"] = *derivationPath
	}
	if mnemonic != nil {
		args["mnemonic"] = *mnemonic
	}
	if index != nil {
		args["index"] = fmt.Sprint(*index)
	}
	return SendGetRequest[api.WalletRecoverData](r, "recover", "Recover", args)
}

// Set the node address back to the wallet address
func (r *WalletRequester) RestoreAddress() (*api.ApiResponse[api.SuccessData], error) {
	return SendGetRequest[api.SuccessData](r, "restore-address", "RestoreAddress", nil)
}

// Search and recover wallet
func (r *WalletRequester) SearchAndRecover(mnemonic string, address common.Address, password string, save bool) (*api.ApiResponse[api.WalletSearchAndRecoverData], error) {
	args := map[string]string{
		"mnemonic":      mnemonic,
		"address":       address.Hex(),
		"password":      password,
		"save-password": fmt.Sprint(save),
	}
	return SendGetRequest[api.WalletSearchAndRecoverData](r, "search-and-recover", "SearchAndRecover", args)
}

// Set an ENS reverse record to a name
func (r *WalletRequester) SetEnsName(name string) (*api.ApiResponse[api.WalletSetEnsNameData], error) {
	args := map[string]string{
		"name": name,
	}
	return SendGetRequest[api.WalletSetEnsNameData](r, "set-ens-name", "SetEnsName", args)
}

// Sets the wallet keystore's password
func (r *WalletRequester) SetPassword(password string, save bool) (*api.ApiResponse[api.SuccessData], error) {
	args := map[string]string{
		"password": password,
		"save":     fmt.Sprint(save),
	}
	return SendGetRequest[api.SuccessData](r, "set-password", "SetPassword", args)
}

// Get wallet status
func (r *WalletRequester) Status() (*api.ApiResponse[api.WalletStatusData], error) {
	return SendGetRequest[api.WalletStatusData](r, "status", "Status", nil)
}

// Search for and recover the wallet in test-mode so none of the artifacts are saved
func (r *WalletRequester) TestSearchAndRecover(mnemonic string, address common.Address) (*api.ApiResponse[api.WalletSearchAndRecoverData], error) {
	args := map[string]string{
		"mnemonic": mnemonic,
		"address":  address.Hex(),
	}
	return SendGetRequest[api.WalletSearchAndRecoverData](r, "test-search-and-recover", "TestSearchAndRecover", args)
}

// Recover wallet in test-mode so none of the artifacts are saved
func (r *WalletRequester) TestRecover(derivationPath *string, mnemonic string, index *uint64) (*api.ApiResponse[api.WalletRecoverData], error) {
	args := map[string]string{
		"mnemonic": mnemonic,
	}
	if derivationPath != nil {
		args["derivation-path"] = *derivationPath
	}
	if index != nil {
		args["index"] = fmt.Sprint(*index)
	}
	return SendGetRequest[api.WalletRecoverData](r, "test-recover", "TestRecover", args)
}

// Sends a zero-value message with a payload
func (r *WalletRequester) SendMessage(message []byte, address common.Address) (*api.ApiResponse[api.TxInfoData], error) {
	args := map[string]string{
		"message": hex.EncodeToString(message),
		"address": address.Hex(),
	}
	return SendGetRequest[api.TxInfoData](r, "send-message", "SendMessage", args)
}

// Use the node private key to sign an arbitrary message
func (r *WalletRequester) SignMessage(message []byte) (*api.ApiResponse[api.WalletSignMessageData], error) {
	args := map[string]string{
		"message": hex.EncodeToString(message),
	}
	return SendGetRequest[api.WalletSignMessageData](r, "sign-message", "SignMessage", args)
}
