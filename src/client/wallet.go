package client

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/nodeset-org/hyperdrive/shared/types/api"
	nmc_client "github.com/rocket-pool/node-manager-core/api/client"
	nmc_types "github.com/rocket-pool/node-manager-core/api/types"
)

type WalletRequester struct {
	context *nmc_client.RequesterContext
}

func NewWalletRequester(context *nmc_client.RequesterContext) *WalletRequester {
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
func (r *WalletRequester) GetContext() *nmc_client.RequesterContext {
	return r.context
}

// Delete the wallet keystore's password from disk
func (r *WalletRequester) DeletePassword() (*nmc_types.ApiResponse[nmc_types.SuccessData], error) {
	return nmc_client.SendGetRequest[nmc_types.SuccessData](r, "delete-password", "DeletePassword", nil)
}

// Export wallet
func (r *WalletRequester) Export() (*nmc_types.ApiResponse[api.WalletExportData], error) {
	return nmc_client.SendGetRequest[api.WalletExportData](r, "export", "Export", nil)
}

// Export the wallet in encrypted ETH key format
func (r *WalletRequester) ExportEthKey() (*nmc_types.ApiResponse[api.WalletExportEthKeyData], error) {
	return nmc_client.SendGetRequest[api.WalletExportEthKeyData](r, "export-eth-key", "ExportEthKey", nil)
}

// Generate a validator key derived from the node wallet's seed
func (r *WalletRequester) GenerateValidatorKey(path string) (*nmc_types.ApiResponse[api.WalletGenerateValidatorKeyData], error) {
	args := map[string]string{
		"path": path,
	}
	return nmc_client.SendGetRequest[api.WalletGenerateValidatorKeyData](r, "generate-validator-key", "GenerateValidatorKey", args)
}

// Initialize the wallet with a new key
func (r *WalletRequester) Initialize(derivationPath *string, index *uint64, saveWallet bool, password string, savePassword bool) (*nmc_types.ApiResponse[api.WalletInitializeData], error) {
	args := map[string]string{
		"password":      password,
		"save-wallet":   strconv.FormatBool(saveWallet),
		"save-password": strconv.FormatBool(savePassword),
	}
	if derivationPath != nil {
		args["derivation-path"] = *derivationPath
	}
	if index != nil {
		args["index"] = fmt.Sprint(*index)
	}
	return nmc_client.SendGetRequest[api.WalletInitializeData](r, "initialize", "Initialize", args)
}

// Set the node address to an arbitrary address
func (r *WalletRequester) Masquerade(address common.Address) (*nmc_types.ApiResponse[nmc_types.SuccessData], error) {
	args := map[string]string{
		"address": address.Hex(),
	}
	return nmc_client.SendGetRequest[nmc_types.SuccessData](r, "masquerade", "Masquerade", args)
}

// Rebuild the validator keys associated with the wallet
func (r *WalletRequester) Rebuild() (*nmc_types.ApiResponse[api.WalletRebuildData], error) {
	return nmc_client.SendGetRequest[api.WalletRebuildData](r, "rebuild", "Rebuild", nil)
}

// Recover wallet
func (r *WalletRequester) Recover(derivationPath *string, mnemonic *string, index *uint64, password string, save bool) (*nmc_types.ApiResponse[api.WalletRecoverData], error) {
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
	return nmc_client.SendGetRequest[api.WalletRecoverData](r, "recover", "Recover", args)
}

// Set the node address back to the wallet address
func (r *WalletRequester) RestoreAddress() (*nmc_types.ApiResponse[nmc_types.SuccessData], error) {
	return nmc_client.SendGetRequest[nmc_types.SuccessData](r, "restore-address", "RestoreAddress", nil)
}

// Search and recover wallet
func (r *WalletRequester) SearchAndRecover(mnemonic string, address common.Address, password string, save bool) (*nmc_types.ApiResponse[api.WalletSearchAndRecoverData], error) {
	args := map[string]string{
		"mnemonic":      mnemonic,
		"address":       address.Hex(),
		"password":      password,
		"save-password": fmt.Sprint(save),
	}
	return nmc_client.SendGetRequest[api.WalletSearchAndRecoverData](r, "search-and-recover", "SearchAndRecover", args)
}

// Set an ENS reverse record to a name
func (r *WalletRequester) SetEnsName(name string) (*nmc_types.ApiResponse[api.WalletSetEnsNameData], error) {
	args := map[string]string{
		"name": name,
	}
	return nmc_client.SendGetRequest[api.WalletSetEnsNameData](r, "set-ens-name", "SetEnsName", args)
}

// Sets the wallet keystore's password
func (r *WalletRequester) SetPassword(password string, save bool) (*nmc_types.ApiResponse[nmc_types.SuccessData], error) {
	args := map[string]string{
		"password": password,
		"save":     fmt.Sprint(save),
	}
	return nmc_client.SendGetRequest[nmc_types.SuccessData](r, "set-password", "SetPassword", args)
}

// Get wallet status
func (r *WalletRequester) Status() (*nmc_types.ApiResponse[api.WalletStatusData], error) {
	return nmc_client.SendGetRequest[api.WalletStatusData](r, "status", "Status", nil)
}

// Search for and recover the wallet in test-mode so none of the artifacts are saved
func (r *WalletRequester) TestSearchAndRecover(mnemonic string, address common.Address) (*nmc_types.ApiResponse[api.WalletSearchAndRecoverData], error) {
	args := map[string]string{
		"mnemonic": mnemonic,
		"address":  address.Hex(),
	}
	return nmc_client.SendGetRequest[api.WalletSearchAndRecoverData](r, "test-search-and-recover", "TestSearchAndRecover", args)
}

// Recover wallet in test-mode so none of the artifacts are saved
func (r *WalletRequester) TestRecover(derivationPath *string, mnemonic string, index *uint64) (*nmc_types.ApiResponse[api.WalletRecoverData], error) {
	args := map[string]string{
		"mnemonic": mnemonic,
	}
	if derivationPath != nil {
		args["derivation-path"] = *derivationPath
	}
	if index != nil {
		args["index"] = fmt.Sprint(*index)
	}
	return nmc_client.SendGetRequest[api.WalletRecoverData](r, "test-recover", "TestRecover", args)
}

// Sends a zero-value message with a payload
func (r *WalletRequester) SendMessage(message []byte, address common.Address) (*nmc_types.ApiResponse[nmc_types.TxInfoData], error) {
	args := map[string]string{
		"message": hex.EncodeToString(message),
		"address": address.Hex(),
	}
	return nmc_client.SendGetRequest[nmc_types.TxInfoData](r, "send-message", "SendMessage", args)
}

// Use the node private key to sign an arbitrary message
func (r *WalletRequester) SignMessage(message []byte) (*nmc_types.ApiResponse[api.WalletSignMessageData], error) {
	args := map[string]string{
		"message": hex.EncodeToString(message),
	}
	return nmc_client.SendGetRequest[api.WalletSignMessageData](r, "sign-message", "SignMessage", args)
}

// Use the node private key to sign a transaction
func (r *WalletRequester) SignTx(message []byte) (*nmc_types.ApiResponse[api.WalletSignTxData], error) {
	args := map[string]string{
		"tx": hex.EncodeToString(message),
	}
	return nmc_client.SendGetRequest[api.WalletSignTxData](r, "sign-tx", "SignTx", args)
}
