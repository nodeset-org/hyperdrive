// Beacon-client manager
package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/sync/errgroup"
)

// Config
const (
	RequestUrlFormat   = "%s%s"
	RequestContentType = "application/json"

	// RequestSyncStatusPath                  = "/eth/v1/node/syncing"
	RequestEth2ConfigPath = "/eth/v1/config/spec"
	// RequestEth2DepositContractMethod       = "/eth/v1/config/deposit_contract"
	RequestGenesisPath = "/eth/v1/beacon/genesis"
	// RequestCommitteePath                   = "/eth/v1/beacon/states/%s/committees"
	// RequestFinalityCheckpointsPath         = "/eth/v1/beacon/states/%s/finality_checkpoints"
	// RequestForkPath                        = "/eth/v1/beacon/states/%s/fork"
	// RequestValidatorsPath                  = "/eth/v1/beacon/states/%s/validators"
	// RequestVoluntaryExitPath               = "/eth/v1/beacon/pool/voluntary_exits"
	// RequestAttestationsPath                = "/eth/v1/beacon/blocks/%s/attestations"
	// RequestBeaconBlockPath                 = "/eth/v2/beacon/blocks/%s"
	// RequestValidatorSyncDuties             = "/eth/v1/validator/duties/sync/%s"
	// RequestValidatorProposerDuties         = "/eth/v1/validator/duties/proposer/%s"
	// RequestWithdrawalCredentialsChangePath = "/eth/v1/beacon/pool/bls_to_execution_changes"
	// MaxRequestValidatorsCount     = 600
	// threadLimit               int = 12
)

// Unsigned integer type
type uinteger uint64

// Byte array type
type byteArray []byte

type Eth2ConfigResponse struct {
	Data struct {
		SecondsPerSlot               uinteger `json:"SECONDS_PER_SLOT"`
		SlotsPerEpoch                uinteger `json:"SLOTS_PER_EPOCH"`
		EpochsPerSyncCommitteePeriod uinteger `json:"EPOCHS_PER_SYNC_COMMITTEE_PERIOD"`
	} `json:"data"`
}

type GenesisResponse struct {
	Data struct {
		GenesisTime           uinteger  `json:"genesis_time"`
		GenesisForkVersion    byteArray `json:"genesis_fork_version"`
		GenesisValidatorsRoot byteArray `json:"genesis_validators_root"`
	} `json:"data"`
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

type BeaconClientInterface interface {
	GetEth2Config() (Eth2Config, error)
}

type BeaconClientManager struct {
	PrimaryBc    BeaconClientInterface
	PrimaryReady bool
}

type SimpleBeaconClient struct {
	providerAddress string
}

// TODO: Ask Joe - I'm not sure what this comment means
// Make a GET request but do not read its body yet (allows buffered decoding)
func (sbc *SimpleBeaconClient) getRequestReader(requestPath string) (io.ReadCloser, int, error) {
	fmt.Printf("RequestUrlFormat: %s\n", RequestUrlFormat)
	fmt.Printf("sbc.providerAddress: %s\n", sbc.providerAddress)
	fmt.Printf("requestPath: %s\n", requestPath)
	// Send request
	response, err := http.Get(fmt.Sprintf(RequestUrlFormat, sbc.providerAddress, requestPath))
	if err != nil {
		return nil, 0, err
	}

	return response.Body, response.StatusCode, nil
}

// Make a GET request to the beacon node and read the body of the response
func (sbc *SimpleBeaconClient) getRequest(requestPath string) ([]byte, int, error) {

	// Send request
	reader, status, err := sbc.getRequestReader(requestPath)
	if err != nil {
		return []byte{}, 0, err
	}
	defer func() {
		_ = reader.Close()
	}()

	// Get response
	body, err := io.ReadAll(reader)
	if err != nil {
		return []byte{}, 0, err
	}

	// Return
	return body, status, nil
}

func (sbc *SimpleBeaconClient) getEth2Config() (Eth2ConfigResponse, error) {

	responseBody, status, err := sbc.getRequest(RequestEth2ConfigPath)
	if err != nil {
		return Eth2ConfigResponse{}, fmt.Errorf("could not get eth2 config: %w", err)
	}
	if status != http.StatusOK {
		return Eth2ConfigResponse{}, fmt.Errorf("could not get eth2 config: HTTP status %d; response body: '%s'", status, string(responseBody))
	}
	var eth2Config Eth2ConfigResponse
	if err := json.Unmarshal(responseBody, &eth2Config); err != nil {
		return Eth2ConfigResponse{}, fmt.Errorf("could not decode eth2 config: %w", err)
	}
	return eth2Config, nil
}

func (sbc *SimpleBeaconClient) getGenesis() (GenesisResponse, error) {
	responseBody, status, err := sbc.getRequest(RequestGenesisPath)
	if err != nil {
		return GenesisResponse{}, fmt.Errorf("could not get genesis data: %w", err)
	}
	if status != http.StatusOK {
		return GenesisResponse{}, fmt.Errorf("could not get genesis data: HTTP status %d; response body: '%s'", status, string(responseBody))
	}
	var genesis GenesisResponse
	if err := json.Unmarshal(responseBody, &genesis); err != nil {
		return GenesisResponse{}, fmt.Errorf("could not decode genesis: %w", err)
	}
	return genesis, nil
}

// This is a signature for a wrapped Beacon client function that returns 1 var and an error
// type bcFunction1 func(BeaconClientInterface) (interface{}, error)

// This is a signature for a wrapped Beacon client function that returns 2 vars and an error
// type bcFunction2 func(BeaconClientInterface) (interface{}, interface{}, error)

// Creates a new BeaconClientManager instance based on the Rocket Pool config
func NewBeaconClientManager(providerAddress string) (*BeaconClientManager, error) {

	var primaryBc BeaconClientInterface = &SimpleBeaconClient{
		providerAddress: providerAddress,
	}

	return &BeaconClientManager{
		PrimaryBc:    primaryBc,
		PrimaryReady: true,
	}, nil
}

// Get the Beacon configuration
func (sbc *SimpleBeaconClient) GetEth2Config() (Eth2Config, error) {
	// Data
	var wg errgroup.Group
	var eth2Config Eth2ConfigResponse
	var genesis GenesisResponse

	// Get eth2 config
	// TODO: Just mock eth2Config for now until Beacon Client is ready
	eth2Config.Data.SecondsPerSlot = 12
	eth2Config.Data.SlotsPerEpoch = 32
	eth2Config.Data.EpochsPerSyncCommitteePeriod = 256
	// wg.Go(func() error {
	// 	var err error
	// 	// TODO: Setup Proteus to get Beacon Client node
	// 	eth2Config, err = sbc.getEth2Config()
	// 	return err
	// })

	// Get genesis
	// TODO: Just mock getGenesis for now until Beacon Client is ready
	genesis.Data.GenesisTime = 1606824000
	genesis.Data.GenesisForkVersion = []byte{0x00, 0x00, 0x00, 0x00}
	genesis.Data.GenesisValidatorsRoot = []byte{0x00, 0x00, 0x00, 0x00}
	// wg.Go(func() error {
	// 	var err error
	// 	genesis, err = sbc.getGenesis()
	// 	return err
	// })

	// Wait for data
	if err := wg.Wait(); err != nil {
		return Eth2Config{}, err
	}

	// Return response
	return Eth2Config{
		GenesisForkVersion:           genesis.Data.GenesisForkVersion,
		GenesisValidatorsRoot:        genesis.Data.GenesisValidatorsRoot,
		GenesisEpoch:                 0,
		GenesisTime:                  uint64(genesis.Data.GenesisTime),
		SecondsPerSlot:               uint64(eth2Config.Data.SecondsPerSlot),
		SlotsPerEpoch:                uint64(eth2Config.Data.SlotsPerEpoch),
		SecondsPerEpoch:              uint64(eth2Config.Data.SecondsPerSlot * eth2Config.Data.SlotsPerEpoch),
		EpochsPerSyncCommitteePeriod: uint64(eth2Config.Data.EpochsPerSyncCommitteePeriod),
	}, nil
}
