package metadata

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/alessio/shellescape"
)

type MetadataTemplate struct {
	Content string
}

func (t MetadataTemplate) Write(data interface{}) (string, error) {
	// Map dynamic getters and parse the template
	tmpl, err := template.New("content").Funcs(template.FuncMap{
		"GetValue":      data.(*MetadataDataSource).GetValue,
		"GetValueArray": data.(*MetadataDataSource).GetValueArray,
		"UseDefault":    data.(*MetadataDataSource).UseDefault,
	}).Parse(t.Content)
	if err != nil {
		return "", fmt.Errorf("Error reading template file %s: %w", shellescape.Quote(t.Content), err)
	}

	var output strings.Builder
	err = tmpl.Execute(&output, data)
	if err != nil {
		return "", fmt.Errorf("error executing template: %w", err)
	}

	return output.String(), nil
}
