package config

import (
	"net"

	externalip "github.com/glendc/go-external-ip"
)

// Get the possible RPC port mode options
func getPortModes(warningOverride string) []*ParameterOption[RpcPortMode] {
	if warningOverride == "" {
		warningOverride = "Allow connections from external hosts. This is safe if you're running your node on your local network. If you're a VPS user, this would expose your node to the internet"
	}

	return []*ParameterOption[RpcPortMode]{
		{
			ParameterOptionCommon: &ParameterOptionCommon{
				Name:        "Closed",
				Description: "Do not allow connections to the port",
			},
			Value: RpcPortMode_Closed,
		}, {
			ParameterOptionCommon: &ParameterOptionCommon{
				Name:        "Open to Localhost",
				Description: "Allow connections from this host only",
			},
			Value: RpcPortMode_OpenLocalhost,
		}, {
			ParameterOptionCommon: &ParameterOptionCommon{
				Name:        "Open to External hosts",
				Description: warningOverride,
			},
			Value: RpcPortMode_OpenExternal,
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
