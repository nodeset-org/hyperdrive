package client

import (
	"path/filepath"

	csconfig "github.com/nodeset-org/hyperdrive-constellation/shared/config"
	hdconfig "github.com/nodeset-org/hyperdrive-daemon/shared/config"
	swconfig "github.com/nodeset-org/hyperdrive-stakewise/shared/config"
)

var (
	hdApiKeyRelPath     string = filepath.Join(hdconfig.SecretsDir, hdconfig.DaemonKeyFilename)
	moduleApiKeyRelPath string = filepath.Join(hdconfig.SecretsDir, hdconfig.ModulesName)
	swApiKeyRelPath     string = filepath.Join(moduleApiKeyRelPath, swconfig.ModuleName, hdconfig.DaemonKeyFilename)
	csApiKeyRelPath     string = filepath.Join(moduleApiKeyRelPath, csconfig.ModuleName, hdconfig.DaemonKeyFilename)
)
