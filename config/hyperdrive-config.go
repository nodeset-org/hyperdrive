package config

import (
	"fmt"
	"os"
	"path/filepath"

	"al.essio.dev/pkg/shellescape"
	"github.com/nodeset-org/hyperdrive/config/ids"
	"github.com/nodeset-org/hyperdrive/modules/config"
	"github.com/nodeset-org/hyperdrive/shared"
	"gopkg.in/yaml.v3"
)

const (
	// Tags
	hyperdriveTag string = "nodeset/hyperdrive:v" + shared.HyperdriveVersion

	// Defaults
	DefaultProjectName   string = "hyperdrive"
	DefaultApiPort       uint16 = 8080
	DefaultEnableIPv6    bool   = false
	DefaultClientTimeout uint16 = 30
)

// The base configuration for Hyperdrive
type HyperdriveConfig struct {
	ProjectName              config.StringParameter
	ApiPort                  config.UintParameter
	EnableIPv6               config.BoolParameter
	UserDataPath             config.StringParameter
	AdditionalDockerNetworks config.StringParameter
	ClientTimeout            config.UintParameter

	// The Docker Hub tag for the daemon container
	ContainerTag config.StringParameter

	// Logging
	Logging *LoggingConfig

	// Info about the loaded modules
	Modules map[string]*config.ModuleInfo

	// Internal fields
	Version                 string
	hyperdriveUserDirectory string
	systemDirectory         string
}

// Creates a new Hyperdrive configuration
func NewHyperdriveConfig(hdDir string, systemDir string) *HyperdriveConfig {
	defaultUserDataPath := filepath.Join(hdDir, "data")
	cfg := &HyperdriveConfig{
		hyperdriveUserDirectory: hdDir,
		systemDirectory:         systemDir,
		Modules:                 map[string]*config.ModuleInfo{},
		Version:                 shared.HyperdriveVersion,
	}

	// Project Name
	cfg.ProjectName.ID = config.Identifier(ids.ProjectNameID)
	cfg.ProjectName.Name = "Project Name"
	cfg.ProjectName.Description.Default = "This is the prefix that will be attached to all of the Docker containers managed by Hyperdrive."
	cfg.ProjectName.Default = DefaultProjectName
	cfg.ProjectName.AffectedContainers = []string{string(ContainerID_All)}

	// API Port
	cfg.ApiPort.ID = config.Identifier(ids.ApiPortID)
	cfg.ApiPort.Name = "Service API Port"
	cfg.ApiPort.Description.Default = "The port that Hyperdrive's API server should run on within the internal Docker network. Note this is bound to the local machine only; it cannot be accessed by other machines."
	cfg.ApiPort.Default = uint64(DefaultApiPort)
	cfg.ApiPort.AffectedContainers = []string{string(ContainerID_Daemon)}

	// Enable IPv6
	cfg.EnableIPv6.ID = config.Identifier(ids.EnableIPv6ID)
	cfg.EnableIPv6.Name = "Enable IPv6"
	cfg.EnableIPv6.Description.Default = "Enable IPv6 networking for Hyperdrive services. This is useful if you have an IPv6 network and want to use it for Hyperdrive.\n\nIf this isn't the first time you're starting Hyperdrive, you'll have to recreate the network after changing this box with `hyperdrive service down` and `hyperdrive service start` for it to take effect.\n\n[orange]NOTE: For IPv6 support to work, you must manually set up your Docker daemon to support it. Please follow the instructions at https://docs.docker.com/config/daemon/ipv6/#dynamic-ipv6-subnet-allocation before checking this box."
	cfg.EnableIPv6.Default = DefaultEnableIPv6
	cfg.EnableIPv6.AffectedContainers = []string{string(ContainerID_All)}

	// User Data Path
	cfg.UserDataPath.ID = config.Identifier(ids.UserDataPathID)
	cfg.UserDataPath.Name = "User Data Path"
	cfg.UserDataPath.Description.Default = "The absolute path of your personal `data` folder that contains secrets such as your node wallet's encrypted file, the password for your node wallet, and all of the validator keys for any Hyperdrive modules."
	cfg.UserDataPath.Default = defaultUserDataPath
	cfg.UserDataPath.AffectedContainers = []string{string(ContainerID_Daemon)}

	// Additional Docker Networks
	cfg.AdditionalDockerNetworks.ID = config.Identifier(ids.AdditionalDockerNetworksID)
	cfg.AdditionalDockerNetworks.Name = "Additional Docker Networks"
	cfg.AdditionalDockerNetworks.Description.Default = "List any other externally-managed Docker networks running on this machine that you'd like to give the Hyperdrive services access to here. Use a comma-separated list of network names.\n\nTo get a list of local Docker networks, run `docker network ls`."
	cfg.AdditionalDockerNetworks.AffectedContainers = []string{string(ContainerID_All)}

	// Client Timeout
	cfg.ClientTimeout.ID = config.Identifier(ids.ClientTimeoutID)
	cfg.ClientTimeout.Name = "Client Timeout"
	cfg.ClientTimeout.Description.Default = "The maximum time (in seconds) that Hyperdrive will wait for a response during HTTP requests (such as Execution Client, Beacon Node, or nodeset.io requests) before timing out."
	cfg.ClientTimeout.Default = uint64(DefaultClientTimeout)
	cfg.ClientTimeout.AffectedContainers = []string{string(ContainerID_Daemon)}

	// Container Tag
	cfg.ContainerTag.ID = config.Identifier(ids.ContainerTagID)
	cfg.ContainerTag.Name = "Service Container Tag"
	cfg.ContainerTag.Description.Default = "The tag name of the Hyperdrive Daemon image to use."
	cfg.ContainerTag.AffectedContainers = []string{string(ContainerID_Daemon)}
	cfg.ContainerTag.OverwriteOnUpgrade = true
	cfg.ContainerTag.Default = hyperdriveTag

	// Create the subconfigs
	cfg.Logging = NewLoggingConfig()
	return cfg
}

// Get the config.Parameters for this config
func (cfg *HyperdriveConfig) GetParameters() []config.IParameter {
	return []config.IParameter{
		&cfg.ProjectName,
		&cfg.ApiPort,
		&cfg.EnableIPv6,
		&cfg.UserDataPath,
		&cfg.AdditionalDockerNetworks,
		&cfg.ClientTimeout,
		&cfg.ContainerTag,
	}
}

// Get the subconfigurations for this config
func (cfg *HyperdriveConfig) GetSections() []config.ISection {
	return []config.ISection{
		cfg.Logging,
	}
}

func (cfg *HyperdriveConfig) GetModulePath() string {
	return shared.GetModulesDirectoryPath(cfg.systemDirectory)
}

func (cfg *HyperdriveConfig) GetAdapterKeyPath() string {
	return shared.GetAdapterKeyPath(cfg.hyperdriveUserDirectory)
}

func (cfg *HyperdriveConfig) GetVersion() string {
	return shared.HyperdriveVersion
}

// Load the Hyperdrive configuration from a file; the Hyperdrive user directory will be set to the directory containing the config file.
// In order to process module configuration instances,
func (c *HyperdriveConfig) LoadSettingsFromFile(configFilePath string) (*HyperdriveSettings, error) {
	// Return nil if the file doesn't exist
	_, err := os.Stat(configFilePath)
	if os.IsNotExist(err) {
		return nil, nil
	}

	// Read the file
	configBytes, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("could not read Hyperdrive settings file at %s: %w", shellescape.Quote(configFilePath), err)
	}

	// Attempt to parse it out into a config instance
	cfg := new(HyperdriveSettings)
	if err := yaml.Unmarshal(configBytes, cfg); err != nil {
		return nil, fmt.Errorf("could not parse config file: %w", err)
	}
	return cfg, nil
}
