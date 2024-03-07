package swcommon

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"math/big"
	"os"
	"path/filepath"
	"sort"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/goccy/go-json"
	swconfig "github.com/nodeset-org/hyperdrive/modules/stakewise/shared/config"
	"github.com/nodeset-org/hyperdrive/shared"
	"github.com/nodeset-org/hyperdrive/shared/types"
	nmc_beacon "github.com/rocket-pool/node-manager-core/beacon"
	nmc_config "github.com/rocket-pool/node-manager-core/config"
	nmc_utils "github.com/rocket-pool/node-manager-core/node/validator/utils"
	eth2types "github.com/wealdtech/go-eth2-types/v2"
)

const (
	// Stakewise validators deposit a full 32 ETH
	StakewiseDepositAmount uint64 = 32e9
)

// DepositDataManager manages the aggregated deposit data file that Stakewise uses
type DepositDataManager struct {
	dataPath string
	sp       *StakewiseServiceProvider
}

// Creates a new manager
func NewDepositDataManager(sp *StakewiseServiceProvider) (*DepositDataManager, error) {
	dataPath := filepath.Join(sp.GetModuleDir(), swconfig.DepositDataFile)

	ddMgr := &DepositDataManager{
		dataPath: filepath.Join(sp.GetModuleDir(), swconfig.DepositDataFile),
		sp:       sp,
	}

	// Initialize the file if it's not there
	_, err := os.Stat(dataPath)
	if errors.Is(err, fs.ErrNotExist) {
		// Make a blank one
		err = ddMgr.UpdateDepositData([]types.ExtendedDepositData{})
		return ddMgr, err
	}
	if err != nil {
		return nil, fmt.Errorf("error checking status of wallet file [%s]: %w", dataPath, err)
	}

	return ddMgr, nil
}

// Generates deposit data for the provided keys
func (m *DepositDataManager) GenerateDepositData(keys []*eth2types.BLSPrivateKey) ([]*types.ExtendedDepositData, error) {
	resources := m.sp.GetResources()

	// Stakewise uses the same withdrawal creds for each validator
	withdrawalCreds := nmc_utils.GetWithdrawalCredsFromAddress(resources.Vault)

	// Create the new aggregated deposit data for all generated keys
	dataList := make([]*types.ExtendedDepositData, len(keys))
	for i, key := range keys {
		depositData, err := nmc_utils.GetDepositData(key, withdrawalCreds, resources.GenesisForkVersion, StakewiseDepositAmount, nmc_config.Network(resources.NodesetNetwork))
		if err != nil {
			pubkey := nmc_beacon.ValidatorPubkey(key.PublicKey().Marshal())
			return nil, fmt.Errorf("error getting deposit data for key %s: %w", pubkey.HexWithPrefix(), err)
		}
		dataList[i] = &types.ExtendedDepositData{
			ExtendedDepositData: depositData,
			HyperdriveVersion:   shared.HyperdriveVersion,
		}
	}
	return dataList, nil
}

// Read the deposit data file
func (m *DepositDataManager) GetDepositData() ([]byte, error) {
	// Read the file
	bytes, err := os.ReadFile(m.dataPath)
	if err != nil {
		return nil, fmt.Errorf("error reading deposit data file [%s]: %w", m.dataPath, err)
	}

	// Make sure it can deserialize properly
	var depositData []types.ExtendedDepositData
	err = json.Unmarshal(bytes, &depositData)
	if err != nil {
		return nil, fmt.Errorf("error deserializing deposit data file [%s]: %w", m.dataPath, err)
	}

	return bytes, nil
}

// Save the deposit data file
func (m *DepositDataManager) UpdateDepositData(data []types.ExtendedDepositData) error {
	// Serialize it
	bytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error serializing deposit data: %w", err)
	}

	// Write it
	err = os.WriteFile(m.dataPath, bytes, fileMode)
	if err != nil {
		return fmt.Errorf("error saving deposit data to disk: %w", err)
	}

	return nil
}

// Compute the Merkle root of the aggregated deposit data using the Stakewise rules
// NOTE: reverse engineered from https://github.com/stakewise/v3-operator/blob/fa4ac2673a64a486ced51098005376e56e2ddd19/src/validators/utils.py#L207
func (m *DepositDataManager) ComputeMerkleRoot(data []types.ExtendedDepositData) (common.Hash, error) {
	leafCount := len(data)
	leaves := make([][]byte, leafCount)

	// Create leaf data for each deposit data
	for i, dd := range data {
		// Get the deposit data root for this deposit data
		ddRoot, err := m.regenerateDepositDataRoot(dd)
		if err != nil {
			pubkey := nmc_beacon.ValidatorPubkey(dd.PublicKey)
			return common.Hash{}, fmt.Errorf("error generating deposit data root for validator %d (%s): %w", i, pubkey.Hex(), err)
		}

		// Get the index
		index := big.NewInt(int64(i))

		// The Stakewise tree ABI encodes its leaves in a custom format, so we have to replicate that for Geth to use
		bytesType, _ := abi.NewType("bytes", "bytes", nil)
		uint256Type, _ := abi.NewType("uint256", "uint256", nil)
		args := abi.Arguments{
			{
				// This is a tweaked version of the deposit data at index i, "entry data"
				Type: bytesType,
			},
			{
				// This is index i
				Type: uint256Type,
			},
		}

		// entryData = pubkey :: signature :: ddRoot
		entryData := []byte{}
		entryData = append(entryData, dd.PublicKey...)
		entryData = append(entryData, dd.Signature...)
		entryData = append(entryData, ddRoot[:]...)

		// ABI encode entryData and index to produce the raw leaf node value
		bytes, err := args.Pack(entryData, index)
		if err != nil {
			return common.Hash{}, fmt.Errorf("error packing abi_encode args for %d: %w", i, err)
		}

		// Keccak256 the ABI-encoded data twice
		hash := crypto.Keccak256(bytes)
		hash = crypto.Keccak256(hash)
		leaves[i] = hash
	}

	// Sort the hashes
	sort.SliceStable(leaves, func(i, j int) bool {
		return bytes.Compare(leaves[i], leaves[j]) == -1
	})

	// Add the leaves to the "tree" backwards - note that this is a nonstandard tree; the leaves aren't padded so they don't all necessarily live at the same depth.
	treeLength := 2*leafCount - 1
	tree := make([][]byte, treeLength)
	for i, leaf := range leaves {
		tree[treeLength-1-i] = leaf
	}

	// Traverse up the "tree", calculating nodes from the children that are already present
	for i := treeLength - 1 - leafCount; i > -1; i-- {
		leftChild := tree[2*i+1]
		rightChild := tree[2*i+2]

		// Compute the hash in sorted mode
		var hash []byte
		if bytes.Compare(leftChild, rightChild) < 0 {
			hash = crypto.Keccak256(leftChild, rightChild)
		} else {
			hash = crypto.Keccak256(rightChild, leftChild)
		}
		tree[i] = hash
	}

	return common.Hash(tree[0]), nil
}

// Regenerate the deposit data hash root from a deposit data object instead of explicitly relying on the deposit data root provided in the EDD
func (m *DepositDataManager) regenerateDepositDataRoot(dd types.ExtendedDepositData) (common.Hash, error) {
	var depositData = nmc_beacon.DepositData{
		PublicKey:             dd.PublicKey,
		WithdrawalCredentials: dd.WithdrawalCredentials,
		Amount:                StakewiseDepositAmount, // Note: hardcoded here because Stakewise ignores the actual amount in the deposit data and hardcodes it in their tree generation
		Signature:             dd.Signature,
	}

	// Get deposit data root
	root, err := depositData.HashTreeRoot()
	if err != nil {
		return common.Hash{}, err
	}
	return root, nil
}
