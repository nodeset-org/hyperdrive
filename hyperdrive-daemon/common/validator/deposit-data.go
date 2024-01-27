package validator

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/nodeset-org/eth-utils/beacon"
	"github.com/nodeset-org/hyperdrive/shared/types"
	eth2types "github.com/wealdtech/go-eth2-types/v2"
)

// Get deposit data & root for a given validator key and withdrawal credentials
func GetDepositData(validatorKey *eth2types.BLSPrivateKey, withdrawalCredentials common.Hash, eth2Config types.Eth2Config, depositAmount uint64) (beacon.DepositData, common.Hash, error) {

	// Build deposit data
	dd := beacon.DepositDataNoSignature{
		PublicKey:             validatorKey.PublicKey().Marshal(),
		WithdrawalCredentials: withdrawalCredentials[:],
		Amount:                depositAmount,
	}

	// Get signing root
	or, err := dd.HashTreeRoot()
	if err != nil {
		return beacon.DepositData{}, common.Hash{}, err
	}

	sr := beacon.SigningRoot{
		ObjectRoot: or[:],
		Domain:     eth2types.Domain(eth2types.DomainDeposit, eth2Config.GenesisForkVersion, eth2types.ZeroGenesisValidatorsRoot),
	}

	// Get signing root with domain
	srHash, err := sr.HashTreeRoot()
	if err != nil {
		return beacon.DepositData{}, common.Hash{}, err
	}

	// Build deposit data struct (with signature)
	var depositData = beacon.DepositData{
		PublicKey:             dd.PublicKey,
		WithdrawalCredentials: dd.WithdrawalCredentials,
		Amount:                dd.Amount,
		Signature:             validatorKey.Sign(srHash[:]).Marshal(),
	}

	// Get deposit data root
	depositDataRoot, err := depositData.HashTreeRoot()
	if err != nil {
		return beacon.DepositData{}, common.Hash{}, err
	}

	// Return
	return depositData, depositDataRoot, nil

}
