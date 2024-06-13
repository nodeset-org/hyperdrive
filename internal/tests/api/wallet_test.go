package client

import (
	"math/big"
	"os"
	"runtime/debug"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/nodeset-org/osha"
	"github.com/nodeset-org/osha/keys"
	"github.com/rocket-pool/node-manager-core/eth"
	"github.com/rocket-pool/node-manager-core/wallet"
	"github.com/stretchr/testify/require"
)

const (
	expectedWalletAddressString string  = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
	goodPassword                string  = "some_password123"
	expectedBalanceFloat        float64 = 10000
)

var (
	emptyWalletAddress    common.Address = common.HexToAddress("0x0000000000000000000000000000000000000000")
	expectedWalletAddress common.Address = common.HexToAddress(expectedWalletAddressString)
	expectedBalance       *big.Int       = eth.EthToWei(expectedBalanceFloat)
)

func TestWalletRecover_Success(t *testing.T) {
	// Take a snapshot, revert at the end
	snapshotName, err := testMgr.CreateCustomSnapshot(osha.Service_Filesystem)
	if err != nil {
		fail("Error creating custom snapshot: %v", err)
	}
	defer wallet_cleanup(snapshotName)

	dataDir := testMgr.ServiceProvider.GetConfig().UserDataPath.Value

	// Make sure the data directory exists
	err = os.MkdirAll(dataDir, 0755)
	if err != nil {
		t.Fatalf("Error creating data directory: %v", err)
	}

	// Run the round-trip test
	derivationPath := string(wallet.DerivationPath_Default)
	index := uint64(0)
	response, err := apiClient.Wallet.Recover(&derivationPath, keys.DefaultMnemonic, &index, goodPassword, true)
	require.NoError(t, err)
	t.Log("Recover called")

	// Check the response
	require.Equal(t, expectedWalletAddress, response.Data.AccountAddress)
	t.Log("Received correct wallet address")
}

func TestWalletRecover_WrongIndex(t *testing.T) {
	// Take a snapshot, revert at the end
	snapshotName, err := testMgr.CreateCustomSnapshot(osha.Service_Filesystem)
	if err != nil {
		fail("Error creating custom snapshot: %v", err)
	}
	defer wallet_cleanup(snapshotName)

	dataDir := testMgr.ServiceProvider.GetConfig().UserDataPath.Value

	// Make sure the data directory exists
	err = os.MkdirAll(dataDir, 0755)
	if err != nil {
		t.Fatalf("Error creating data directory: %v", err)
	}

	// Run the round-trip test
	derivationPath := string(wallet.DerivationPath_Default)
	index := uint64(1)
	response, err := apiClient.Wallet.Recover(&derivationPath, keys.DefaultMnemonic, &index, goodPassword, true)
	require.NoError(t, err)
	t.Log("Recover called")

	// Check the response
	require.NotEqual(t, expectedWalletAddress, response.Data.AccountAddress)
	t.Logf("Wallet address doesn't match as expected (expected %s, got %s)", expectedWalletAddress.Hex(), response.Data.AccountAddress.Hex())
}

func TestWalletRecover_WrongDerivationPath(t *testing.T) {
	snapshotName, err := testMgr.CreateCustomSnapshot(osha.Service_Filesystem)
	if err != nil {
		fail("Error creating custom snapshot: %v", err)
	}
	defer wallet_cleanup(snapshotName)

	dataDir := testMgr.ServiceProvider.GetConfig().UserDataPath.Value

	// Make sure the data directory exists
	err = os.MkdirAll(dataDir, 0755)
	if err != nil {
		t.Fatalf("Error creating data directory: %v", err)
	}

	// Run the round-trip test
	derivationPath := string(wallet.DerivationPath_LedgerLive)
	index := uint64(0)
	response, err := apiClient.Wallet.Recover(&derivationPath, keys.DefaultMnemonic, &index, goodPassword, true)
	require.NoError(t, err)
	t.Log("Recover called")

	// Check the response
	require.NotEqual(t, expectedWalletAddress, response.Data.AccountAddress)
	t.Logf("Wallet address doesn't match as expected (expected %s, got %s)", expectedWalletAddress.Hex(), response.Data.AccountAddress.Hex())
}

func TestWalletStatus_NotLoaded(t *testing.T) {
	response, err := apiClient.Wallet.Status()
	t.Log("Status called")

	require.NoError(t, err)
	require.Equal(t, response.Data.WalletStatus.Address.NodeAddress, emptyWalletAddress)
	require.False(t, response.Data.WalletStatus.Address.HasAddress)

	require.Equal(t, response.Data.WalletStatus.Wallet.Type, wallet.WalletType(""))
	require.False(t, response.Data.WalletStatus.Wallet.IsLoaded)
	require.False(t, response.Data.WalletStatus.Wallet.IsOnDisk)
	require.Equal(t, response.Data.WalletStatus.Wallet.WalletAddress, emptyWalletAddress)

	t.Log("Received correct wallet status")
}

func TestWalletStatus_Loaded(t *testing.T) {
	// Take a snapshot, revert at the end
	snapshotName, err := testMgr.CreateCustomSnapshot(osha.Service_EthClients | osha.Service_Filesystem)
	if err != nil {
		fail("Error creating custom snapshot: %v", err)
	}
	defer wallet_cleanup(snapshotName)

	// Commit a block just so the latest block is fresh - otherwise the sync progress check will
	// error out because the block is too old and it thinks the client just can't find any peers
	err = testMgr.CommitBlock()
	if err != nil {
		t.Fatalf("Error committing block: %v", err)
	}

	// Make sure the data directory exists
	dataDir := testMgr.ServiceProvider.GetConfig().UserDataPath.Value
	err = os.MkdirAll(dataDir, 0755)
	if err != nil {
		t.Fatalf("Error creating data directory: %v", err)
	}

	// Regen the wallet
	derivationPath := string(wallet.DerivationPath_Default)
	index := uint64(0)
	_, err = apiClient.Wallet.Recover(&derivationPath, keys.DefaultMnemonic, &index, goodPassword, true)
	require.NoError(t, err)
	t.Log("Recover called")

	response, err := apiClient.Wallet.Status()
	t.Log("Status called")

	require.NoError(t, err)
	require.Equal(t, response.Data.WalletStatus.Address.NodeAddress, expectedWalletAddress)
	require.True(t, response.Data.WalletStatus.Address.HasAddress)

	require.Equal(t, response.Data.WalletStatus.Wallet.Type, wallet.WalletType("local"))
	require.True(t, response.Data.WalletStatus.Wallet.IsLoaded)
	require.True(t, response.Data.WalletStatus.Wallet.IsOnDisk)
	require.Equal(t, response.Data.WalletStatus.Wallet.WalletAddress, expectedWalletAddress)

	t.Log("Received correct wallet status")
}

func TestWalletBalance(t *testing.T) {
	// Take a snapshot, revert at the end
	snapshotName, err := testMgr.CreateCustomSnapshot(osha.Service_EthClients | osha.Service_Filesystem)
	if err != nil {
		fail("Error creating custom snapshot: %v", err)
	}
	defer wallet_cleanup(snapshotName)

	// Commit a block just so the latest block is fresh - otherwise the sync progress check will
	// error out because the block is too old and it thinks the client just can't find any peers
	err = testMgr.CommitBlock()
	if err != nil {
		t.Fatalf("Error committing block: %v", err)
	}

	// Make sure the data directory exists
	dataDir := testMgr.ServiceProvider.GetConfig().UserDataPath.Value
	err = os.MkdirAll(dataDir, 0755)
	if err != nil {
		t.Fatalf("Error creating data directory: %v", err)
	}

	// Regen the wallet
	derivationPath := string(wallet.DerivationPath_Default)
	index := uint64(0)
	_, err = apiClient.Wallet.Recover(&derivationPath, keys.DefaultMnemonic, &index, goodPassword, true)
	require.NoError(t, err)
	t.Log("Recover called")

	// Run the round-trip test
	response, err := apiClient.Wallet.Balance()
	require.NoError(t, err)
	t.Log("Balance called")

	// Check the response
	require.Equal(t, 0, response.Data.Balance.Cmp(expectedBalance))
	t.Logf("Received correct balance (%s)", response.Data.Balance.String())
}

// Clean up after each test
func wallet_cleanup(snapshotName string) {
	// Handle panics
	r := recover()
	if r != nil {
		debug.PrintStack()
		fail("Recovered from panic: %v", r)
	}

	// Revert to the snapshot taken at the start of the test
	err := testMgr.RevertToCustomSnapshot(snapshotName)
	if err != nil {
		fail("Error reverting to custom snapshot: %v", err)
	}

	// Reload the wallet to undo any changes made during the test
	err = testMgr.ServiceProvider.GetWallet().Reload(testMgr.GetLogger())
	if err != nil {
		fail("Error reloading wallet: %v", err)
	}
}
