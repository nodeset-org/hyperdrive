package wallet

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/nodeset-org/hyperdrive/shared/utils/log"
)

// Recover a wallet keystore from a mnemonic - only used for testing mnemonics
func TestRecovery(derivationPath string, walletIndex uint, mnemonic string, chainID uint) (*Wallet, error) {
	// Create a new dummy wallet with a fake password
	log := log.NewColorLogger(color.FgHiWhite)
	w, err := NewWallet(&log, "", "", "", chainID)
	if err != nil {
		return nil, fmt.Errorf("error creating new test node wallet: %w", err)
	}

	err = w.Recover(derivationPath, walletIndex, mnemonic, "test password", false, true)
	if err != nil {
		return nil, fmt.Errorf("error test recovering mnemonic: %w", err)
	}
	return w, nil
}
