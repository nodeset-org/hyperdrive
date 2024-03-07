package services

import (
	"context"
	"time"

	"github.com/rocket-pool/node-manager-core/eth"
)

const (
	ethClientRecentBlockThreshold time.Duration = 5 * time.Minute
)

// Confirm the EC's latest block is within the threshold of the current system clock
func IsSyncWithinThreshold(ec eth.IExecutionClient) (bool, time.Time, error) {
	timestamp, err := GetEthClientLatestBlockTimestamp(ec)
	if err != nil {
		return false, time.Time{}, err
	}

	// Return true if the latest block is under the threshold
	blockTime := time.Unix(int64(timestamp), 0)
	if time.Since(blockTime) < ethClientRecentBlockThreshold {
		return true, blockTime, nil
	}

	return false, blockTime, nil
}

func GetEthClientLatestBlockTimestamp(ec eth.IExecutionClient) (uint64, error) {
	// Get latest block
	header, err := ec.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return 0, err
	}

	// Return block timestamp
	return header.Time, nil
}
