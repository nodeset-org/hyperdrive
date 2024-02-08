package swclient

import (
	"strconv"

	"github.com/nodeset-org/eth-utils/beacon"
	"github.com/nodeset-org/hyperdrive/client"
	swapi "github.com/nodeset-org/hyperdrive/modules/stakewise/shared/api"
	"github.com/nodeset-org/hyperdrive/shared/types/api"
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
func (r *ValidatorRequester) GetSignedExitMessage(pubkeys []beacon.ValidatorPubkey, epoch *uint64) (*api.ApiResponse[swapi.ValidatorGetSignedExitMessagesData], error) {
	args := map[string]string{
		"pubkeys": client.MakeBatchArg(pubkeys),
	}
	if epoch != nil {
		args["epoch"] = strconv.FormatUint(*epoch, 10)
	}
	return client.SendGetRequest[swapi.ValidatorGetSignedExitMessagesData](r, "get-signed-exit-messages", "GetSignedExitMessage", args)
}
