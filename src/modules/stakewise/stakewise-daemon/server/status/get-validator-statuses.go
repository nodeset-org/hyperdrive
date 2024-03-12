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
	registeredPubkeys, err := nc.GetRegisteredValidators()
	if err != nil {
		return fmt.Errorf("error getting registered validators: %w", err)
	}
	privateKeys, err := w.GetAllPrivateKeys()
	if err != nil {
		return fmt.Errorf("error getting private keys: %w", err)
	}
	publicKeys, err := w.DerivePubKeys(privateKeys)
	if err != nil {
		return fmt.Errorf("error getting public keys: %w", err)
	}
	statuses, err := bc.GetValidatorStatuses(context.Background(), publicKeys, nil)
	if err != nil {
		return fmt.Errorf("error getting validator statuses: %w", err)
	}

	beaconStatuses := make(map[beacon.ValidatorPubkey]types.ValidatorState)
	nodesetStatuses := make(map[beacon.ValidatorPubkey]swapi.NodesetStatus)

	for _, pubKey := range publicKeys {
		status, exists := statuses[pubKey]
		if exists {
			beaconStatuses[pubKey] = status.Status
		}
	}

	for _, pubKey := range publicKeys {
		switch {
		case IsRegisteredToStakewise(pubKey, statuses):
			nodesetStatuses[pubKey] = swapi.RegisteredToStakewise
		case IsUploadedStakewise(pubKey, statuses):
			nodesetStatuses[pubKey] = swapi.UploadedStakewise
		case IsUploadedToNodeset(pubKey, statuses, registeredPubkeys):
			nodesetStatuses[pubKey] = swapi.UploadedToNodeset
		default:
			nodesetStatuses[pubKey] = swapi.Generated
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
