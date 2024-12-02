package modules

import (
	"fmt"
	"regexp"

	"github.com/blang/semver/v4"
)

const (
	// The regex pattern for a dependency
	DependencyPattern string = `^(?P<author>[^\s]+)\/(?P<name>[^\s]+)( (?P<version_op><|<=|=|>=|>) (?P<version>([^\s]+)))?$`
)

var (
	// The regex for a dependency
	dependencyRegex = regexp.MustCompile(DependencyPattern)
)

// An operator signifying how to interpret the version of a dependency required by a module
type VersionOperator string

const (
	// Unknown / unset operator
	VersionOperator_Unknown VersionOperator = ""

	// Less than
	VersionOperator_LessThan VersionOperator = "<"

	// Less than or equal to
	VersionOperator_LessThanOrEqual VersionOperator = "<="

	// Equal to
	VersionOperator_Equal VersionOperator = "="

	// Greater than or equal to
	VersionOperator_GreaterThanOrEqual VersionOperator = ">="

	// Greater than
	VersionOperator_GreaterThan VersionOperator = ">"
)

// A dependency string - including the author, module name, and optional version with an operator
type Dependency struct {
	// The author of the dependency
	Author DescriptorString `json:"author" yaml:"author"`

	// The name of the dependency
	Name DescriptorString `json:"name" yaml:"name"`

	// The operator to use when determining how to interpret the dependency's version
	VersionOp VersionOperator `json:"versionOp,omitempty" yaml:"versionOp,omitempty"`

	// The version of the dependency
	Version semver.Version `json:"version,omitempty" yaml:"version,omitempty"`
}

// String representation of a dependency
func (d Dependency) String() string {
	if d.VersionOp == "" {
		return fmt.Sprintf("%s/%s", d.Author, d.Name)
	}
	return fmt.Sprintf("%s/%s %s %s", d.Author, d.Name, d.VersionOp, d.Version)
}

// Marshal the dependency to text
func (d Dependency) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

// Unmarshal the dependency from text
func (d *Dependency) UnmarshalText(text []byte) error {
	// Check the text against the dependency regex
	submatches := dependencyRegex.FindStringSubmatch(string(text))
	if len(submatches) == 0 {
		return fmt.Errorf("invalid dependency string: %s", text)
	}

	// Parse the text into the dependency struct
	var dependency Dependency
	versionSet := false
	for i, name := range dependencyRegex.SubexpNames() {
		if i != 0 && name != "" && i < len(submatches) {
			captureGroup := submatches[i]
			switch name {
			case "author":
				err := dependency.Author.UnmarshalText([]byte(captureGroup))
				if err != nil {
					return fmt.Errorf("error parsing dependency author [%s] for [%s]: %w", captureGroup, string(text), err)
				}
			case "name":
				err := dependency.Name.UnmarshalText([]byte(captureGroup))
				if err != nil {
					return fmt.Errorf("error parsing dependency name [%s] for [%s]: %w", captureGroup, string(text), err)
				}
			case "version_op":
				dependency.VersionOp = VersionOperator(captureGroup)
			case "version":
				version, err := semver.Parse(captureGroup)
				if err != nil {
					return fmt.Errorf("error parsing dependency version [%s] for [%s]: %w", captureGroup, string(text), err)
				}
				dependency.Version = version
				versionSet = true
			}
		}
	}

	// Make sure the author and name are set
	if dependency.Author == "" {
		return fmt.Errorf("no author was provided for dependency [%s]", string(text))
	}
	if dependency.Name == "" {
		return fmt.Errorf("no name was provided for dependency [%s]", string(text))
	}

	// If a version was set, ensure an operator was also set
	if versionSet && dependency.VersionOp == VersionOperator_Unknown {
		return fmt.Errorf("a version was set for dependency [%s], but no comparison operator was provided", string(text))
	}

	// Make sure there's a version if there's an operator
	if !versionSet && dependency.VersionOp != VersionOperator_Unknown {
		return fmt.Errorf("a comparison operator was set for dependency [%s], but no version was provided", string(text))
	}

	*d = dependency
	return nil
}
