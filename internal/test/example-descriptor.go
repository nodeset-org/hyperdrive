package internal_test

import (
	"github.com/blang/semver/v4"
	"github.com/nodeset-org/hyperdrive/modules"
)

var (
	ExampleDescriptor modules.ModuleDescriptor = modules.ModuleDescriptor{
		Name:         "example-module",
		Shortcut:     "em",
		Description:  "Simple example of a Hyperdrive module",
		Version:      semver.MustParse("0.2.0"),
		Author:       "NodeSet",
		URL:          "https://github.com/nodeset-org/hyperdrive-example",
		Email:        "info@nodeset.io",
		Dependencies: []modules.Dependency{},
	}
)
