package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/alessio/shellescape"
	"github.com/nodeset-org/hyperdrive/shared"
	"github.com/nodeset-org/hyperdrive/shared/config/ids"
	"github.com/nodeset-org/hyperdrive/shared/config/migration"
	"github.com/pbnjay/memory"
	"github.com/rocket-pool/node-manager-core/config"
	"gopkg.in/yaml.v3"
)

// =========================
// === Hyperdrive Config ===
// =========================

const (
	// Tags
	hyperdriveTag string = "nodeset/hyperdrive:v" + shared.HyperdriveVersion
)

// The master configuration struct
type HyperdriveConfig struct {
	// General settings
	DebugMode          config.Parameter[bool]
	Network            config.Parameter[config.Network]
	ClientMode         config.Parameter[config.ClientMode]
	ProjectName        config.Parameter[string]
	UserDataPath       config.Parameter[string]
	AutoTxMaxFee       config.Parameter[float64]
	MaxPriorityFee     config.Parameter[float64]
	AutoTxGasThreshold config.Parameter[float64]

	// Execution client settings
	LocalExecutionClient    *config.LocalExecutionConfig
	ExternalExecutionClient *config.ExternalExecutionConfig

	// Beacon node settings
	LocalBeaconClient    *config.LocalBeaconConfig
	ExternalBeaconClient *config.ExternalBeaconConfig

	// Fallback clients
	Fallback *config.FallbackConfig

	// Metrics
	Metrics *config.MetricsConfig

	// Modules
	Modules map[string]any

	// Internal fields
	Version                 string
	HyperdriveUserDirectory string
	resources               *config.NetworkResources
}

// Load configuration settings from a file
func LoadFromFile(path string) (*HyperdriveConfig, error) {
	// Return nil if the file doesn't exist
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, nil
	}

	// Read the file
	configBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read Hyperdrive settings file at %s: %w", shellescape.Quote(path), err)
	}

	// Attempt to parse it out into a settings map
	var settings map[string]any
	if err := yaml.Unmarshal(configBytes, &settings); err != nil {
		return nil, fmt.Errorf("could not parse settings file: %w", err)
	}

	// Deserialize it into a config object
	cfg := NewHyperdriveConfig(filepath.Dir(path))
	err = cfg.Deserialize(settings)
	if err != nil {
		return nil, fmt.Errorf("could not deserialize settings file: %w", err)
	}

	return cfg, nil
}

// Creates a new Hyperdrive configuration instance
func NewHyperdriveConfig(hdDir string) *HyperdriveConfig {
	cfg := &HyperdriveConfig{
		HyperdriveUserDirectory: hdDir,
		Modules:                 map[string]any{},

		ProjectName: config.Parameter[string]{
			ParameterCommon: &config.ParameterCommon{
				ID:                 ids.ProjectNameID,
				Name:               "Project Name",
				Description:        "This is the prefix that will be attached to all of the Docker containers managed by Hyperdrive.",
				AffectsContainers:  []config.ContainerID{config.ContainerID_BeaconNode, config.ContainerID_Daemon, config.ContainerID_ExecutionClient, config.ContainerID_Exporter, config.ContainerID_Grafana, config.ContainerID_Prometheus, config.ContainerID_ValidatorClient},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[config.Network]string{
				config.Network_All: "hyperdrive",
			},
		},

		Network: config.Parameter[config.Network]{
			ParameterCommon: &config.ParameterCommon{
				ID:                 ids.NetworkID,
				Name:               "Network",
				Description:        "The Ethereum network you want to use - select Holesky Testnet to practice with fake ETH, or Mainnet to stake on the real network using real ETH.",
				AffectsContainers:  []config.ContainerID{config.ContainerID_Daemon, config.ContainerID_ExecutionClient, config.ContainerID_BeaconNode, config.ContainerID_ValidatorClient},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Options: getNetworkOptions(),
			Default: map[config.Network]config.Network{
				config.Network_All: config.Network_Mainnet,
			},
		},

		ClientMode: config.Parameter[config.ClientMode]{
			ParameterCommon: &config.ParameterCommon{
				ID:                 ids.ClientModeID,
				Name:               "Client Mode",
				Description:        "Choose which mode to use for your Execution Client and Beacon Node - locally managed (Docker Mode), or externally managed (Hybrid Mode).",
				AffectsContainers:  []config.ContainerID{config.ContainerID_Daemon, config.ContainerID_ExecutionClient, config.ContainerID_BeaconNode},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Options: []*config.ParameterOption[config.ClientMode]{
				{
					ParameterOptionCommon: &config.ParameterOptionCommon{
						Name:        "Locally Managed",
						Description: "Allow Hyperdrive to manage the Execution Client and Beacon Node for you (Docker Mode)",
					},
					Value: config.ClientMode_Local,
				}, {
					ParameterOptionCommon: &config.ParameterOptionCommon{
						Name:        "Externally Managed",
						Description: "Use an existing Execution Client and Beacon Node that you manage on your own (Hybrid Mode)",
					},
					Value: config.ClientMode_External,
				}},
			Default: map[config.Network]config.ClientMode{
				config.Network_All: config.ClientMode_Local,
			},
		},

		AutoTxMaxFee: config.Parameter[float64]{
			ParameterCommon: &config.ParameterCommon{
				ID:                 ids.AutoTxMaxFeeID,
				Name:               "Auto TX Max Fee",
				Description:        "Set this if you want all of Hyperdrive's automatic transactions to use this specific max fee value (in gwei), which is the most you'd be willing to pay (*including the priority fee*).\n\nA value of 0 will use the suggested max fee based on the current network conditions.\n\nAny other value will ignore the network suggestion and use this value instead.",
				AffectsContainers:  []config.ContainerID{config.ContainerID_Daemon},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[config.Network]float64{
				config.Network_All: float64(0),
			},
		},

		MaxPriorityFee: config.Parameter[float64]{
			ParameterCommon: &config.ParameterCommon{
				ID:                 ids.MaxPriorityFeeID,
				Name:               "Max Priority Fee",
				Description:        "The default value for the priority fee (in gwei) for all of your transactions, including automatic ones. This describes how much you're willing to pay *above the network's current base fee* - the higher this is, the more ETH you give to the validators for including your transaction, which generally means it will be included in a block faster (as long as your max fee is sufficiently high to cover the current network conditions).\n\nMust be larger than 0.",
				AffectsContainers:  []config.ContainerID{config.ContainerID_Daemon},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[config.Network]float64{
				config.Network_All: float64(1),
			},
		},

		AutoTxGasThreshold: config.Parameter[float64]{
			ParameterCommon: &config.ParameterCommon{
				ID:                 ids.AutoTxGasThresholdID,
				Name:               "Automatic TX Gas Threshold",
				Description:        "The threshold (in gwei) that the recommended network gas price must be under in order for automated transactions to be submitted when due.\n\nA value of 0 will disable non-essential automatic transactions.",
				AffectsContainers:  []config.ContainerID{config.ContainerID_Daemon},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[config.Network]float64{
				config.Network_All: float64(100),
			},
		},

		UserDataPath: config.Parameter[string]{
			ParameterCommon: &config.ParameterCommon{
				ID:                 ids.UserDataPathID,
				Name:               "User Data Path",
				Description:        "The absolute path of your personal `data` folder that contains secrets such as your node wallet's encrypted file, the password for your node wallet, and all of the validator keys for any Hyperdrive modules.",
				AffectsContainers:  []config.ContainerID{config.ContainerID_Daemon, config.ContainerID_ValidatorClient},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[config.Network]string{
				config.Network_All: filepath.Join(hdDir, "data"),
			},
		},

		DebugMode: config.Parameter[bool]{
			ParameterCommon: &config.ParameterCommon{
				ID:                 ids.DebugModeID,
				Name:               "Debug Mode",
				Description:        "Enable debug log printing in the daemon.",
				AffectsContainers:  []config.ContainerID{config.ContainerID_Daemon},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[config.Network]bool{
				config.Network_All: false,
			},
		},
	}

	// Create the subconfigs
	cfg.LocalExecutionClient = NewLocalExecutionClient()
	cfg.ExternalExecutionClient = config.NewExternalExecutionConfig()
	cfg.LocalBeaconClient = NewLocalBeaconClient()
	cfg.ExternalBeaconClient = config.NewExternalBeaconConfig()
	cfg.Fallback = config.NewFallbackConfig()
	cfg.Metrics = NewMetricsConfig()

	// Apply the default values for mainnet
	cfg.Network.Value = config.Network_Mainnet
	cfg.applyAllDefaults()

	return cfg
}

// Get the title for this config
func (cfg *HyperdriveConfig) GetTitle() string {
	return "Hyperdrive"
}

// Get the config.Parameters for this config
func (cfg *HyperdriveConfig) GetParameters() []config.IParameter {
	return []config.IParameter{
		&cfg.ProjectName,
		&cfg.Network,
		&cfg.ClientMode,
		&cfg.AutoTxMaxFee,
		&cfg.MaxPriorityFee,
		&cfg.AutoTxGasThreshold,
		&cfg.UserDataPath,
		&cfg.DebugMode,
	}
}

// Get the subconfigurations for this config
func (cfg *HyperdriveConfig) GetSubconfigs() map[string]config.IConfigSection {
	return map[string]config.IConfigSection{
		ids.FallbackID:          cfg.Fallback,
		ids.LocalExecutionID:    cfg.LocalExecutionClient,
		ids.ExternalExecutionID: cfg.ExternalExecutionClient,
		ids.LocalBeaconID:       cfg.LocalBeaconClient,
		ids.ExternalBeaconID:    cfg.ExternalBeaconClient,
		ids.MetricsID:           cfg.Metrics,
	}
}

// Serializes the configuration into a map of maps, compatible with a settings file
func (cfg *HyperdriveConfig) Serialize(modules []IModuleConfig) map[string]any {
	masterMap := map[string]any{}

	hdMap := config.Serialize(cfg)
	masterMap[ids.UserDirectoryKey] = cfg.HyperdriveUserDirectory
	masterMap[ids.VersionID] = fmt.Sprintf("v%s", shared.HyperdriveVersion)
	masterMap[ids.RootConfigID] = hdMap

	// Handle modules
	modulesMap := map[string]any{}
	for modName, value := range cfg.Modules {
		// Copy the module configs already on-board
		modulesMap[modName] = value
	}
	for _, module := range modules {
		// Serialize / overwrite them with explictly provided ones
		modMap := config.Serialize(module)
		modulesMap[module.GetModuleName()] = modMap
	}
	masterMap[ModulesName] = modulesMap
	return masterMap
}

// Deserializes a settings file into this config
func (cfg *HyperdriveConfig) Deserialize(masterMap map[string]any) error {
	// Upgrade the config to the latest version
	err := migration.UpdateConfig(masterMap)
	if err != nil {
		return fmt.Errorf("error upgrading configuration to v%s: %w", shared.HyperdriveVersion, err)
	}

	// Get the network
	network := config.Network_Mainnet
	hyperdriveParams, exists := masterMap[ids.RootConfigID]
	if !exists {
		return fmt.Errorf("config is missing the [%s] section", ids.RootConfigID)
	}
	hdMap, isMap := hyperdriveParams.(map[string]any)
	if !isMap {
		return fmt.Errorf("config has an entry named [%s] but it is not a map, it's a %s", ids.RootConfigID, reflect.TypeOf(hyperdriveParams))
	}
	networkVal, exists := hdMap[cfg.Network.ID]
	if exists {
		networkString, isString := networkVal.(string)
		if !isString {
			return fmt.Errorf("expected [%s - %s] to be a string but it is not", ids.RootConfigID, cfg.Network.ID)
		}
		network = config.Network(networkString)
	}

	// Deserialize the params and subconfigs
	err = config.Deserialize(cfg, hdMap, network)
	if err != nil {
		return fmt.Errorf("error deserializing [%s]: %w", ids.RootConfigID, err)
	}

	// Get the special fields
	udKey, exists := masterMap[ids.UserDirectoryKey]
	if !exists {
		return fmt.Errorf("expected a user directory config.Parameter named [%s] but it was not found", ids.UserDirectoryKey)
	}
	cfg.HyperdriveUserDirectory = udKey.(string)
	version, exists := masterMap[ids.VersionID]
	if !exists {
		return fmt.Errorf("expected a version config.Parameter named [%s] but it was not found", ids.VersionID)
	}
	cfg.Version = version.(string)

	// Handle modules
	modules, exists := masterMap[ModulesName]
	if exists {
		if modMap, ok := modules.(map[string]any); ok {
			cfg.Modules = modMap
		} else {
			return fmt.Errorf("config has an entry named [%s] but it is not a map, it's a %s", ModulesName, reflect.TypeOf(modules))
		}
	} else {
		cfg.Modules = map[string]any{}
	}

	cfg.updateResources()
	return nil
}

// Changes the current network, propagating new parameter settings if they are affected
func (cfg *HyperdriveConfig) ChangeNetwork(newNetwork config.Network) {
	// Get the current network
	oldNetwork := cfg.Network.Value
	if oldNetwork == newNetwork {
		return
	}
	cfg.Network.Value = newNetwork

	// Run the changes
	config.ChangeNetwork(cfg, oldNetwork, newNetwork)
	cfg.updateResources()
}

// =====================
// === Field Helpers ===
// =====================

// Applies all of the defaults to all of the settings that have them defined
func (cfg *HyperdriveConfig) applyAllDefaults() {
	network := cfg.Network.Value
	config.ApplyDefaults(cfg, network)
}

// Get the list of options for networks to run on
func getNetworkOptions() []*config.ParameterOption[config.Network] {
	options := []*config.ParameterOption[config.Network]{
		/*{
			ParameterOptionCommon: &config.ParameterOptionCommon{
				Name:        "Ethereum Mainnet",
				Description: "This is the real Ethereum main network, using real ETH to make real validators.",
			},
			Value: config.Network_Mainnet,
		},*/
		{
			ParameterOptionCommon: &config.ParameterOptionCommon{
				Name:        "Holesky Testnet",
				Description: "This is the Holešky (Holešovice) test network, which is the next generation of long-lived testnets for Ethereum. It uses free fake ETH to make fake validators.\nUse this if you want to practice running Hyperdrive in a free, safe environment before moving to Mainnet.",
			},
			Value: config.Network_Holesky,
		},
	}

	if strings.HasSuffix(shared.HyperdriveVersion, "-dev") {
		options = append(options, &config.ParameterOption[config.Network]{
			ParameterOptionCommon: &config.ParameterOptionCommon{
				Name:        "Devnet",
				Description: "This is a development network used by Hyperdrive engineers to test new features and contract upgrades before they are promoted to Holesky for staging. You should not use this network unless invited to do so by the developers.",
			},
			Value: Network_HoleskyDev,
		})
	}

	return options
}

// Get a more verbose client description, including warnings
func getAugmentedEcDescription(client config.ExecutionClient, originalDescription string) string {
	switch client {
	case config.ExecutionClient_Nethermind:
		totalMemoryGB := memory.TotalMemory() / 1024 / 1024 / 1024
		if totalMemoryGB < 9 {
			return fmt.Sprintf("%s\n\n[red]WARNING: Nethermind currently requires over 8 GB of RAM to run smoothly. We do not recommend it for your system. This may be improved in a future release.", originalDescription)
		}
	case config.ExecutionClient_Besu:
		totalMemoryGB := memory.TotalMemory() / 1024 / 1024 / 1024
		if totalMemoryGB < 9 {
			return fmt.Sprintf("%s\n\n[red]WARNING: Besu currently requires over 8 GB of RAM to run smoothly. We do not recommend it for your system. This may be improved in a future release.", originalDescription)
		}
	}

	return originalDescription
}

func (cfg *HyperdriveConfig) updateResources() {
	switch cfg.Network.Value {
	case Network_HoleskyDev:
		cfg.resources = config.NewResources(config.Network_Holesky)
	default:
		cfg.resources = config.NewResources(cfg.Network.Value)
	}
}

// ==============================
// === IConfig Implementation ===
// ==============================

func (cfg *HyperdriveConfig) GetApiLogFilePath() string {
	return filepath.Join(cfg.HyperdriveUserDirectory, LogDir, ApiLogName)
}

func (cfg *HyperdriveConfig) GetTasksLogFilePath() string {
	return filepath.Join(cfg.HyperdriveUserDirectory, LogDir, TasksLogName)
}

func (cfg *HyperdriveConfig) GetNodeAddressFilePath() string {
	return filepath.Join(cfg.UserDataPath.Value, UserAddressFilename)
}

func (cfg *HyperdriveConfig) GetWalletFilePath() string {
	return filepath.Join(cfg.UserDataPath.Value, UserWalletDataFilename)
}

func (cfg *HyperdriveConfig) GetPasswordFilePath() string {
	return filepath.Join(cfg.UserDataPath.Value, UserPasswordFilename)
}

func (cfg *HyperdriveConfig) GetNetworkResources() *config.NetworkResources {
	return cfg.resources
}

func (cfg *HyperdriveConfig) GetExecutionClientUrls() (string, string) {
	primaryEcUrl := cfg.GetEcHttpEndpoint()
	var fallbackEcUrl string
	if cfg.Fallback.UseFallbackClients.Value {
		fallbackEcUrl = cfg.Fallback.EcHttpUrl.Value
	}
	return primaryEcUrl, fallbackEcUrl
}

func (cfg *HyperdriveConfig) GetBeaconNodeUrls() (string, string) {
	primaryBnUrl := cfg.GetBnHttpEndpoint()
	var fallbackBnUrl string
	if cfg.Fallback.UseFallbackClients.Value {
		fallbackBnUrl = cfg.Fallback.BnHttpUrl.Value
	}
	return primaryBnUrl, fallbackBnUrl
}
