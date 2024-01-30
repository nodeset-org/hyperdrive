package utils

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/nodeset-org/eth-utils/beacon"
	"github.com/nodeset-org/hyperdrive/shared"
	"github.com/nodeset-org/hyperdrive/shared/types"
	eth2types "github.com/wealdtech/go-eth2-types/v2"
)

// Get deposit data & root for a given validator key and withdrawal credentials
func GetDepositData(validatorKey *eth2types.BLSPrivateKey, withdrawalCredentials common.Hash, genesisForkVersion []byte, depositAmount uint64, network types.Network) (types.ExtendedDepositData, error) {
	// Build deposit data
	dd := beacon.DepositDataNoSignature{
		PublicKey:             validatorKey.PublicKey().Marshal(),
		WithdrawalCredentials: withdrawalCredentials[:],
		Amount:                depositAmount,
	}
	domain, err := eth2types.ComputeDomain(eth2types.DomainDeposit, genesisForkVersion, eth2types.ZeroGenesisValidatorsRoot)
	if err != nil {
		return types.ExtendedDepositData{}, fmt.Errorf("error computing domain: %w", err)
	}

	// Get signing root
	messageRoot, err := dd.HashTreeRoot()
	if err != nil {
		return types.ExtendedDepositData{}, fmt.Errorf("error getting message root: %w", err)
	}
	dataRoot := beacon.SigningRoot{
		ObjectRoot: messageRoot[:],
		Domain:     domain,
	}

	// Get signing root with domain
	dataRootHash, err := dataRoot.HashTreeRoot()
	if err != nil {
		return types.ExtendedDepositData{}, err
	}

	// Build deposit data struct (with signature)
	var depositData = beacon.DepositData{
		PublicKey:             dd.PublicKey,
		WithdrawalCredentials: dd.WithdrawalCredentials,
		Amount:                dd.Amount,
		Signature:             validatorKey.Sign(dataRootHash[:]).Marshal(),
	}

	// Get deposit data root
	depositDataRoot, err := depositData.HashTreeRoot()
	if err != nil {
		return types.ExtendedDepositData{}, err
	}

	// Create the extended data
	return types.ExtendedDepositData{
		PublicKey:             depositData.PublicKey,
		WithdrawalCredentials: depositData.WithdrawalCredentials,
		Amount:                depositData.Amount,
		Signature:             depositData.Signature,
		DepositMessageRoot:    messageRoot[:],
		DepositDataRoot:       depositDataRoot[:],
		ForkVersion:           genesisForkVersion,
		NetworkName:           string(network),
		HyperdriveVersion:     shared.HyperdriveVersion,
	}, nil
}
