package service

import (
	"fmt"

	hdconfig "github.com/nodeset-org/hyperdrive-daemon/shared/config"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

// Generate a YAML file that shows the current configuration schema, including all of the parameters and their descriptions
func getConfigYaml(c *cli.Context) error {
	cfg, err := hdconfig.NewHyperdriveConfig("", []*hdconfig.HyperdriveSettings{})
	if err != nil {
		return fmt.Errorf("error creating Hyperdrive configuration: %w", err)
	}
	bytes, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("error serializing configuration schema: %w", err)
	}

	fmt.Println(string(bytes))
	return nil
}
