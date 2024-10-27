package tx

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/gas"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/terminal"
	"github.com/rocket-pool/node-manager-core/eth"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
)

// Handle a transaction, either printing its details, signing it, or submitting it and waiting for it to be included
func HandleTx(c *cli.Context, hd *client.HyperdriveClient, txInfo *eth.TransactionInfo, confirmMessage string, identifier string, submissionMessage string) (bool, error) {
	return handleTxImpl(c, hd, txInfo, confirmMessage, identifier, submissionMessage, false)
}

// Handle a transaction, either printing its details, signing it, or submitting it and waiting for it to be included, explicitly requiring the user to enter "I Agree" in the confirmation prompt
func HandleTxWithIAgree(c *cli.Context, hd *client.HyperdriveClient, txInfo *eth.TransactionInfo, confirmMessage string, identifier string, submissionMessage string) (bool, error) {
	return handleTxImpl(c, hd, txInfo, confirmMessage, identifier, submissionMessage, true)
}

// Implementation of the transaction handling logic
func handleTxImpl(c *cli.Context, hd *client.HyperdriveClient, txInfo *eth.TransactionInfo, confirmMessage string, identifier string, submissionMessage string, useIAgree bool) (bool, error) {
	// Overwrite the gas limit if requested
	if c.IsSet(utils.ForceGasLimitFlag.Name) {
		manualLimit := c.Uint64(utils.ForceGasLimitFlag.Name)
		txInfo.SimulationResult.EstimatedGasLimit = manualLimit
		txInfo.SimulationResult.SafeGasLimit = manualLimit
	}

	// Print the TX data if requested
	if c.Bool(utils.PrintTxDataFlag.Name) {
		fmt.Printf("TX Data for %s:\n", identifier)
		fmt.Printf("\tTo:       %s\n", txInfo.To.Hex())
		fmt.Printf("\tData:     %s\n", hexutil.Encode(txInfo.Data))
		fmt.Printf("\tValue:    %s\n", txInfo.Value.String())
		fmt.Printf("\tEst. Gas: %d\n", txInfo.SimulationResult.EstimatedGasLimit)
		fmt.Printf("\tSafe Gas: %d\n", txInfo.SimulationResult.SafeGasLimit)

		// Warn if the TX failed simulation
		if txInfo.SimulationResult.SimulationError != "" {
			fmt.Printf("%sWARNING: '%s' failed simulation: %s\nThis transaction will likely revert if you submit it.%s\n", terminal.ColorYellow, identifier, txInfo.SimulationResult.SimulationError, terminal.ColorReset)
		}
		return false, nil
	}

	// Make sure the TX was successful
	if txInfo.SimulationResult.SimulationError != "" {
		if c.Bool(utils.IgnoreTxSimFailureFlag.Name) {
			fmt.Printf("%sWARNING: '%s' failed simulation: %s\nThis transaction will likely revert if you submit it.%s\n", terminal.ColorYellow, identifier, txInfo.SimulationResult.SimulationError, terminal.ColorReset)
		} else {
			return false, fmt.Errorf("simulating %s failed: %s", identifier, txInfo.SimulationResult.SimulationError)
		}
	}

	// Assign max fees
	maxFee, maxPrioFee, err := gas.GetMaxFees(c, hd, txInfo.SimulationResult)
	if err != nil {
		return false, fmt.Errorf("error getting fee information: %w", err)
	}

	// Check the nonce flag
	var nonce *big.Int
	if hd.Context.Nonce.Cmp(common.Big0) > 0 {
		nonce = hd.Context.Nonce
	}

	// Create the submission from the TX info
	submission, _ := eth.CreateTxSubmissionFromInfo(txInfo, nil)

	// Sign only (no submission) if requested
	if c.Bool(utils.SignTxOnlyFlag.Name) {
		response, err := hd.Api.Tx.SignTx(submission, nonce, maxFee, maxPrioFee)
		if err != nil {
			return false, fmt.Errorf("error signing transaction: %w", err)
		}
		fmt.Printf("Signed transaction (%s):\n", identifier)
		fmt.Println(response.Data.SignedTx)
		fmt.Println()
		updateCustomNonce(hd)
		return false, nil
	}

	// Confirm submission
	var confirmFunc func(string) bool
	if useIAgree {
		confirmFunc = utils.ConfirmWithIAgree
	} else {
		confirmFunc = utils.Confirm
	}
	if !(c.Bool(utils.YesFlag.Name) || confirmFunc(confirmMessage)) {
		fmt.Println("Cancelled.")
		return false, nil
	}

	// Submit it
	fmt.Println(submissionMessage)
	response, err := hd.Api.Tx.SubmitTx(submission, nonce, maxFee, maxPrioFee)
	if err != nil {
		return false, fmt.Errorf("error submitting transaction: %w", err)
	}

	// Wait for it
	utils.PrintTransactionHash(hd, response.Data.TxHash)
	if _, err = hd.Api.Tx.WaitForTransaction(response.Data.TxHash); err != nil {
		return false, fmt.Errorf("error waiting for transaction: %w", err)
	}

	updateCustomNonce(hd)
	return true, nil
}

// Handle a batch of transactions, either printing their details, signing them, or submitting them and waiting for them to be included
func HandleTxBatch(c *cli.Context, hd *client.HyperdriveClient, txInfos []*eth.TransactionInfo, confirmMessage string, identifierFunc func(int) string, submissionMessage string) (bool, error) {
	// Overwrite the gas limit if requested
	if c.IsSet(utils.ForceGasLimitFlag.Name) {
		manualLimit := c.Uint64(utils.ForceGasLimitFlag.Name)
		for _, txInfo := range txInfos {
			txInfo.SimulationResult.EstimatedGasLimit = manualLimit
			txInfo.SimulationResult.SafeGasLimit = manualLimit
		}
	}

	// Print the TX data if requested
	if c.Bool(utils.PrintTxDataFlag.Name) {
		for i, info := range txInfos {
			id := identifierFunc(i)
			fmt.Printf("Data for TX %d (%s):\n", i, identifierFunc(i))
			fmt.Printf("\tTo:       %s\n", info.To.Hex())
			fmt.Printf("\tData:     %s\n", hexutil.Encode(info.Data))
			fmt.Printf("\tValue:    %s\n", info.Value.String())
			fmt.Printf("\tEst. Gas: %d\n", info.SimulationResult.EstimatedGasLimit)
			fmt.Printf("\tSafe Gas: %d\n", info.SimulationResult.SafeGasLimit)
			fmt.Println()

			// Warn if the TX failed simulation
			if info.SimulationResult.SimulationError != "" {
				fmt.Printf("%sWARNING: '%s' failed simulation: %s\nThis transaction will likely revert if you submit it.%s\n", terminal.ColorYellow, id, info.SimulationResult.SimulationError, terminal.ColorReset)
				fmt.Println()
			}
		}
		return false, nil
	}

	// Make sure the TXs were successful
	for i, txInfo := range txInfos {
		if txInfo.SimulationResult.SimulationError != "" {
			if c.Bool(utils.IgnoreTxSimFailureFlag.Name) {
				fmt.Printf("%sWARNING: '%s' failed simulation: %s\nThis transaction will likely revert if you submit it.%s\n", terminal.ColorYellow, identifierFunc(i), txInfo.SimulationResult.SimulationError, terminal.ColorReset)
			} else {
				return false, fmt.Errorf("simulating %s failed: %s", identifierFunc(i), txInfo.SimulationResult.SimulationError)
			}
		}
	}

	// Assign max fees
	var simResult eth.SimulationResult
	for _, info := range txInfos {
		simResult.EstimatedGasLimit += info.SimulationResult.EstimatedGasLimit
		simResult.SafeGasLimit += info.SimulationResult.SafeGasLimit
	}
	maxFee, maxPrioFee, err := gas.GetMaxFees(c, hd, simResult)
	if err != nil {
		return false, fmt.Errorf("error getting fee information: %w", err)
	}

	// Check the nonce flag
	var nonce *big.Int
	if hd.Context.Nonce.Cmp(common.Big0) > 0 {
		nonce = hd.Context.Nonce
	}

	// Create the submissions from the TX infos
	submissions := make([]*eth.TransactionSubmission, len(txInfos))
	for i, info := range txInfos {
		submission, _ := eth.CreateTxSubmissionFromInfo(info, nil)
		submissions[i] = submission
	}

	// Sign only (no submission) if requested
	if c.Bool(utils.SignTxOnlyFlag.Name) {
		response, err := hd.Api.Tx.SignTxBatch(submissions, nonce, maxFee, maxPrioFee)
		if err != nil {
			return false, fmt.Errorf("error signing transactions: %w", err)
		}

		for i, tx := range response.Data.SignedTxs {
			fmt.Printf("Signed transaction (%s):\n", identifierFunc(i))
			fmt.Println(tx)
			fmt.Println()
		}
		return false, nil
	}

	// Confirm submission
	if !(c.Bool(utils.YesFlag.Name) || utils.Confirm(confirmMessage)) {
		fmt.Println("Cancelled.")
		return false, nil
	}

	// Submit them
	fmt.Println(submissionMessage)
	response, err := hd.Api.Tx.SubmitTxBatch(submissions, nonce, maxFee, maxPrioFee)
	if err != nil {
		return false, fmt.Errorf("error submitting transactions: %w", err)
	}

	// Wait for them
	utils.PrintTransactionBatchHashes(hd, response.Data.TxHashes)
	return true, waitForTransactions(hd, response.Data.TxHashes, identifierFunc)
}

// Wait for a batch of transactions to get included in blocks
func waitForTransactions(hd *client.HyperdriveClient, hashes []common.Hash, identifierFunc func(int) string) error {
	var wg errgroup.Group
	var lock sync.Mutex
	total := len(hashes)
	successCount := 0

	// Create waiters for each TX
	for i, hash := range hashes {
		i := i
		hash := hash
		wg.Go(func() error {
			if _, err := hd.Api.Tx.WaitForTransaction(hash); err != nil {
				return fmt.Errorf("error waiting for transaction %s: %w", hash.Hex(), err)
			}
			lock.Lock()
			successCount++
			fmt.Printf("TX %s (%s) complete (%d/%d)\n", hash.Hex(), identifierFunc(i), successCount, total)
			lock.Unlock()
			return nil
		})
	}
	if err := wg.Wait(); err != nil {
		return fmt.Errorf("error waiting for transactions: %w", err)
	}
	return nil
}

// If a custom nonce is set, increment it for the next transaction
func updateCustomNonce(hd *client.HyperdriveClient) {
	if hd.Context.Nonce.Cmp(common.Big0) > 0 {
		hd.Context.Nonce.Add(hd.Context.Nonce, common.Big1)
	}
}
