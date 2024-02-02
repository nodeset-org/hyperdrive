package config

import (
	"github.com/gdamore/tcell/v2"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	swconfig "github.com/nodeset-org/hyperdrive/modules/stakewise/shared/config"
	"github.com/nodeset-org/hyperdrive/shared/types"
	"github.com/rivo/tview"
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

	// Return to the home page after pressing Escape
	configPage.layout.form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Return to the modules page
		if event.Key() == tcell.KeyEsc {
			// Close all dropdowns and break if one was open
			for _, param := range configPage.layout.parameters {
				dropDown, ok := param.item.(*DropDown)
				if ok && dropDown.open {
					dropDown.CloseList(configPage.modulesPage.home.md.app)
					return nil
				}
			}

			configPage.modulesPage.home.md.setPage(configPage.modulesPage.home.homePage)
			return nil
		}
		return event
	})

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

	if configPage.masterConfig.Stakewise.Enabled.Value == true {
		// Remove the Stakewise enable param since it's already there
		stakewiseItems := []*parameterizedFormItem{}
		for _, item := range configPage.stakewiseItems {
			if item.parameter.GetCommon().ID == swconfig.StakewiseEnableID {
				continue
			}
			stakewiseItems = append(stakewiseItems, item)
		}
		configPage.layout.addFormItems(stakewiseItems)

		// Display the relevant VC items
		configPage.layout.addFormItems(configPage.vcCommonItems)

		bn := configPage.masterConfig.Hyperdrive.GetSelectedBeaconNode()
		switch bn {
		case types.BeaconNode_Lighthouse:
			configPage.layout.addFormItems(configPage.lighthouseItems)
		case types.BeaconNode_Lodestar:
			configPage.layout.addFormItems(configPage.lodestarItems)
		case types.BeaconNode_Nimbus:
			configPage.layout.addFormItems(configPage.nimbusItems)
		case types.BeaconNode_Prysm:
			configPage.layout.addFormItems(configPage.prysmItems)
		case types.BeaconNode_Teku:
			configPage.layout.addFormItems(configPage.tekuItems)
		}
	}

	configPage.layout.refresh()
}
