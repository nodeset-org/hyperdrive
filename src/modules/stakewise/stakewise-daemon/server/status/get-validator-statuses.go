package swstatus

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/nodeset-org/eth-utils/beacon"
	swapi "github.com/nodeset-org/hyperdrive/modules/stakewise/shared/api"
	"github.com/nodeset-org/hyperdrive/shared/types"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/daemon-utils/server"
)

// ===============
// === Factory ===
// ===============

type statusGetValidatorsStatusesContextFactory struct {
	handler *StatusHandler
}

func (f *statusGetValidatorsStatusesContextFactory) Create(args url.Values) (*statusGetValidatorsStatusesContext, error) {
	c := &statusGetValidatorsStatusesContext{
		handler: f.handler,
	}
	inputErrs := []error{}
	return c, errors.Join(inputErrs...)
}

func (f *statusGetValidatorsStatusesContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterQuerylessGet[*statusGetValidatorsStatusesContext, swapi.ValidatorStatusData](
		router, "status", f, f.handler.serviceProvider.ServiceProvider,
	)
}

// ===============
// === Context ===
// ===============

type statusGetValidatorsStatusesContext struct {
	handler *StatusHandler
}

func (c *statusGetValidatorsStatusesContext) PrepareData(data *swapi.ValidatorStatusData, opts *bind.TransactOpts) error {
	sp := c.handler.serviceProvider
	bc := sp.GetBeaconClient()
	w := sp.GetWallet()
	nc := sp.GetNodesetClient()
	nodesetStatusResponse, err := nc.GetRegisteredValidators()
	if err != nil {
		return fmt.Errorf("error getting nodeset statuses: %w", err)
	}
	fmt.Printf("!!! Nodeset statuses: %v\n", nodesetStatusResponse)
	privateKeys, err := w.GetAllPrivateKeys()
	if err != nil {
		return fmt.Errorf("error getting private keys: %w", err)
	}
	publicKeys, err := w.DerivePubKeys(privateKeys)
	if err != nil {
		return fmt.Errorf("error getting public keys: %w", err)
	}
	statusResponse, err := bc.GetValidatorStatuses(context.Background(), publicKeys, nil)
	if err != nil {
		return fmt.Errorf("error getting validator statuses: %w", err)
	}

	registeredPubkeys := make([]beacon.ValidatorPubkey, 0)
	for _, pubkeyStatus := range nodesetStatusResponse {
		registeredPubkeys = append(registeredPubkeys, pubkeyStatus.Pubkey)
	}

	beaconStatuses := make(map[string]types.ValidatorState)
	nodesetStatuses := make(map[string]swapi.NodesetStatus)

	for _, pubKey := range publicKeys {
		status, exists := statusResponse[pubKey]
		if exists {
			beaconStatuses[pubKey.HexWithPrefix()] = status.Status
		}
	}

	for _, pubKey := range publicKeys {
		switch {
		case IsRegisteredToStakewise(pubKey, statusResponse):
			nodesetStatuses[pubKey.HexWithPrefix()] = swapi.RegisteredToStakewise
		case IsUploadedStakewise(pubKey, statusResponse):
			nodesetStatuses[pubKey.HexWithPrefix()] = swapi.UploadedStakewise
		case IsUploadedToNodeset(pubKey, statusResponse, registeredPubkeys):
			nodesetStatuses[pubKey.HexWithPrefix()] = swapi.UploadedToNodeset
		default:
			nodesetStatuses[pubKey.HexWithPrefix()] = swapi.Generated
		}
	}

	data.BeaconStatus = beaconStatuses
	data.NodesetStatus = nodesetStatuses

	return nil
}

func IsRegisteredToStakewise(pubKey beacon.ValidatorPubkey, statuses map[beacon.ValidatorPubkey]types.ValidatorStatus) bool {
	// TODO: Implement
	return false
}

func IsUploadedStakewise(pubKey beacon.ValidatorPubkey, statuses map[beacon.ValidatorPubkey]types.ValidatorStatus) bool {
	// TODO: Implement
	return false
}

func IsUploadedToNodeset(pubKey beacon.ValidatorPubkey, statuses map[beacon.ValidatorPubkey]types.ValidatorStatus, registeredPubkeys []beacon.ValidatorPubkey) bool {
	for _, registeredPubKey := range registeredPubkeys {
		if registeredPubKey == pubKey {
			return true
		}
	}
	return false
}
