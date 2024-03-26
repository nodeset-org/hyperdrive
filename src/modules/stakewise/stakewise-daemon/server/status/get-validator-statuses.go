package swstatus

import (
	"errors"
	"fmt"
	"net/url"

	swapi "github.com/nodeset-org/hyperdrive/modules/stakewise/shared/api"
	swtypes "github.com/nodeset-org/hyperdrive/modules/stakewise/shared/types"
	"github.com/rocket-pool/node-manager-core/api/types"
	"github.com/rocket-pool/node-manager-core/beacon"

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

func (c *statusGetValidatorsStatusesContext) PrepareData(data *swapi.ValidatorStatusData, opts *bind.TransactOpts) (types.ResponseStatus, error) {
	sp := c.handler.serviceProvider
	bc := sp.GetBeaconClient()
	w := sp.GetWallet()
	nc := sp.GetNodesetClient()
	ctx := sp.GetContext()

	nodesetStatusResponse, err := nc.GetRegisteredValidators()
	if err != nil {
		return types.ResponseStatus_Error, fmt.Errorf("error getting nodeset statuses: %w", err)
	}

	privateKeys, err := w.GetAllPrivateKeys()
	if err != nil {
		return types.ResponseStatus_Error, fmt.Errorf("error getting private keys: %w", err)
	}

	publicKeys, err := w.DerivePubKeys(privateKeys)
	if err != nil {
		return types.ResponseStatus_Error, fmt.Errorf("error getting public keys: %w", err)
	}

	beaconStatusResponse, err := bc.GetValidatorStatuses(ctx, publicKeys, nil)
	if err != nil {
		return types.ResponseStatus_Error, fmt.Errorf("error getting validator statuses: %w", err)
	}

	registeredPubkeys := make([]beacon.ValidatorPubkey, 0)
	for _, pubkeyStatus := range nodesetStatusResponse {
		registeredPubkeys = append(registeredPubkeys, pubkeyStatus.Pubkey)
	}

	// Get status info for each pubkey
	data.States = make([]swapi.ValidatorStateInfo, len(publicKeys))
	for i, pubkey := range publicKeys {
		state := &data.States[i]
		state.Pubkey = pubkey

		// Beacon status
		status, exists := beaconStatusResponse[pubkey]
		if exists {
			state.BeaconStatus = status.Status
			state.Index = status.Index
		} else {
			state.BeaconStatus = ""
		}

		// NodeSet status
		switch {
		case isRegisteredToStakewise(pubkey, beaconStatusResponse):
			state.NodesetStatus = swtypes.NodesetStatus_RegisteredToStakewise
		case isUploadedStakewise(pubkey, beaconStatusResponse):
			state.NodesetStatus = swtypes.NodesetStatus_UploadedStakewise
		case isUploadedToNodeset(pubkey, registeredPubkeys):
			state.NodesetStatus = swtypes.NodesetStatus_UploadedToNodeset
		default:
			state.NodesetStatus = swtypes.NodesetStatus_Generated
		}
	}

	return types.ResponseStatus_Success, nil
}

func isRegisteredToStakewise(pubKey beacon.ValidatorPubkey, statuses map[beacon.ValidatorPubkey]beacon.ValidatorStatus) bool {
	// TODO: Implement
	return false
}

func isUploadedStakewise(pubKey beacon.ValidatorPubkey, statuses map[beacon.ValidatorPubkey]beacon.ValidatorStatus) bool {
	// TODO: Implement
	return false
}

func isUploadedToNodeset(pubKey beacon.ValidatorPubkey, registeredPubkeys []beacon.ValidatorPubkey) bool {
	for _, registeredPubKey := range registeredPubkeys {
		if registeredPubKey == pubKey {
			return true
		}
	}
	return false
}
