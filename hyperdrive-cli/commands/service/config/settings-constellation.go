package config

import (
	csids "github.com/nodeset-org/hyperdrive-constellation/shared/config/ids"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/rivo/tview"
	"github.com/rocket-pool/node-manager-core/config"
)

// The page wrapper for the Constellation config
type ConstellationConfigPage struct {
	modulesPage            *ModulesPage
	page                   *page
	layout                 *standardLayout
	masterConfig           *client.GlobalConfig
	enableConstellationBox *parameterizedFormItem

	constellationItems []*parameterizedFormItem
	vcCommonItems      []*parameterizedFormItem
	lighthouseItems    []*parameterizedFormItem
	lodestarItems      []*parameterizedFormItem
	nimbusItems        []*parameterizedFormItem
	prysmItems         []*parameterizedFormItem
	tekuItems          []*parameterizedFormItem
}

// Creates a new page for the Constellation settings
func NewConstellationConfigPage(modulesPage *ModulesPage) *ConstellationConfigPage {

	configPage := &ConstellationConfigPage{
		modulesPage:  modulesPage,
		masterConfig: modulesPage.home.md.Config,
	}
	configPage.createContent()

	configPage.page = newPage(
		modulesPage.page,
		"settings-constellation",
		"Constellation",
		"Select this to manage the Constellation module and the Validator Client it uses.",
		configPage.layout.grid,
	)

	return configPage

}

// Get the underlying page
func (configPage *ConstellationConfigPage) getPage() *page {
	return configPage.page
}

// Creates the content for the Constellation settings page
func (configPage *ConstellationConfigPage) createContent() {

	// Create the layout
	configPage.layout = newStandardLayout()
	configPage.layout.createForm(&configPage.masterConfig.Hyperdrive.Network, "Constellation Settings")
	configPage.layout.setupEscapeReturnHomeHandler(configPage.modulesPage.home.md, configPage.modulesPage.page)

	// Set up the form items
	configPage.enableConstellationBox = createParameterizedCheckbox(&configPage.masterConfig.Constellation.Enabled)
	configPage.constellationItems = createParameterizedFormItems(configPage.masterConfig.Constellation.GetParameters(), configPage.layout.descriptionBox)
	configPage.vcCommonItems = createParameterizedFormItems(configPage.masterConfig.Constellation.VcCommon.GetParameters(), configPage.layout.descriptionBox)
	configPage.lighthouseItems = createParameterizedFormItems(configPage.masterConfig.Constellation.Lighthouse.GetParameters(), configPage.layout.descriptionBox)
	configPage.lodestarItems = createParameterizedFormItems(configPage.masterConfig.Constellation.Lodestar.GetParameters(), configPage.layout.descriptionBox)
	configPage.nimbusItems = createParameterizedFormItems(configPage.masterConfig.Constellation.Nimbus.GetParameters(), configPage.layout.descriptionBox)
	configPage.prysmItems = createParameterizedFormItems(configPage.masterConfig.Constellation.Prysm.GetParameters(), configPage.layout.descriptionBox)
	configPage.tekuItems = createParameterizedFormItems(configPage.masterConfig.Constellation.Teku.GetParameters(), configPage.layout.descriptionBox)

	// Map the parameters to the form items in the layout
	configPage.layout.mapParameterizedFormItems(configPage.enableConstellationBox)
	configPage.layout.mapParameterizedFormItems(configPage.constellationItems...)
	configPage.layout.mapParameterizedFormItems(configPage.vcCommonItems...)
	configPage.layout.mapParameterizedFormItems(configPage.lighthouseItems...)
	configPage.layout.mapParameterizedFormItems(configPage.lodestarItems...)
	configPage.layout.mapParameterizedFormItems(configPage.nimbusItems...)
	configPage.layout.mapParameterizedFormItems(configPage.prysmItems...)
	configPage.layout.mapParameterizedFormItems(configPage.tekuItems...)

	// Set up the setting callbacks
	configPage.enableConstellationBox.item.(*tview.Checkbox).SetChangedFunc(func(checked bool) {
		if configPage.masterConfig.Constellation.Enabled.Value == checked {
			return
		}
		configPage.masterConfig.Constellation.Enabled.Value = checked
		configPage.handleLayoutChanged()
	})

	// Do the initial draw
	configPage.handleLayoutChanged()
}

// Handle all of the form changes when the Enable Metrics box has changed
func (configPage *ConstellationConfigPage) handleLayoutChanged() {
	configPage.layout.form.Clear(true)
	configPage.layout.form.AddFormItem(configPage.enableConstellationBox.item)

	if configPage.masterConfig.Constellation.Enabled.Value {
		// Remove the Constellation enable param since it's already there
		csItems := []*parameterizedFormItem{}
		for _, item := range configPage.constellationItems {
			if item.parameter.GetCommon().ID == csids.ConstellationEnableID {
				continue
			}
			csItems = append(csItems, item)
		}
		configPage.layout.addFormItems(csItems)

		// Display the relevant VC items
		configPage.layout.addFormItems(configPage.vcCommonItems)

		bn := configPage.masterConfig.Hyperdrive.GetSelectedBeaconNode()
		switch bn {
		case config.BeaconNode_Lighthouse:
			configPage.layout.addFormItems(configPage.lighthouseItems)
		case config.BeaconNode_Lodestar:
			configPage.layout.addFormItems(configPage.lodestarItems)
		case config.BeaconNode_Nimbus:
			configPage.layout.addFormItems(configPage.nimbusItems)
		case config.BeaconNode_Prysm:
			configPage.layout.addFormItems(configPage.prysmItems)
		case config.BeaconNode_Teku:
			configPage.layout.addFormItems(configPage.tekuItems)
		}
	}

	configPage.layout.refresh()
}
