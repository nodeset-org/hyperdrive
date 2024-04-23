package config

import (
	"github.com/nodeset-org/hyperdrive-daemon/shared/config/ids"
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
		"Hyperdrive and TX Fees",
		"Select this to configure the settings for Hyperdrive itself, including the defaults and limits on transaction fees.",
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
	masterConfig := configPage.home.md.Config
	layout := newStandardLayout()
	configPage.layout = layout
	layout.createForm(&masterConfig.Hyperdrive.Network, "Hyperdrive and TX Fee Settings")
	layout.setupEscapeReturnHomeHandler(configPage.home.md, configPage.home.homePage)

	// Set up the form items
	formItems := createParameterizedFormItems(masterConfig.Hyperdrive.GetParameters(), layout.descriptionBox)
	for _, formItem := range formItems {
		// Ignore the client mode item since it's presented in the EC / BN sections
		if formItem.parameter.GetCommon().ID == ids.ClientModeID {
			continue
		}

		// Handle the rest
		layout.form.AddFormItem(formItem.item)
		layout.parameters[formItem.item] = formItem
		if formItem.parameter.GetCommon().ID == ids.NetworkID {
			dropDown := formItem.item.(*DropDown)
			dropDown.SetSelectedFunc(func(text string, index int) {
				newNetwork := configPage.home.md.Config.Hyperdrive.Network.Options[index].Value
				configPage.home.md.Config.ChangeNetwork(newNetwork)
				configPage.home.refresh()
			})
		}
	}
	layout.refresh()

}

// Handle a bulk redraw request
func (configPage *HyperdriveConfigPage) handleLayoutChanged() {
	configPage.layout.refresh()
}
