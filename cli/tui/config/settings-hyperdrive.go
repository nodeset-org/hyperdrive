package config

import (
	"fmt"

	modconfig "github.com/nodeset-org/hyperdrive/modules/config"
)

// The page wrapper for the Hyperdrive config
type HyperdriveConfigPage struct {
	home   *settingsHome
	page   *page
	layout *standardLayout
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

	return configPage
}

// Get the underlying page
func (configPage *HyperdriveConfigPage) getPage() *page {
	return configPage.page
}

// Creates the content for the Hyperdrive settings page
func (configPage *HyperdriveConfigPage) createContent() {
	// Create the layout
	layout := newStandardLayout()
	configPage.layout = layout
	layout.createForm("Hyperdrive Settings")
	layout.setupEscapeReturnHomeHandler(configPage.home.md, configPage.home.homePage)

	params := configPage.home.md.Config.GetParameters()
	newInstance := configPage.home.md.newInstance
	settings := []modconfig.IParameterSetting{}
	for _, param := range params {
		id := param.GetID()
		setting, err := newInstance.GetParameter(id)
		if err != nil {
			panic(fmt.Errorf("error getting base parameter setting [%s]: %w", id, err))
		}
		settings = append(settings, setting)
	}

	// Set up the form items
	formItems := createParameterizedFormItems(settings, layout.descriptionBox)
	for _, formItem := range formItems {
		layout.form.AddFormItem(formItem.item)
		layout.parameters[formItem.item] = formItem
	}
	configPage.layout.mapParameterizedFormItems(formItems...)
	layout.refresh()
}

// Handle a bulk redraw request
func (configPage *HyperdriveConfigPage) handleLayoutChanged() {
	configPage.layout.refresh()
}
