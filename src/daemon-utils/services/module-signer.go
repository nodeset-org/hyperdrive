package services

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/nodeset-org/hyperdrive/client"
)

// Used to request TX signatures from the node wallet
type ModuleSigner struct {
	hd *client.ApiClient
}

// Creates a new ModuleSigner
func NewModuleSigner(hd *client.ApiClient) *ModuleSigner {
	return &ModuleSigner{
		hd: hd,
	}
}

// Gets a transactor for signing transactions
func (s *ModuleSigner) GetTransactor(walletAddress common.Address) *bind.TransactOpts {
	return &bind.TransactOpts{
		From: walletAddress,
		Signer: func(from common.Address, tx *types.Transaction) (*types.Transaction, error) {
			if from != walletAddress {
				return nil, fmt.Errorf("cannot sign transactions from address %s (wallet address is %s)", from.Hex(), walletAddress.Hex())
			}
			return s.signTx(tx)
		},
	}
}

// Signs the TX by asking the Hyperdrive daemon to perform the actual signature
func (s *ModuleSigner) signTx(tx *types.Transaction) (*types.Transaction, error) {
	// Serialize it
	bytes, err := tx.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("error serializing TX: %w", err)
	}

	// Sign it
	response, err := s.hd.Wallet.SignTx(bytes)
	if err != nil {
		return nil, fmt.Errorf("error requesting TX signature: %w", err)
	}

	// Deserialize it
	var signedTx types.Transaction
	err = signedTx.UnmarshalBinary(response.Data.SignedTx)
	if err != nil {
		return nil, fmt.Errorf("error deserializing signed TX: %w", err)
	}

	return &signedTx, nil
}
