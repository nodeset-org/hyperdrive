package config

import (
	"fmt"
	"net"

	externalip "github.com/glendc/go-external-ip"
	"github.com/nodeset-org/hyperdrive/shared/types"
)

const (
	// Common param IDs across configs
	MaxPeersID        string = "maxPeers"
	ContainerTagID    string = "containerTag"
	AdditionalFlagsID string = "additionalFlags"
	HttpPortID        string = "httpPort"
	OpenHttpPortsID   string = "openHttpPort"
	P2pPortID         string = "p2pPort"
	PortID            string = "port"
	OpenPortID        string = "openPort"
	HttpUrlID         string = "httpUrl"
	EcID              string = "executionClient"
	BnID              string = "beaconNode"
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

// Get the external IP address. Try finding an IPv4 address first to:
// * Improve peer discovery and node performance
// * Avoid unnecessary container restarts caused by switching between IPv4 and IPv6
func getExternalIP() (net.IP, error) {
	// Try IPv4 first
	ip4Consensus := externalip.DefaultConsensus(nil, nil)
	ip4Consensus.UseIPProtocol(4)
	if ip, err := ip4Consensus.ExternalIP(); err == nil {
		return ip, nil
	}

	// Try IPv6 as fallback
	ip6Consensus := externalip.DefaultConsensus(nil, nil)
	ip6Consensus.UseIPProtocol(6)
	return ip6Consensus.ExternalIP()
}

// Check a port setting to see if it's already used elsewhere
func addAndCheckForDuplicate(portMap map[uint16]bool, param types.Parameter[uint16], errors []string) (map[uint16]bool, []string) {
	port := param.Value
	if portMap[port] {
		return portMap, append(errors, fmt.Sprintf("Port %s for %s is already in use", port, param.GetCommon().Name))
	} else {
		portMap[port] = true
	}
	return portMap, errors
}
