package modules

import (
	"fmt"
	"regexp"

	"github.com/blang/semver/v4"
)

const (
	// The regex pattern for a standard descriptor string
	DescriptorStringPattern string = `^([a-zA-Z0-9-_\.]+$`
)

var (
	// The regex for a standard descriptor string
	descriptorStringRegex = regexp.MustCompile(DescriptorStringPattern)
)

// A standard descriptor string
type DescriptorString string

// Marshal the descriptor string to text
func (d DescriptorString) MarshalText() ([]byte, error) {
	return []byte(d), nil
}

// Unmarshal the descriptor string from text
func (d *DescriptorString) UnmarshalText(text []byte) error {
	if !descriptorStringRegex.Match(text) {
		return fmt.Errorf("invalid descriptor string: %s", text)
	}
	*d = DescriptorString(text)
	return nil
}

// Descriptor for a Hyperdrive module
type HyperdriveModuleDescriptor struct {
	// The name of the module
	Name DescriptorString `json:"name" yaml:"name"`

	// A shortcut to access the module's CLI in the terminal
	Shortcut DescriptorString `json:"shortcut" yaml:"shortcut"`

	// A description of the module
	Description DescriptorString `json:"description" yaml:"description"`

	// The version of the module
	Version semver.Version `json:"version,omitempty" yaml:"version,omitempty"`

	// The author of the module
	Author DescriptorString `json:"author" yaml:"author"`

	// A URL for more information about the module (repository, docs, etc).
	URL string `json:"url,omitempty" yaml:"url,omitempty"`

	// An optional email to use for contacting the authors
	Email string `json:"email,omitempty" yaml:"email,omitempty"`

	// Info about the module's CLI
	CLI CliInfo `json:"cli" yaml:"cli"`

	// A list of the module's dependencies
	Dependencies []Dependency `json:"dependencies,omitempty" yaml:"dependencies,omitempty"`
}

// Info about a module's CLI binary
type CliInfo struct {
	// The name of the CLI binary
	Filename string `json:"filename" yaml:"filename"`

	// The name of the Docker image tag if the CLI is included in a Docker image.
	// Ignored for standalone binaries installed on the native system.
	DockerTag string `json:"dockerTag,omitempty" yaml:"dockerTag,omitempty"`

	// The path to the CLI binary within the Docker image.
	// Ignored for standalone binaries installed on the native system.
	ImagePath string `json:"imagePath,omitempty" yaml:"imagePath,omitempty"`
}
