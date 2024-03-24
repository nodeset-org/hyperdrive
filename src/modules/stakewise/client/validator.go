package swclient

import (
	"strconv"

	swapi "github.com/nodeset-org/hyperdrive/modules/stakewise/shared/api"
	"github.com/rocket-pool/node-manager-core/api/client"
	"github.com/rocket-pool/node-manager-core/api/types"
	"github.com/rocket-pool/node-manager-core/beacon"
)

type ValidatorRequester struct {
	context *client.RequesterContext
}

func NewValidatorRequester(context *client.RequesterContext) *ValidatorRequester {
	return &ValidatorRequester{
		context: context,
	}
}

func (r *ValidatorRequester) GetName() string {
	return "Validator"
}
func (r *ValidatorRequester) GetRoute() string {
	return "validator"
}
func (r *ValidatorRequester) GetContext() *client.RequesterContext {
	return r.context
}

// Exit the provided validators from the Beacon Chain (or simply return their signed exit messages for later use without broadcasting),
// with an optional epoch parameter. If not specified, the epoch from the current chain head will be used.
func (r *ValidatorRequester) Exit(pubkeys []beacon.ValidatorPubkey, epoch *uint64, noBroadcastBool bool) (*types.ApiResponse[swapi.ValidatorExitData], error) {
	args := map[string]string{
		"pubkeys":      client.MakeBatchArg(pubkeys),
		"no-broadcast": strconv.FormatBool(noBroadcastBool),
	}
	if epoch != nil {
		args["epoch"] = strconv.FormatUint(*epoch, 10)
	}
	return client.SendGetRequest[swapi.ValidatorExitData](r, "exit", "Exit", args)
}
