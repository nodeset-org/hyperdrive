package config

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/config"
	"github.com/nodeset-org/hyperdrive/config/ids"
	modconfig "github.com/nodeset-org/hyperdrive/modules/config"
)

// The page wrapper for the logging config
type LoggingConfigPage struct {
	home         *settingsHome
	page         *page
	layout       *standardLayout
	masterConfig *config.HyperdriveConfig
	loggingItems []*parameterizedFormItem
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

	return configPage
}

// Get the underlying page
func (configPage *LoggingConfigPage) getPage() *page {
	return configPage.page
}

// Creates the content for the logging settings page
func (configPage *LoggingConfigPage) createContent() {
	// Create the layout
	configPage.layout = newStandardLayout()
	configPage.layout.createForm("Logging Settings")
	configPage.layout.setupEscapeReturnHomeHandler(configPage.home.md, configPage.home.homePage)

	// Create form items for each parameter
	newInstance, err := configPage.home.md.newInstance.GetSection(modconfig.Identifier(ids.LoggingSectionID))
	if err != nil {
		panic(fmt.Errorf("error getting logging section: %w", err))
	}
	loggingCfg := configPage.home.md.Config.Logging
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
	configPage.loggingItems = createParameterizedFormItems(settings, configPage.layout.descriptionBox)
	for _, formItem := range configPage.loggingItems {
		configPage.layout.form.AddFormItem(formItem.item)
		configPage.layout.parameters[formItem.item] = formItem
	}
	configPage.layout.refresh()

	// Map the parameters to the form items in the layout
	configPage.layout.mapParameterizedFormItems(configPage.loggingItems...)

	// Do the initial draw
	configPage.handleLayoutChanged()
}

// Handle a bulk redraw request
func (configPage *LoggingConfigPage) handleLayoutChanged() {
	configPage.layout.form.Clear(true)
	configPage.layout.addFormItems(configPage.loggingItems)
	configPage.layout.refresh()
}
