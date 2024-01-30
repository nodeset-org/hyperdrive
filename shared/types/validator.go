package types

import (
	"encoding/hex"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/nodeset-org/eth-utils/beacon"
)

const (
	// Hyperdrive distinguishes keys by module to prevent overlapping between modules.
	// It uses the `use` field of the path, as defined in EIP-2334, to represent each module.

	RocketPoolValidatorPath    string = "m/12381/3600/%d/0/0"
	StakewiseValidatorPath     string = "m/12381/3600/%d/1/0"
	ConstellationValidatorPath string = "m/12381/3600/%d/2/0"
	SoloValidatorPath          string = "m/12381/3600/%d/3/0"
)

// Encrypted validator keystore following the EIP-2335 standard
// (https://eips.ethereum.org/EIPS/eip-2335)
type ValidatorKeystore struct {
	Crypto  map[string]interface{} `json:"crypto"`
	Name    string                 `json:"name,omitempty"` // Technically not part of the spec but Prysm needs it
	Version uint                   `json:"version"`
	UUID    uuid.UUID              `json:"uuid"`
	Path    string                 `json:"path"`
	Pubkey  beacon.ValidatorPubkey `json:"pubkey,omitempty"`
}

// Extended deposit data beyond what is required in an actual deposit message to Beacon, emulating what the deposit CLI produces
type ExtendedDepositData struct {
	PublicKey             ByteArray `json:"pubkey"`
	WithdrawalCredentials ByteArray `json:"withdrawal_credentials"`
	Amount                uint64    `json:"amount"`
	Signature             ByteArray `json:"signature"`
	DepositMessageRoot    ByteArray `json:"deposit_message_root"`
	DepositDataRoot       ByteArray `json:"deposit_data_root"`
	ForkVersion           ByteArray `json:"fork_version"`
	NetworkName           string    `json:"network_name"`
	HyperdriveVersion     string    `json:"hyperdrive_version,omitempty"`
}

// Byte array type
type ByteArray []byte

func (b ByteArray) MarshalJSON() ([]byte, error) {
	return json.Marshal(hex.EncodeToString(b))
}
func (b *ByteArray) UnmarshalJSON(data []byte) error {

	// Unmarshal string
	var dataStr string
	if err := json.Unmarshal(data, &dataStr); err != nil {
		return err
	}

	// Decode hex
	value, err := hex.DecodeString(dataStr)
	if err != nil {
		return err
	}

	// Set value and return
	*b = value
	return nil
}
