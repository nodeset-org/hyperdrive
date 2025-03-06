package config

import (
	"fmt"
	"strconv"
	"strings"
	"text/template"

	"github.com/nodeset-org/hyperdrive/modules"
	modconfig "github.com/nodeset-org/hyperdrive/modules/config"
)

// The page wrapper for the Hyperdrive config
type HyperdriveConfigPage struct {
	home   *settingsHome
	page   *page
	layout *standardLayout

	params    []*parameterizedFormItem
	redrawing bool
}

// Creates a new page for the Hyperdrive settings
func NewHyperdriveConfigPage(home *settingsHome) *HyperdriveConfigPage {
	configPage := &HyperdriveConfigPage{
		home: home,
	}

	configPage.createContent()
	configPage.page = newPage(
		home.homePage,
		"settings-hyperdrive",
		"Hyperdrive",
		"Select this to configure the settings for Hyperdrive itself.",
		configPage.layout.grid,
	)

	// Do the initial draw
	configPage.handleLayoutChanged()
	return configPage
}

// Get the underlying page
func (p *HyperdriveConfigPage) getPage() *page {
	return p.page
}

// Creates the content for the Hyperdrive settings page
func (p *HyperdriveConfigPage) createContent() {
	// Create the layout
	layout := newStandardLayout(p.home.md)
	p.layout = layout
	layout.createForm("Hyperdrive Settings")
	layout.setupEscapeReturnHomeHandler(p.home.md, p.home.homePage)

	params := p.home.md.Config.GetParameters()
	newInstance := p.home.md.newInstance
	settings := []modconfig.IParameterSetting{}
	for _, param := range params {
		id := param.GetID()
		setting, err := newInstance.GetParameter(id)
		if err != nil {
			panic(fmt.Errorf("error getting base parameter setting [%s]: %w", id, err))
		}
		settings = append(settings, setting)
	}
	p.params = createParameterizedFormItems(settings, layout.descriptionBox, p.handleLayoutChanged)
	p.layout.mapParameterizedFormItems(p.params...)
}

// Handle a bulk redraw request
func (p *HyperdriveConfigPage) handleLayoutChanged() {
	if p.redrawing {
		return
	}
	p.redrawing = true
	defer func() {
		p.redrawing = false
	}()

	p.layout.form.Clear(true)

	params := []*parameterizedFormItem{}
	md := p.home.md
	for _, param := range p.params {
		metadata := param.parameter.GetMetadata()
		hidden := metadata.GetHidden()

		// Handle parameters that don't have a hidden template
		if hidden.Template == "" {
			if !hidden.Default {
				params = append(params, param)
			}
			continue
		}

		// Generate a template source for the parameter
		templateSource := parameterTemplateSource{
			configurationTemplateSource: configurationTemplateSource{
				fqmn:              modules.HyperdriveFqmn,
				hdSettings:        md.newInstance,
				moduleSettingsMap: md.moduleSettingsMap,
			},
			parameter: param.parameter.GetMetadata(),
		}

		// Update the hidden status
		template, err := template.New(string(metadata.GetID())).Parse(hidden.Template)
		if err != nil {
			fqmn := modules.HyperdriveFqmn
			panic(fmt.Errorf("error parsing hidden template for parameter [%s:%s]: %w", fqmn, metadata.GetID(), err))
		}
		result := &strings.Builder{}
		err = template.Execute(result, templateSource)
		if err != nil {
			fqmn := modules.HyperdriveFqmn
			panic(fmt.Errorf("error executing hidden template for parameter [%s:%s]: %w", fqmn, metadata.GetID(), err))
		}

		hiddenResult, err := strconv.ParseBool(result.String())
		if err != nil {
			fqmn := modules.HyperdriveFqmn
			panic(fmt.Errorf("error parsing hidden template result for parameter [%s:%s]: %w", fqmn, metadata.GetID(), err))
		}
		if !hiddenResult {
			params = append(params, param)
		}
	}
	p.layout.addFormItems(params)

	p.layout.refresh()
}
