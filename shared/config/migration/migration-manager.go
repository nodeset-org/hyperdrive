package migration

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/nodeset-org/hyperdrive/shared/config/ids"
)

type ConfigUpgrader struct {
	Version     *version.Version
	UpgradeFunc func(serializedConfig map[string]any) error
}

func UpdateConfig(serializedConfig map[string]any) error {
	// Get the config's version
	configVersion, err := getVersionFromConfig(serializedConfig)
	if err != nil {
		return err
	}

	// Create the collection of upgraders
	upgraders := []ConfigUpgrader{}

	// Find the index of the provided config's version
	targetIndex := -1
	for i, upgrader := range upgraders {
		if configVersion.LessThanOrEqual(upgrader.Version) {
			targetIndex = i
		}
	}

	// If there are no upgrades to apply, return
	if targetIndex == -1 {
		return nil
	}

	// If there are upgrades, start at the first applicable index and apply them all in series
	for i := targetIndex; i < len(upgraders); i++ {
		upgrader := upgraders[i]
		err = upgrader.UpgradeFunc(serializedConfig)
		if err != nil {
			return fmt.Errorf("error applying upgrade for config version %s: %w", upgrader.Version.String(), err)
		}
	}

	return nil

}

// Get the Smartnode version that the given config was built with
func getVersionFromConfig(serializedConfig map[string]any) (*version.Version, error) {
	configVersionString, exists := serializedConfig[ids.VersionID].(string)
	if !exists {
		return nil, fmt.Errorf("expected a top-level setting named '%s' but it didn't exist", ids.VersionID)
	}

	configVersion, err := version.NewVersion(strings.TrimPrefix(configVersionString, "v"))
	if err != nil {
		return nil, fmt.Errorf("error parsing version [%s] from config file: %w", configVersionString, err)
	}

	return configVersion, nil
}

// Parses a version string into a semantic version
func parseVersion(versionString string) (*version.Version, error) {
	parsedVersion, err := version.NewSemver(versionString)
	if err != nil {
		return nil, fmt.Errorf("error parsing version %s: %w", versionString, err)
	}
	return parsedVersion, nil
}
