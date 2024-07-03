package client

import (
	"github.com/nodeset-org/hyperdrive-daemon/shared/types/api"
	"github.com/rocket-pool/node-manager-core/api/client"
	"github.com/rocket-pool/node-manager-core/api/types"
)

// Requester for core calls to the nodeset.io service
type NodeSetRequester struct {
	context client.IRequesterContext
}

func NewNodeSetRequester(context client.IRequesterContext) *NodeSetRequester {
	return &NodeSetRequester{
		context: context,
	}
}

func (r *NodeSetRequester) GetName() string {
	return "NodeSet"
}
func (r *NodeSetRequester) GetRoute() string {
	return "nodeset"
}
func (r *NodeSetRequester) GetContext() client.IRequesterContext {
	return r.context
}

// Gets the node's registration status with the NodeSet service
func (r *ServiceRequester) GetRegistrationStatus() (*types.ApiResponse[api.NodeSetGetRegistrationStatusData], error) {
	return client.SendGetRequest[api.NodeSetGetRegistrationStatusData](r, "get-registration-status", "GetRegistrationStatus", nil)
}

// Registers the node with the NodeSet service
func (r *ServiceRequester) RegisterNode(email string) (*types.ApiResponse[api.NodeSetRegisterNodeData], error) {
	args := map[string]string{
		"email": email,
	}
	return client.SendGetRequest[api.NodeSetRegisterNodeData](r, "register-node", "RegisterNode", args)
}
