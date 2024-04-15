package config

import (
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	swids "github.com/nodeset-org/hyperdrive/modules/stakewise/shared/config/ids"
	"github.com/rivo/tview"
	"github.com/rocket-pool/node-manager-core/config"
)

// The page wrapper for the Stakewise config
type StakewiseConfigPage struct {
	modulesPage        *ModulesPage
	page               *page
	layout             *standardLayout
	masterConfig       *client.GlobalConfig
	enableStakewiseBox *parameterizedFormItem

	stakewiseItems  []*parameterizedFormItem
	vcCommonItems   []*parameterizedFormItem
	lighthouseItems []*parameterizedFormItem
	lodestarItems   []*parameterizedFormItem
	nimbusItems     []*parameterizedFormItem
	prysmItems      []*parameterizedFormItem
	tekuItems       []*parameterizedFormItem
}

// Creates a new page for the Stakewise settings
func NewStakewiseConfigPage(modulesPage *ModulesPage) *StakewiseConfigPage {

	configPage := &StakewiseConfigPage{
		modulesPage:  modulesPage,
		masterConfig: modulesPage.home.md.Config,
	}
	configPage.createContent()

	configPage.page = newPage(
		modulesPage.page,
		"settings-stakewise",
		"Stakewise",
		"Select this to manage the Stakewise module and the Validator Client it uses.",
		configPage.layout.grid,
	)

	return configPage

}

// Get the underlying page
func (configPage *StakewiseConfigPage) getPage() *page {
	return configPage.page
}

// Creates the content for the Stakewise settings page
func (configPage *StakewiseConfigPage) createContent() {

	// Create the layout
	configPage.layout = newStandardLayout()
	configPage.layout.createForm(&configPage.masterConfig.Hyperdrive.Network, "Stakewise Settings")
	configPage.layout.setupEscapeReturnHomeHandler(configPage.modulesPage.home.md, configPage.modulesPage.page)

	// Set up the form items
	configPage.enableStakewiseBox = createParameterizedCheckbox(&configPage.masterConfig.Stakewise.Enabled)
	configPage.stakewiseItems = createParameterizedFormItems(configPage.masterConfig.Stakewise.GetParameters(), configPage.layout.descriptionBox)
	configPage.vcCommonItems = createParameterizedFormItems(configPage.masterConfig.Stakewise.VcCommon.GetParameters(), configPage.layout.descriptionBox)
	configPage.lighthouseItems = createParameterizedFormItems(configPage.masterConfig.Stakewise.Lighthouse.GetParameters(), configPage.layout.descriptionBox)
	configPage.lodestarItems = createParameterizedFormItems(configPage.masterConfig.Stakewise.Lodestar.GetParameters(), configPage.layout.descriptionBox)
	configPage.nimbusItems = createParameterizedFormItems(configPage.masterConfig.Stakewise.Nimbus.GetParameters(), configPage.layout.descriptionBox)
	configPage.prysmItems = createParameterizedFormItems(configPage.masterConfig.Stakewise.Prysm.GetParameters(), configPage.layout.descriptionBox)
	configPage.tekuItems = createParameterizedFormItems(configPage.masterConfig.Stakewise.Teku.GetParameters(), configPage.layout.descriptionBox)

	// Map the parameters to the form items in the layout
	configPage.layout.mapParameterizedFormItems(configPage.enableStakewiseBox)
	configPage.layout.mapParameterizedFormItems(configPage.stakewiseItems...)
	configPage.layout.mapParameterizedFormItems(configPage.vcCommonItems...)
	configPage.layout.mapParameterizedFormItems(configPage.lighthouseItems...)
	configPage.layout.mapParameterizedFormItems(configPage.lodestarItems...)
	configPage.layout.mapParameterizedFormItems(configPage.nimbusItems...)
	configPage.layout.mapParameterizedFormItems(configPage.prysmItems...)
	configPage.layout.mapParameterizedFormItems(configPage.tekuItems...)

	// Set up the setting callbacks
	configPage.enableStakewiseBox.item.(*tview.Checkbox).SetChangedFunc(func(checked bool) {
		if configPage.masterConfig.Stakewise.Enabled.Value == checked {
			return
		}
		configPage.masterConfig.Stakewise.Enabled.Value = checked
		configPage.handleLayoutChanged()
	})

	// Do the initial draw
	configPage.handleLayoutChanged()
}

// Handle all of the form changes when the Enable Metrics box has changed
func (configPage *StakewiseConfigPage) handleLayoutChanged() {
	configPage.layout.form.Clear(true)
	configPage.layout.form.AddFormItem(configPage.enableStakewiseBox.item)

	if configPage.masterConfig.Stakewise.Enabled.Value {
		// Remove the Stakewise enable param since it's already there
		stakewiseItems := []*parameterizedFormItem{}
		for _, item := range configPage.stakewiseItems {
			if item.parameter.GetCommon().ID == swids.StakewiseEnableID {
				continue
			}
			stakewiseItems = append(stakewiseItems, item)
		}
		configPage.layout.addFormItems(stakewiseItems)

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
