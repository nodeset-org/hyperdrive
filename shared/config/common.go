package config

import "github.com/nodeset-org/hyperdrive/shared/types"

const (
	// Common param IDs across configs
	MaxPeersID        string = "maxPeers"
	ContainerTagID    string = "containerTag"
	AdditionalFlagsID string = "additionalFlags"
	HttpPortID        string = "httpPort"
	OpenHttpPortID    string = "openHttpPort"
	P2pPortID         string = "p2pPort"
	PortID            string = "port"
	OpenPortID        string = "openPort"
	HttpUrlID         string = "httpUrl"
)

// Get the possible RPC port mode options
func getPortModes(warningOverride string) []*types.ParameterOption[types.RpcPortMode] {
	if warningOverride == "" {
		warningOverride = "Allow connections from external hosts. This is safe if you're running your node on your local network. If you're a VPS user, this would expose your node to the internet"
	}

	return []*types.ParameterOption[types.RpcPortMode]{
		{
			ParameterOptionCommon: &types.ParameterOptionCommon{
				Name:        "Closed",
				Description: "Do not allow connections to the port",
			},
			Value: types.RpcPortMode_Closed,
		}, {
			ParameterOptionCommon: &types.ParameterOptionCommon{
				Name:        "Open to Localhost",
				Description: "Allow connections from this host only",
			},
			Value: types.RpcPortMode_OpenLocalhost,
		}, {
			ParameterOptionCommon: &types.ParameterOptionCommon{
				Name:        "Open to External hosts",
				Description: warningOverride,
			},
			Value: types.RpcPortMode_OpenExternal,
		},
	}
}
