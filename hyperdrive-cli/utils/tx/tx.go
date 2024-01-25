package tx

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/nodeset-org/eth-utils/eth"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/gas"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
)

// Handle a transaction, either printing its details, signing it, or submitting it and waiting for it to be included
func HandleTx(c *cli.Context, hd *client.Client, txInfo *eth.TransactionInfo, confirmMessage string, identifier string, submissionMessage string) error {
	// Make sure the TX was successful
	if txInfo.SimulationResult.SimulationError != "" {
		return fmt.Errorf("simulating %s failed: %s", identifier, txInfo.SimulationResult.SimulationError)
	}

	// Print the TX data if requested
	if c.Bool(utils.PrintTxDataFlag) {
		fmt.Printf("TX Data for %s:\n", identifier)
		fmt.Printf("\tTo:       %s\n", txInfo.To.Hex())
		fmt.Printf("\tData:     %s\n", hexutil.Encode(txInfo.Data))
		fmt.Printf("\tValue:    %s\n", txInfo.Value.String())
		fmt.Printf("\tEst. Gas: %d\n", txInfo.SimulationResult.EstimatedGasLimit)
		fmt.Printf("\tSafe Gas: %d\n", txInfo.SimulationResult.SafeGasLimit)
		return nil
	}

	// Assign max fees
	maxFee, maxPrioFee, err := gas.GetMaxFees(c, hd, txInfo.SimulationResult)
	if err != nil {
		return fmt.Errorf("error getting fee information: %w", err)
	}

	// Check the nonce flag
	var nonce *big.Int
	if hd.Context.Nonce.Cmp(common.Big0) > 0 {
		nonce = hd.Context.Nonce
	}

	// Create the submission from the TX info
	submission, _ := eth.CreateTxSubmissionFromInfo(txInfo, nil)

	// Sign only (no submission) if requested
	if c.Bool(utils.SignTxOnlyFlag) {
		response, err := hd.Api.Tx.SignTx(submission, nonce, maxFee, maxPrioFee)
		if err != nil {
			return fmt.Errorf("error signing transaction: %w", err)
		}
		fmt.Printf("Signed transaction (%s):\n", identifier)
		fmt.Println(response.Data.SignedTx)
		fmt.Println()
		updateCustomNonce(hd)
		return nil
	}

	// Confirm submission
	if !(c.Bool(utils.YesFlag.Name) || utils.Confirm(confirmMessage)) {
		fmt.Println("Cancelled.")
		return nil
	}

	// Submit it
	fmt.Println(submissionMessage)
	response, err := hd.Api.Tx.SubmitTx(submission, nonce, maxFee, maxPrioFee)
	if err != nil {
		return fmt.Errorf("error submitting transaction: %w", err)
	}

	// Wait for it
	utils.PrintTransactionHash(hd, response.Data.TxHash)
	if _, err = hd.Api.Tx.WaitForTransaction(response.Data.TxHash); err != nil {
		return fmt.Errorf("error waiting for transaction: %w", err)
	}

	updateCustomNonce(hd)
	return nil
}

// Handle a batch of transactions, either printing their details, signing them, or submitting them and waiting for them to be included
func HandleTxBatch(c *cli.Context, hd *client.Client, txInfos []*eth.TransactionInfo, confirmMessage string, identifierFunc func(int) string, submissionMessage string) error {
	// Make sure the TXs were successful
	for i, txInfo := range txInfos {
		if txInfo.SimulationResult.SimulationError != "" {
			return fmt.Errorf("simulating %s failed: %s", identifierFunc(i), txInfo.SimulationResult.SimulationError)
		}
	}

	// Print the TX data if requested
	if c.Bool(utils.PrintTxDataFlag) {
		for i, info := range txInfos {
			fmt.Printf("Data for TX %d (%s):\n", i, identifierFunc(i))
			fmt.Printf("\tTo:       %s\n", info.To.Hex())
			fmt.Printf("\tData:     %s\n", hexutil.Encode(info.Data))
			fmt.Printf("\tValue:    %s\n", info.Value.String())
			fmt.Printf("\tEst. Gas: %d\n", info.SimulationResult.EstimatedGasLimit)
			fmt.Printf("\tSafe Gas: %d\n", info.SimulationResult.SafeGasLimit)
			fmt.Println()
		}
		return nil
	}

	// Assign max fees
	var simResult eth.SimulationResult
	for _, info := range txInfos {
		simResult.EstimatedGasLimit += info.SimulationResult.EstimatedGasLimit
		simResult.SafeGasLimit += info.SimulationResult.SafeGasLimit
	}
	maxFee, maxPrioFee, err := gas.GetMaxFees(c, hd, simResult)
	if err != nil {
		return fmt.Errorf("error getting fee information: %w", err)
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
	if c.Bool(utils.SignTxOnlyFlag) {
		response, err := hd.Api.Tx.SignTxBatch(submissions, nonce, maxFee, maxPrioFee)
		if err != nil {
			return fmt.Errorf("error signing transactions: %w", err)
		}

		for i, tx := range response.Data.SignedTxs {
			fmt.Printf("Signed transaction (%s):\n", identifierFunc(i))
			fmt.Println(tx)
			fmt.Println()
		}
		return nil
	}

	// Confirm submission
	if !(c.Bool(utils.YesFlag.Name) || utils.Confirm(confirmMessage)) {
		fmt.Println("Cancelled.")
		return nil
	}

	// Submit them
	fmt.Println(submissionMessage)
	response, err := hd.Api.Tx.SubmitTxBatch(submissions, nonce, maxFee, maxPrioFee)
	if err != nil {
		return fmt.Errorf("error submitting transactions: %w", err)
	}

	// Wait for them
	utils.PrintTransactionBatchHashes(hd, response.Data.TxHashes)
	return waitForTransactions(hd, response.Data.TxHashes, identifierFunc)
}

// Wait for a batch of transactions to get included in blocks
func waitForTransactions(hd *client.Client, hashes []common.Hash, identifierFunc func(int) string) error {
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
func updateCustomNonce(hd *client.Client) {
	if hd.Context.Nonce.Cmp(common.Big0) > 0 {
		hd.Context.Nonce.Add(hd.Context.Nonce, common.Big1)
	}
}
