package config

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/config/ids"
	"github.com/nodeset-org/hyperdrive/modules"
	modconfig "github.com/nodeset-org/hyperdrive/modules/config"
)

// The page wrapper for the logging config
type LoggingConfigPage struct {
	home      *settingsHome
	page      *page
	layout    *standardLayout
	formItems []*parameterizedFormItem
}

// Creates a new page for the logging settings
func NewLoggingConfigPage(home *settingsHome) *LoggingConfigPage {
	configPage := &LoggingConfigPage{
		home: home,
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
	md := p.home.md
	p.layout = newStandardLayout(md, modules.HyperdriveFqmn)
	p.layout.createForm("Logging Settings")
	p.layout.setupEscapeReturnHomeHandler(md, p.home.homePage)

	// Create form items for each parameter
	newInstance, err := md.newInstance.GetSection(modconfig.Identifier(ids.LoggingSectionID))
	if err != nil {
		panic(fmt.Errorf("error getting logging section: %w", err))
	}
	loggingCfg := md.Config.Logging
	params := loggingCfg.GetParameters()
	for _, param := range params {
		id := param.GetID()
		paramSetting, err := newInstance.GetParameter(id)
		if err != nil {
			panic(fmt.Errorf("error getting logging parameter setting [%s]: %w", id, err))
		}

		// Create the form item for the parameter
		pfi := createParameterizedFormItem(paramSetting, p.layout, p.handleLayoutChanged)
		p.layout.registerFormItems(pfi)
		p.formItems = append(p.formItems, pfi)
	}
}

// Handle a bulk redraw request
func (p *LoggingConfigPage) handleLayoutChanged() {
	p.layout.redrawForm(p.formItems, nil, nil)
}
