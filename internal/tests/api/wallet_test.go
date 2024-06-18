package client

import (
	"math/big"
	"os"
	"runtime/debug"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/nodeset-org/nodeset-svc-mock/auth"
	"github.com/nodeset-org/osha"
	"github.com/nodeset-org/osha/keys"
	"github.com/rocket-pool/node-manager-core/eth"
	"github.com/rocket-pool/node-manager-core/node/services"
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
	require.Equal(t, emptyWalletAddress, response.Data.WalletStatus.Address.NodeAddress)
	require.False(t, response.Data.WalletStatus.Address.HasAddress)

	require.Equal(t, wallet.WalletType(""), response.Data.WalletStatus.Wallet.Type)
	require.False(t, response.Data.WalletStatus.Wallet.IsLoaded)
	require.False(t, response.Data.WalletStatus.Wallet.IsOnDisk)
	require.Equal(t, emptyWalletAddress, response.Data.WalletStatus.Wallet.WalletAddress)

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
	require.Equal(t, expectedWalletAddress, response.Data.WalletStatus.Address.NodeAddress)
	require.True(t, response.Data.WalletStatus.Address.HasAddress)

	require.Equal(t, wallet.WalletType_Local, response.Data.WalletStatus.Wallet.Type)
	require.True(t, response.Data.WalletStatus.Wallet.IsLoaded)
	require.True(t, response.Data.WalletStatus.Wallet.IsOnDisk)
	require.Equal(t, expectedWalletAddress, response.Data.WalletStatus.Wallet.WalletAddress)

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

func TestWalletSignMessage(t *testing.T) {
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

	message := []byte("hello world")
	response, err := apiClient.Wallet.SignMessage(message)
	require.NoError(t, err)
	t.Log("SignMessage called")

	require.NotEmpty(t, response.Data.SignedMessage)

	// Make sure that the recovered address is the signer address
	recoveredAddr, err := auth.GetAddressFromSignature(message, response.Data.SignedMessage)
	require.NoError(t, err)

	require.Equal(t, expectedWalletAddress, recoveredAddr)
	t.Logf("Successfully signed message")

}

func TestWalletSend_EthSuccess(t *testing.T) {
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

	targetAddress := common.HexToAddress("0x95222290dd7278aa3ddd389cc1e1d165cc4bafe5")
	response, err := apiClient.Wallet.Send(eth.EthToWei(1), "eth", targetAddress)
	require.NoError(t, err)
	t.Log("Send called")

	require.Equal(t, targetAddress, response.Data.TxInfo.To)
	require.Equal(t, eth.EthToWei(1), response.Data.TxInfo.Value)
	require.NotEmpty(t, response.Data.TxInfo.SimulationResult)

	require.True(t, response.Data.CanSend)
	require.False(t, response.Data.InsufficientBalance)
	t.Logf("Successfully generated transaction info for sending ETH")

	sub, _ := eth.CreateTxSubmissionFromInfo(response.Data.TxInfo, nil)
	submitResponse, err := apiClient.Tx.SubmitTx(sub, nil, eth.GweiToWei(10), eth.GweiToWei(1))
	require.NoError(t, err)
	t.Log("SubmitTx called")

	err = testMgr.CommitBlock()
	require.NoError(t, err)

	_, err = apiClient.Tx.WaitForTransaction(submitResponse.Data.TxHash)
	require.NoError(t, err)
	t.Log("Waiting complete")

	// Check the balance
	tm, err := osha.NewTestManager()
	require.NoError(t, err)
	beaconCfg := tm.GetBeaconMockManager().GetConfig()
	sp := testMgr.ServiceProvider
	ctx := sp.GetBaseContext()

	ecManager := services.NewExecutionClientManager(tm.GetExecutionClient(), uint(beaconCfg.ChainID), time.Minute)
	targetAddressBalance, err := ecManager.BalanceAt(ctx, targetAddress, nil)
	require.NoError(t, err)
	require.Equal(t, eth.EthToWei(1), targetAddressBalance)

	expectedWalletBalance, err := ecManager.BalanceAt(ctx, expectedWalletAddress, nil)
	require.NoError(t, err)

	require.True(t, expectedWalletBalance.Cmp(eth.EthToWei(99999)) < 0)
	t.Logf("Successfully sent ETH to target address")
}

func TestWalletSend_EthFailure(t *testing.T) {
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

	// Attempt to send too much ETH
	targetAddress := common.HexToAddress("0x95222290dd7278aa3ddd389cc1e1d165cc4bafe5")
	response, err := apiClient.Wallet.Send(eth.EthToWei(99999), "eth", targetAddress)
	require.NoError(t, err)
	t.Log("Send called")

	require.Empty(t, response.Data.TxInfo)

	require.False(t, response.Data.CanSend)
	require.True(t, response.Data.InsufficientBalance)
	t.Logf("Response correctly indicates insufficient balance")

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
