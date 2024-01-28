package types

import (
	"github.com/google/uuid"
	"github.com/nodeset-org/eth-utils/beacon"
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
