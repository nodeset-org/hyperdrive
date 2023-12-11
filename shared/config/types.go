package config

type ParameterType string
type ContainerID string

type Network string

// A single option in a choice parameter
type ParameterOption struct {
	Name        string      `yaml:"name,omitempty"`
	Description string      `yaml:"description,omitempty"`
	Value       interface{} `yaml:"value,omitempty"`
}

// A parameter that can be configured by the user
type Parameter struct {
	ID                    string                  `yaml:"id,omitempty"`
	Name                  string                  `yaml:"name,omitempty"`
	Description           string                  `yaml:"description,omitempty"`
	Type                  ParameterType           `yaml:"type,omitempty"`
	Default               map[Network]interface{} `yaml:"default,omitempty"`
	MaxLength             int                     `yaml:"maxLength,omitempty"`
	Regex                 string                  `yaml:"regex,omitempty"`
	Advanced              bool                    `yaml:"advanced,omitempty"`
	AffectsContainers     []ContainerID           `yaml:"affectsContainers,omitempty"`
	EnvironmentVariables  []string                `yaml:"environmentVariables,omitempty"`
	CanBeBlank            bool                    `yaml:"canBeBlank,omitempty"`
	OverwriteOnUpgrade    bool                    `yaml:"overwriteOnUpgrade,omitempty"`
	Options               []ParameterOption       `yaml:"options,omitempty"`
	Value                 interface{}             `yaml:"-"`
	DescriptionsByNetwork map[Network]string      `yaml:"-"`
}
