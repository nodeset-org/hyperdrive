package swcommon

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/goccy/go-json"
	"github.com/nodeset-org/eth-utils/beacon"
	"github.com/nodeset-org/hyperdrive/daemon-utils/validator/utils"
	swconfig "github.com/nodeset-org/hyperdrive/modules/stakewise/shared/config"
	"github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/nodeset-org/hyperdrive/shared/types"
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
	withdrawalCreds := utils.GetWithdrawalCredsFromAddress(resources.Vault)

	// Create the new aggregated deposit data for all generated keys
	dataList := make([]*types.ExtendedDepositData, len(keys))
	for i, key := range keys {
		depositData, err := utils.GetDepositData(key, withdrawalCreds, resources.GenesisForkVersion, StakewiseDepositAmount, config.Network(resources.NodesetNetwork))
		if err != nil {
			pubkey := beacon.ValidatorPubkey(key.PublicKey().Marshal())
			return nil, fmt.Errorf("error getting deposit data for key %s: %w", pubkey.HexWithPrefix(), err)
		}
		dataList[i] = &depositData
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
