package config

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

	// Set up the form items
	/*
		masterConfig := configPage.home.md.Instance
		formItems := createParameterizedFormItems(masterConfig.GetParameters(), layout.descriptionBox)
		for _, formItem := range formItems {
			layout.form.AddFormItem(formItem.item)
			layout.parameters[formItem.item] = formItem
		}
	*/
	layout.refresh()

}

// Handle a bulk redraw request
func (configPage *HyperdriveConfigPage) handleLayoutChanged() {
	configPage.layout.refresh()
}
