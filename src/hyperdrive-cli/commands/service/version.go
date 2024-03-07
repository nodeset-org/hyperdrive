package service

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/rocket-pool/node-manager-core/config"
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
	err = utils.PrintNetwork(cfg.Hyperdrive.Network.Value, isNew)
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
	clientMode := cfg.Hyperdrive.ClientMode.Value
	switch clientMode {
	case config.ClientMode_Local:
		format := "%s (Locally managed)\n\tImage: %s"

		// Execution client
		ec := cfg.Hyperdrive.LocalExecutionConfig.ExecutionClient.Value
		switch ec {
		case config.ExecutionClient_Geth:
			executionClientString = fmt.Sprintf(format, "Geth", cfg.Hyperdrive.LocalExecutionConfig.Geth.ContainerTag.Value)
		case config.ExecutionClient_Nethermind:
			executionClientString = fmt.Sprintf(format, "Nethermind", cfg.Hyperdrive.LocalExecutionConfig.Nethermind.ContainerTag.Value)
		case config.ExecutionClient_Besu:
			executionClientString = fmt.Sprintf(format, "Besu", cfg.Hyperdrive.LocalExecutionConfig.Besu.ContainerTag.Value)
		default:
			return fmt.Errorf("unknown local execution client [%v]", ec)
		}

		// Beacon node
		bn := cfg.Hyperdrive.LocalBeaconConfig.BeaconNode.Value
		switch bn {
		case config.BeaconNode_Lighthouse:
			beaconNodeString = fmt.Sprintf(format, "Lighthouse", cfg.Hyperdrive.LocalBeaconConfig.Lighthouse.ContainerTag.Value)
		case config.BeaconNode_Lodestar:
			beaconNodeString = fmt.Sprintf(format, "Lodestar", cfg.Hyperdrive.LocalBeaconConfig.Lodestar.ContainerTag.Value)
		case config.BeaconNode_Nimbus:
			beaconNodeString = fmt.Sprintf(format, "Nimbus", cfg.Hyperdrive.LocalBeaconConfig.Nimbus.ContainerTag.Value)
		case config.BeaconNode_Prysm:
			beaconNodeString = fmt.Sprintf(format, "Prysm", cfg.Hyperdrive.LocalBeaconConfig.Prysm.ContainerTag.Value)
		case config.BeaconNode_Teku:
			beaconNodeString = fmt.Sprintf(format, "Teku", cfg.Hyperdrive.LocalBeaconConfig.Teku.ContainerTag.Value)
		default:
			return fmt.Errorf("unknown local Beacon Node [%v]", bn)
		}

	case config.ClientMode_External:
		format := "Externally managed (%s)"

		// Execution client
		ec := cfg.Hyperdrive.ExternalExecutionConfig.ExecutionClient.Value
		switch ec {
		case config.ExecutionClient_Geth:
			executionClientString = fmt.Sprintf(format, "Geth")
		case config.ExecutionClient_Nethermind:
			executionClientString = fmt.Sprintf(format, "Nethermind")
		case config.ExecutionClient_Besu:
			executionClientString = fmt.Sprintf(format, "Besu")
		default:
			return fmt.Errorf("unknown external Execution Client [%v]", ec)
		}

		// Beacon node
		bn := cfg.Hyperdrive.ExternalBeaconConfig.BeaconNode.Value
		switch bn {
		case config.BeaconNode_Lighthouse:
			beaconNodeString = fmt.Sprintf(format, "Lighthouse")
		case config.BeaconNode_Lodestar:
			beaconNodeString = fmt.Sprintf(format, "Lodestar")
		case config.BeaconNode_Nimbus:
			beaconNodeString = fmt.Sprintf(format, "Nimbus")
		case config.BeaconNode_Prysm:
			beaconNodeString = fmt.Sprintf(format, "Prysm")
		case config.BeaconNode_Teku:
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
