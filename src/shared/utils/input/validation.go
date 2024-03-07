package input

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/goccy/go-json"
	"github.com/nodeset-org/eth-utils/beacon"
	"github.com/rocket-pool/node-manager-core/eth"
	"github.com/tyler-smith/go-bip39"
	"github.com/urfave/cli/v2"
)

// Config
const (
	MinPasswordLength int = 12
)

//
// General types
//

// Validate command argument count
func ValidateArgCount(c *cli.Context, expectedCount int) error {
	argCount := c.Args().Len()
	if argCount != expectedCount {
		return fmt.Errorf("Incorrect argument count; expected %d but have %d", expectedCount, argCount)
	}
	return nil
}

// Validate a comma-delimited batch of inputs
func ValidateBatch[ReturnType any](name string, value string, validate func(string, string) (ReturnType, error)) ([]ReturnType, error) {
	elements := strings.Split(value, ",")
	results := make([]ReturnType, len(elements))
	for i, element := range elements {
		element = strings.TrimSpace(element)
		result, err := validate(name, element)
		if err != nil {
			return nil, fmt.Errorf("invalid element at index %d in %s: %w", i, name, err)
		}
		results[i] = result
	}
	return results, nil
}

// Validate a big int
func ValidateBigInt(name, value string) (*big.Int, error) {
	val, success := big.NewInt(0).SetString(value, 0)
	if !success {
		return nil, fmt.Errorf("Invalid %s '%s'", name, value)
	}
	return val, nil
}

// Validate a boolean value
func ValidateBool(name, value string) (bool, error) {
	val := strings.ToLower(value)
	if !(val == "true" || val == "yes" || val == "false" || val == "no") {
		return false, fmt.Errorf("Invalid %s '%s' - valid values are 'true', 'yes', 'false' and 'no'", name, value)
	}
	if val == "true" || val == "yes" {
		return true, nil
	}
	return false, nil
}

// Validate an unsigned integer value
func ValidateUint(name, value string) (uint64, error) {
	val, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("Invalid %s '%s'", name, value)
	}
	return val, nil
}

// Validate an unsigned integer value
func ValidateUint32(name, value string) (uint32, error) {
	val, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("Invalid %s '%s'", name, value)
	}
	return uint32(val), nil
}

// Validate an address
func ValidateAddress(name, value string) (common.Address, error) {
	if !common.IsHexAddress(value) {
		return common.Address{}, fmt.Errorf("Invalid %s '%s'", name, value)
	}
	return common.HexToAddress(value), nil
}

// Validate a wei amount
func ValidateWeiAmount(name, value string) (*big.Int, error) {
	val := new(big.Int)
	if _, ok := val.SetString(value, 10); !ok {
		return nil, fmt.Errorf("Invalid %s '%s'", name, value)
	}
	return val, nil
}

// Validate an ether amount
func ValidateEthAmount(name, value string) (float64, error) {
	val, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf("Invalid %s '%s'", name, value)
	}
	return val, nil
}

// Validate a fraction
func ValidateFraction(name, value string) (float64, error) {
	val, err := strconv.ParseFloat(value, 64)
	if err != nil || val < 0 || val > 1 {
		return 0, fmt.Errorf("Invalid %s '%s' - must be a number between 0 and 1", name, value)
	}
	return val, nil
}

// Validate a percentage
func ValidatePercentage(name, value string) (float64, error) {
	val, err := strconv.ParseFloat(value, 64)
	if err != nil || val < 0 || val > 100 {
		return 0, fmt.Errorf("Invalid %s '%s' - must be a number between 0 and 100", name, value)
	}
	return val, nil
}

//
// Command specific types
//

// Validate a positive unsigned integer value
func ValidatePositiveUint(name, value string) (uint64, error) {
	val, err := ValidateUint(name, value)
	if err != nil {
		return 0, err
	}
	if val == 0 {
		return 0, fmt.Errorf("Invalid %s '%s' - must be greater than 0", name, value)
	}
	return val, nil
}

// Validate a positive 32-bit unsigned integer value
func ValidatePositiveUint32(name, value string) (uint32, error) {
	val, err := ValidateUint32(name, value)
	if err != nil {
		return 0, err
	}
	if val == 0 {
		return 0, fmt.Errorf("Invalid %s '%s' - must be greater than 0", name, value)
	}
	return val, nil
}

// Validate a positive wei amount
func ValidatePositiveWeiAmount(name, value string) (*big.Int, error) {
	val, err := ValidateWeiAmount(name, value)
	if err != nil {
		return nil, err
	}
	if val.Cmp(big.NewInt(0)) < 1 {
		return nil, fmt.Errorf("Invalid %s '%s' - must be greater than 0", name, value)
	}
	return val, nil
}

// Validate a positive or zero wei amount
func ValidatePositiveOrZeroWeiAmount(name, value string) (*big.Int, error) {
	val, err := ValidateWeiAmount(name, value)
	if err != nil {
		return nil, err
	}
	if val.Cmp(big.NewInt(0)) < 0 {
		return nil, fmt.Errorf("Invalid %s '%s' - must be greater or equal to 0", name, value)
	}
	return val, nil
}

// Validate a positive ether amount
func ValidatePositiveEthAmount(name, value string) (float64, error) {
	val, err := ValidateEthAmount(name, value)
	if err != nil {
		return 0, err
	}
	if val <= 0 {
		return 0, fmt.Errorf("Invalid %s '%s' - must be greater than 0", name, value)
	}
	return val, nil
}

// Validate a node password
func ValidateNodePassword(name string, value string) (string, error) {
	if len(value) < MinPasswordLength {
		return "", fmt.Errorf("Invalid %s '%s' - must be at least %d characters long", name, value, MinPasswordLength)
	}
	return value, nil
}

// Validate a wallet mnemonic phrase
func ValidateWalletMnemonic(name, value string) (string, error) {
	if !bip39.IsMnemonicValid(value) {
		return "", fmt.Errorf("Invalid %s '%s'", name, value)
	}
	return value, nil
}

// Validate a timezone location
func ValidateTimezoneLocation(name, value string) (string, error) {
	if !regexp.MustCompile("^([a-zA-Z_]{2,}\\/)+[a-zA-Z_]{2,}$").MatchString(value) {
		return "", fmt.Errorf("Invalid %s '%s' - must be in the format 'Country/City'", name, value)
	}
	return value, nil
}

// Validate a hash
func ValidateHash(name, value string) (common.Hash, error) {
	// Remove a 0x prefix if present
	value = strings.TrimPrefix(value, "0x")

	// Hash should be 64 characters long
	if len(value) != hex.EncodedLen(common.HashLength) {
		return common.Hash{}, fmt.Errorf("Invalid %s '%s': it must have 64 characters.", name, value)
	}

	// Try to parse the string (removing the prefix)
	bytes, err := hex.DecodeString(value)
	if err != nil {
		return common.Hash{}, fmt.Errorf("Invalid %s '%s': %w", name, value, err)
	}
	hash := common.BytesToHash(bytes)

	return hash, nil
}

// Validate TX info
func ValidateTxInfo(name string, value string) (*eth.TransactionInfo, error) {
	// Remove a 0x prefix if present
	value = strings.TrimPrefix(value, "0x")

	// Try to parse the string
	bytes, err := hex.DecodeString(value)
	if err != nil {
		return nil, fmt.Errorf("Invalid %s '%s': %w", name, value, err)
	}

	// Deserialize it
	var info eth.TransactionInfo
	err = json.Unmarshal(bytes, &info)
	if err != nil {
		return nil, fmt.Errorf("Deserializing %s failed: %w", name, err)
	}

	return &info, nil
}

// Validate a validator pubkey
func ValidatePubkey(name, value string) (beacon.ValidatorPubkey, error) {
	pubkey, err := beacon.HexToValidatorPubkey(value)
	if err != nil {
		return beacon.ValidatorPubkey{}, fmt.Errorf("Invalid %s '%s': %w", name, value, err)
	}
	return pubkey, nil
}

// Validate a hex-encoded byte array
func ValidateByteArray(name, value string) ([]byte, error) {
	// Remove a 0x prefix if present
	if strings.HasPrefix(value, "0x") {
		value = value[2:]
	}

	// Try to parse the string (removing the prefix)
	bytes, err := hex.DecodeString(value)
	if err != nil {
		return nil, fmt.Errorf("Invalid %s '%s': %w", name, value, err)
	}

	return bytes, nil
}

// Validate a duration
func ValidateDuration(name, value string) (time.Duration, error) {
	duration, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("Invalid %s '%s': %w", name, value, err)
	}
	return duration, nil
}

// Validate a timestamp using RFC3339
func ValidateTime(name, value string) (time.Time, error) {
	val, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("Invalid %s '%s': %w", name, value, err)
	}
	return val, nil
}
