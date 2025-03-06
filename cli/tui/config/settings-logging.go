package config

import (
	"fmt"
	"strconv"
	"strings"
	"text/template"

	"github.com/nodeset-org/hyperdrive/config"
	"github.com/nodeset-org/hyperdrive/config/ids"
	"github.com/nodeset-org/hyperdrive/modules"
	modconfig "github.com/nodeset-org/hyperdrive/modules/config"
)

// The page wrapper for the logging config
type LoggingConfigPage struct {
	home         *settingsHome
	page         *page
	layout       *standardLayout
	masterConfig *config.HyperdriveConfig
	params       []*parameterizedFormItem
	redrawing    bool
}

// Creates a new page for the logging settings
func NewLoggingConfigPage(home *settingsHome) *LoggingConfigPage {
	configPage := &LoggingConfigPage{
		home:         home,
		masterConfig: home.md.Config,
	}
	configPage.createContent()

	configPage.page = newPage(
		home.homePage,
		"settings-logging",
		"Logging",
		"Configure Hyperdrive's daemon and module logs.",
		configPage.layout.grid,
	)

	// Do the initial draw
	configPage.handleLayoutChanged()
	return configPage
}

// Get the underlying page
func (p *LoggingConfigPage) getPage() *page {
	return p.page
}

// Creates the content for the logging settings page
func (p *LoggingConfigPage) createContent() {
	// Create the layout
	p.layout = newStandardLayout(p.home.md)
	p.layout.createForm("Logging Settings")
	p.layout.setupEscapeReturnHomeHandler(p.home.md, p.home.homePage)

	// Create form items for each parameter
	newInstance, err := p.home.md.newInstance.GetSection(modconfig.Identifier(ids.LoggingSectionID))
	if err != nil {
		panic(fmt.Errorf("error getting logging section: %w", err))
	}
	loggingCfg := p.home.md.Config.Logging
	params := loggingCfg.GetParameters()
	settings := []modconfig.IParameterSetting{}
	for _, param := range params {
		id := param.GetID()
		setting, err := newInstance.GetParameter(id)
		if err != nil {
			panic(fmt.Errorf("error getting logging parameter setting [%s]: %w", id, err))
		}
		settings = append(settings, setting)
	}

	// Set up the form items
	p.params = createParameterizedFormItems(settings, p.layout.descriptionBox, p.handleLayoutChanged)
	p.layout.mapParameterizedFormItems(p.params...)
}

// Handle a bulk redraw request
func (p *LoggingConfigPage) handleLayoutChanged() {
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
