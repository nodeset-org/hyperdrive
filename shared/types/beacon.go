package types

import (
	"encoding/hex"

	"github.com/ethereum/go-ethereum/common"
	"github.com/prysmaticlabs/go-bitfield"
)

// Validator pubkey
const ValidatorPubkeyLength = 48 // bytes
type ValidatorPubkey [ValidatorPubkeyLength]byte

// Validator signature
const ValidatorSignatureLength = 96 // bytes
type ValidatorSignature [ValidatorSignatureLength]byte

// Bytes conversion
func (v ValidatorPubkey) Bytes() []byte {
	return v[:]
}
func BytesToValidatorPubkey(value []byte) ValidatorPubkey {
	var pubkey ValidatorPubkey
	copy(pubkey[:], value)
	return pubkey
}

// String conversion
func (v ValidatorPubkey) Hex() string {
	return hex.EncodeToString(v.Bytes())
}
func (v ValidatorPubkey) String() string {
	return v.Hex()
}

// Bytes conversion
func (v ValidatorSignature) Bytes() []byte {
	return v[:]
}
func BytesToValidatorSignature(value []byte) ValidatorSignature {
	var signature ValidatorSignature
	copy(signature[:], value)
	return signature
}

// String conversion
func (v ValidatorSignature) Hex() string {
	return hex.EncodeToString(v.Bytes())
}
func (v ValidatorSignature) String() string {
	return v.Hex()
}

// API request options
type ValidatorStatusOptions struct {
	Epoch *uint64
	Slot  *uint64
}

// API response types
type SyncStatus struct {
	Syncing  bool
	Progress float64
}
type Eth2Config struct {
	GenesisForkVersion           []byte
	GenesisValidatorsRoot        []byte
	GenesisEpoch                 uint64
	GenesisTime                  uint64
	SecondsPerSlot               uint64
	SlotsPerEpoch                uint64
	SecondsPerEpoch              uint64
	EpochsPerSyncCommitteePeriod uint64
}
type Eth2DepositContract struct {
	ChainID uint64
	Address common.Address
}
type BeaconHead struct {
	Epoch                  uint64
	FinalizedEpoch         uint64
	JustifiedEpoch         uint64
	PreviousJustifiedEpoch uint64
}
type ValidatorStatus struct {
	Pubkey                     ValidatorPubkey
	Index                      string
	WithdrawalCredentials      common.Hash
	Balance                    uint64
	Status                     ValidatorState
	EffectiveBalance           uint64
	Slashed                    bool
	ActivationEligibilityEpoch uint64
	ActivationEpoch            uint64
	ExitEpoch                  uint64
	WithdrawableEpoch          uint64
	Exists                     bool
}
type Eth1Data struct {
	DepositRoot  common.Hash
	DepositCount uint64
	BlockHash    common.Hash
}
type BeaconBlock struct {
	Slot                 uint64
	ProposerIndex        string
	HasExecutionPayload  bool
	Attestations         []AttestationInfo
	FeeRecipient         common.Address
	ExecutionBlockNumber uint64
}
type AttestationInfo struct {
	AggregationBits bitfield.Bitlist
	SlotIndex       uint64
	CommitteeIndex  uint64
}

type ValidatorState string

const (
	ValidatorState_PendingInitialized ValidatorState = "pending_initialized"
	ValidatorState_PendingQueued      ValidatorState = "pending_queued"
	ValidatorState_ActiveOngoing      ValidatorState = "active_ongoing"
	ValidatorState_ActiveExiting      ValidatorState = "active_exiting"
	ValidatorState_ActiveSlashed      ValidatorState = "active_slashed"
	ValidatorState_ExitedUnslashed    ValidatorState = "exited_unslashed"
	ValidatorState_ExitedSlashed      ValidatorState = "exited_slashed"
	ValidatorState_WithdrawalPossible ValidatorState = "withdrawal_possible"
	ValidatorState_WithdrawalDone     ValidatorState = "withdrawal_done"
)

// Beacon client interface
type IBeaconClient interface {
	GetSyncStatus() (SyncStatus, error)
	GetEth2Config() (Eth2Config, error)
	GetEth2DepositContract() (Eth2DepositContract, error)
	GetAttestations(blockId string) ([]AttestationInfo, bool, error)
	GetBeaconBlock(blockId string) (BeaconBlock, bool, error)
	GetBeaconHead() (BeaconHead, error)
	GetValidatorStatusByIndex(index string, opts *ValidatorStatusOptions) (ValidatorStatus, error)
	GetValidatorStatus(pubkey ValidatorPubkey, opts *ValidatorStatusOptions) (ValidatorStatus, error)
	GetValidatorStatuses(pubkeys []ValidatorPubkey, opts *ValidatorStatusOptions) (map[ValidatorPubkey]ValidatorStatus, error)
	GetValidatorIndex(pubkey ValidatorPubkey) (string, error)
	GetValidatorSyncDuties(indices []string, epoch uint64) (map[string]bool, error)
	GetValidatorProposerDuties(indices []string, epoch uint64) (map[string]uint64, error)
	GetDomainData(domainType []byte, epoch uint64, useGenesisFork bool) ([]byte, error)
	ExitValidator(validatorIndex string, epoch uint64, signature ValidatorSignature) error
	Close() error
	GetEth1DataForEth2Block(blockId string) (Eth1Data, bool, error)
	ChangeWithdrawalCredentials(validatorIndex string, fromBlsPubkey ValidatorPubkey, toExecutionAddress common.Address, signature ValidatorSignature) error
}
