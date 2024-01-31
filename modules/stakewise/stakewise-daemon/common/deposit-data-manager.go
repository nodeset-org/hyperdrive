package swcommon

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-json"
	"github.com/nodeset-org/eth-utils/beacon"
	"github.com/nodeset-org/hyperdrive/daemon-utils/validator/utils"
	swconfig "github.com/nodeset-org/hyperdrive/modules/stakewise/shared/config"
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
func NewDepositDataManager(sp *StakewiseServiceProvider) *DepositDataManager {
	return &DepositDataManager{
		dataPath: filepath.Join(sp.GetModuleDir(), swconfig.DepositDataFile),
		sp:       sp,
	}
}

// Generates deposit data for the provided keys
func (m *DepositDataManager) GenerateDepositData(keys []*eth2types.BLSPrivateKey) ([]*types.ExtendedDepositData, error) {
	resources := m.sp.GetResources()

	// Stakewise uses the same withdrawal creds for each validator
	withdrawalCreds := utils.GetWithdrawalCredsFromAddress(resources.Vault)

	// Create the new aggregated deposit data for all generated keys
	dataList := make([]*types.ExtendedDepositData, len(keys))
	for i, key := range keys {
		depositData, err := utils.GetDepositData(key, withdrawalCreds, resources.GenesisForkVersion, StakewiseDepositAmount, resources.Network)
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
