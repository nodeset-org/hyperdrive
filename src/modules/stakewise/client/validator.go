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

// Get signed exit messages for the provided validators, with an optional epoch parameter. If not specified, the epoch from the current chain head will be used.
func (r *ValidatorRequester) GetSignedExitMessage(pubkeys []beacon.ValidatorPubkey, epoch *uint64, noBroadcastBool bool) (*types.ApiResponse[swapi.ValidatorGetSignedExitMessagesData], error) {
	args := map[string]string{
		"pubkeys":      client.MakeBatchArg(pubkeys),
		"no-broadcast": strconv.FormatBool(noBroadcastBool),
	}
	if epoch != nil {
		args["epoch"] = strconv.FormatUint(*epoch, 10)
	}
	return client.SendGetRequest[swapi.ValidatorGetSignedExitMessagesData](r, "get-signed-exit-messages", "GetSignedExitMessage", args)
}
