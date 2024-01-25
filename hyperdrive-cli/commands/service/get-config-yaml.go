package service

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

// Generate a YAML file that shows the current configuration schema, including all of the parameters and their descriptions
func getConfigYaml(c *cli.Context) error {
	cfg := config.NewHyperdriveConfig("")
	bytes, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("error serializing configuration schema: %w", err)
	}

	fmt.Println(string(bytes))
	return nil
}
