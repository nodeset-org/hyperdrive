package template

import (
	"fmt"
	"path/filepath"
)

const (
	TemplateSuffix    string = ".tmpl"
	ComposeFileSuffix string = ".yml"
)

type ComposePaths struct {
	RuntimePath  string
	TemplatePath string
	OverridePath string
}

type ComposeFile struct {
	name  string
	paths *ComposePaths
}

func (c *ComposePaths) File(name string) *ComposeFile {
	return &ComposeFile{
		name:  name,
		paths: c,
	}
}

// Given a ComposeFile returned by ComposePaths.File, find and parse the .tmpl
// from the TemplatePath, populate and save to the RuntimePath, and return a
// slice of compose definitions pertaining to the container (including the override).
func (c *ComposeFile) Write(data interface{}) ([]string, error) {
	composePath := filepath.Join(c.paths.RuntimePath, c.name+ComposeFileSuffix)
	tmpl := Template{
		Src: filepath.Join(c.paths.TemplatePath, c.name+TemplateSuffix),
		Dst: composePath,
	}
	err := tmpl.Write(data)
	if err != nil {
		return nil, fmt.Errorf("error writing %s compose definition: %w", c.name, err)
	}

	return []string{
		composePath,
		filepath.Join(c.paths.OverridePath, c.name+ComposeFileSuffix),
	}, nil
}
