package service

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/nodeset-org/hyperdrive/shared/types"
	"github.com/urfave/cli/v2"
)

// View the Hyperdrive service version information
func serviceVersion(c *cli.Context) error {
	// Get Hyperdrive client
	hd := client.NewHyperdriveClientFromCtx(c)

	// Get the config
	cfg, isNew, err := hd.LoadConfig()
	if err != nil {
		return fmt.Errorf("Error loading configuration: %w", err)
	}

	// Print what network we're on
	err = utils.PrintNetwork(cfg.Network.Value, isNew)
	if err != nil {
		return err
	}

	// Get Hyperdrive service version
	serviceVersion, err := hd.GetServiceVersion()
	if err != nil {
		return err
	}

	// Get the execution client string
	var executionClientString string
	var beaconNodeString string
	clientMode := cfg.ClientMode.Value
	switch clientMode {
	case types.ClientMode_Local:
		format := "%s (Locally managed)\n\tImage: %s"

		// Execution client
		ec := cfg.LocalExecutionConfig.ExecutionClient.Value
		switch ec {
		case types.ExecutionClient_Geth:
			executionClientString = fmt.Sprintf(format, "Geth", cfg.LocalExecutionConfig.Geth.ContainerTag.Value)
		case types.ExecutionClient_Nethermind:
			executionClientString = fmt.Sprintf(format, "Nethermind", cfg.LocalExecutionConfig.Nethermind.ContainerTag.Value)
		case types.ExecutionClient_Besu:
			executionClientString = fmt.Sprintf(format, "Besu", cfg.LocalExecutionConfig.Besu.ContainerTag.Value)
		default:
			return fmt.Errorf("unknown local execution client [%v]", ec)
		}

		// Beacon node
		bn := cfg.LocalBeaconConfig.BeaconNode.Value
		switch bn {
		case types.BeaconNode_Lighthouse:
			beaconNodeString = fmt.Sprintf(format, "Lighthouse", cfg.LocalBeaconConfig.Lighthouse.ContainerTag.Value)
		case types.BeaconNode_Lodestar:
			beaconNodeString = fmt.Sprintf(format, "Lodestar", cfg.LocalBeaconConfig.Lodestar.ContainerTag.Value)
		case types.BeaconNode_Nimbus:
			beaconNodeString = fmt.Sprintf(format, "Nimbus", cfg.LocalBeaconConfig.Nimbus.ContainerTag.Value)
		case types.BeaconNode_Prysm:
			beaconNodeString = fmt.Sprintf(format, "Prysm", cfg.LocalBeaconConfig.Prysm.ContainerTag.Value)
		case types.BeaconNode_Teku:
			beaconNodeString = fmt.Sprintf(format, "Teku", cfg.LocalBeaconConfig.Teku.ContainerTag.Value)
		default:
			return fmt.Errorf("unknown local Beacon Node [%v]", bn)
		}

	case types.ClientMode_External:
		format := "Externally managed (%s)"

		// Execution client
		ec := cfg.ExternalExecutionConfig.ExecutionClient.Value
		switch ec {
		case types.ExecutionClient_Geth:
			executionClientString = fmt.Sprintf(format, "Geth")
		case types.ExecutionClient_Nethermind:
			executionClientString = fmt.Sprintf(format, "Nethermind")
		case types.ExecutionClient_Besu:
			executionClientString = fmt.Sprintf(format, "Besu")
		default:
			return fmt.Errorf("unknown external Execution Client [%v]", ec)
		}

		// Beacon node
		bn := cfg.ExternalBeaconConfig.BeaconNode.Value
		switch bn {
		case types.BeaconNode_Lighthouse:
			beaconNodeString = fmt.Sprintf(format, "Lighthouse")
		case types.BeaconNode_Lodestar:
			beaconNodeString = fmt.Sprintf(format, "Lodestar")
		case types.BeaconNode_Nimbus:
			beaconNodeString = fmt.Sprintf(format, "Nimbus")
		case types.BeaconNode_Prysm:
			beaconNodeString = fmt.Sprintf(format, "Prysm")
		case types.BeaconNode_Teku:
			beaconNodeString = fmt.Sprintf(format, "Teku")
		default:
			return fmt.Errorf("unknown external Beacon Node [%v]", bn)
		}

	default:
		return fmt.Errorf("unknown client mode [%v]", clientMode)
	}

	// Print version info
	fmt.Printf("Hyperdrive client version: %s\n", c.App.Version)
	fmt.Printf("Hyperdrive daemon version: %s\n", serviceVersion)
	fmt.Printf("Selected Execution Client: %s\n", executionClientString)
	fmt.Printf("Selected Beacon Node: %s\n", beaconNodeString)
	return nil
}
