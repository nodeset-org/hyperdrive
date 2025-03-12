package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const (
	DefaultNodesetApiUrl string = "https://nodeset.io/api"

	DefaultEncryptionPubkey string = "age1hs87f8pl369tel4xpprtmaxwhg2e3gynvz8q45h48hed4m4h4d0s6t8caa"
)

type HyperdriveResources struct {
	// The URL for the NodeSet API server
	NodeSetApiUrl string `yaml:"nodeSetApiUrl" json:"nodeSetApiUrl"`

	// The pubkey used to encrypt messages to nodeset.io
	EncryptionPubkey string `yaml:"encryptionPubkey" json:"encryptionPubkey"`
}

// Load network settings from a folder
func LoadResources(sourceDir string) ([]*HyperdriveResources, error) {
	// Make sure the folder exists
	_, err := os.Stat(sourceDir)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("network settings folder \"%s\" does not exist", sourceDir)
	}

	// Enumerate the dir
	files, err := os.ReadDir(sourceDir)
	if err != nil {
		return nil, fmt.Errorf("error enumerating network settings folder: %w", err)
	}

	settingsList := []*HyperdriveResources{}
	for _, file := range files {
		// Ignore dirs and nonstandard files
		if file.IsDir() || !file.Type().IsRegular() {
			continue
		}

		// Load the file
		filename := file.Name()
		ext := filepath.Ext(filename)
		if ext != ".yaml" && ext != ".yml" {
			// Only load YAML files
			continue
		}
		settingsFilePath := filepath.Join(sourceDir, filename)
		bytes, err := os.ReadFile(settingsFilePath)
		if err != nil {
			return nil, fmt.Errorf("error reading network settings file \"%s\": %w", settingsFilePath, err)
		}

		// Unmarshal the settings
		settings := new(HyperdriveResources)
		err = yaml.Unmarshal(bytes, settings)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling network settings file \"%s\": %w", settingsFilePath, err)
		}
		settingsList = append(settingsList, settings)
	}
	return settingsList, nil
}
