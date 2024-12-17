package module

import (
	"fmt"
	"os"
	"text/template"

	"github.com/alessio/shellescape"
)

type Template struct {
	Src string
	Dst string
}

func (t Template) Write(data interface{}) error {
	// Open the output file, creating it if it doesn't exist
	runtimeFile, err := os.OpenFile(t.Dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0664)
	if err != nil {
		return fmt.Errorf("Could not open templated file %s for writing: %w", shellescape.Quote(t.Dst), err)
	}
	defer runtimeFile.Close()

	// Map dynamic getters and parse the template
	// tmpl, err := template.New("template").Funcs(template.FuncMap{
	// 	"GetValue":      data.(*TemplateDataSource).GetValue,
	// 	"GetValueArray": data.(*TemplateDataSource).GetValueArray,
	// 	"UseDefault":    data.(*TemplateDataSource).UseDefault,
	// }).ParseFiles(t.Src)
	tmpl, err := template.ParseFiles(t.Src)
	if err != nil {
		return fmt.Errorf("Error reading template file %s: %w", shellescape.Quote(t.Src), err)
	}

	// Replace template variables and write the result
	err = tmpl.Execute(runtimeFile, data)
	if err != nil {
		return fmt.Errorf("Error writing and substituting template: %w", err)
	}

	// If the file was newly created, 0664 may have been altered by umask, so chmod back to 0664.
	err = os.Chmod(t.Dst, 0664)
	if err != nil {
		return fmt.Errorf("Could not set templated file (%s) permissions: %w", shellescape.Quote(t.Dst), err)
	}

	return nil
}
