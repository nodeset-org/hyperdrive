package config

import (
	"github.com/gdamore/tcell/v2"
	"github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/nodeset-org/hyperdrive/shared/types"
)

// The page wrapper for the BN configs
type BeaconConfigPage struct {
	home               *settingsHome
	page               *page
	layout             *standardLayout
	masterConfig       *config.HyperdriveConfig
	clientModeDropdown *parameterizedFormItem
	localBnDropdown    *parameterizedFormItem
	externalBnDropdown *parameterizedFormItem
	localBnItems       []*parameterizedFormItem
	lighthouseItems    []*parameterizedFormItem
	lodestarItems      []*parameterizedFormItem
	nimbusItems        []*parameterizedFormItem
	prysmItems         []*parameterizedFormItem
	tekuItems          []*parameterizedFormItem
	externalBnItems    []*parameterizedFormItem
}

// Creates a new page for the Beacon Node settings
func NewBeaconConfigPage(home *settingsHome) *BeaconConfigPage {

	configPage := &BeaconConfigPage{
		home:         home,
		masterConfig: home.md.Config,
	}
	configPage.createContent()

	configPage.page = newPage(
		home.homePage,
		"settings-consensus",
		"Beacon Node",
		"Select this to choose your Beacon Node and configure its settings.",
		configPage.layout.grid,
	)

	return configPage

}

// Get the underlying page
func (configPage *BeaconConfigPage) getPage() *page {
	return configPage.page
}

// Creates the content for the Beacon Node settings page
func (configPage *BeaconConfigPage) createContent() {

	// Create the layout
	configPage.layout = newStandardLayout()
	configPage.layout.createForm(&configPage.masterConfig.Network, "Beacon Node Settings")

	// Return to the home page after pressing Escape
	configPage.layout.form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			// Close all dropdowns and break if one was open
			for _, param := range configPage.layout.parameters {
				dropDown, ok := param.item.(*DropDown)
				if ok && dropDown.open {
					dropDown.CloseList(configPage.home.md.app)
					return nil
				}
			}

			// Return to the home page
			configPage.home.md.setPage(configPage.home.homePage)
			return nil
		}
		return event
	})

	// Set up the form items
	configPage.clientModeDropdown = createParameterizedDropDown(&configPage.masterConfig.ClientMode, configPage.layout.descriptionBox)
	configPage.localBnDropdown = createParameterizedDropDown(&configPage.masterConfig.LocalBeaconConfig.BeaconNode, configPage.layout.descriptionBox)
	configPage.externalBnDropdown = createParameterizedDropDown(&configPage.masterConfig.ExternalBeaconConfig.BeaconNode, configPage.layout.descriptionBox)
	configPage.localBnItems = createParameterizedFormItems(configPage.masterConfig.LocalBeaconConfig.GetParameters(), configPage.layout.descriptionBox)
	configPage.lighthouseItems = createParameterizedFormItems(configPage.masterConfig.LocalBeaconConfig.Lighthouse.GetParameters(), configPage.layout.descriptionBox)
	configPage.lodestarItems = createParameterizedFormItems(configPage.masterConfig.LocalBeaconConfig.Lodestar.GetParameters(), configPage.layout.descriptionBox)
	configPage.nimbusItems = createParameterizedFormItems(configPage.masterConfig.LocalBeaconConfig.Nimbus.GetParameters(), configPage.layout.descriptionBox)
	configPage.prysmItems = createParameterizedFormItems(configPage.masterConfig.LocalBeaconConfig.Prysm.GetParameters(), configPage.layout.descriptionBox)
	configPage.tekuItems = createParameterizedFormItems(configPage.masterConfig.LocalBeaconConfig.Teku.GetParameters(), configPage.layout.descriptionBox)
	configPage.externalBnItems = createParameterizedFormItems(configPage.masterConfig.ExternalBeaconConfig.GetParameters(), configPage.layout.descriptionBox)

	// Map the parameters to the form items in the layout
	configPage.layout.mapParameterizedFormItems(configPage.clientModeDropdown, configPage.localBnDropdown, configPage.externalBnDropdown)
	configPage.layout.mapParameterizedFormItems(configPage.localBnItems...)
	configPage.layout.mapParameterizedFormItems(configPage.lighthouseItems...)
	configPage.layout.mapParameterizedFormItems(configPage.lodestarItems...)
	configPage.layout.mapParameterizedFormItems(configPage.nimbusItems...)
	configPage.layout.mapParameterizedFormItems(configPage.prysmItems...)
	configPage.layout.mapParameterizedFormItems(configPage.tekuItems...)
	configPage.layout.mapParameterizedFormItems(configPage.externalBnItems...)

	// Set up the setting callbacks
	configPage.clientModeDropdown.item.(*DropDown).SetSelectedFunc(func(text string, index int) {
		if configPage.masterConfig.ClientMode.Value == configPage.masterConfig.ClientMode.Options[index].Value {
			return
		}
		configPage.masterConfig.ClientMode.Value = configPage.masterConfig.ClientMode.Options[index].Value
		configPage.handleClientModeChanged()
	})
	configPage.localBnDropdown.item.(*DropDown).SetSelectedFunc(func(text string, index int) {
		if configPage.masterConfig.LocalBeaconConfig.BeaconNode.Value == configPage.masterConfig.LocalBeaconConfig.BeaconNode.Options[index].Value {
			return
		}
		configPage.masterConfig.LocalBeaconConfig.BeaconNode.Value = configPage.masterConfig.LocalBeaconConfig.BeaconNode.Options[index].Value
		configPage.handleLocalBnChanged()
	})
	configPage.externalBnDropdown.item.(*DropDown).SetSelectedFunc(func(text string, index int) {
		if configPage.masterConfig.ExternalBeaconConfig.BeaconNode.Value == configPage.masterConfig.ExternalBeaconConfig.BeaconNode.Options[index].Value {
			return
		}
		configPage.masterConfig.ExternalBeaconConfig.BeaconNode.Value = configPage.masterConfig.ExternalBeaconConfig.BeaconNode.Options[index].Value
		configPage.handleExternalBnChanged()
	})

	// Do the initial draw
	configPage.handleClientModeChanged()

}

// Handle all of the form changes when the client mode has changed
func (configPage *BeaconConfigPage) handleClientModeChanged() {
	configPage.layout.form.Clear(true)
	configPage.layout.form.AddFormItem(configPage.clientModeDropdown.item)

	selectedMode := configPage.masterConfig.ClientMode.Value
	switch selectedMode {
	case types.ClientMode_Local:
		// Local (Docker mode)
		configPage.handleLocalBnChanged()

	case types.ClientMode_External:
		// External (Hybrid mode)
		configPage.handleExternalBnChanged()
	}
}

// Handle all of the form changes when the BN has changed (local mode)
func (configPage *BeaconConfigPage) handleLocalBnChanged() {
	configPage.layout.form.Clear(true)
	configPage.layout.form.AddFormItem(configPage.clientModeDropdown.item)
	configPage.layout.form.AddFormItem(configPage.localBnDropdown.item)
	selectedBn := configPage.masterConfig.LocalBeaconConfig.BeaconNode.Value

	switch selectedBn {
	case types.BeaconNode_Lighthouse:
		configPage.layout.addFormItemsWithCommonParams(configPage.localBnItems, configPage.lighthouseItems, nil)
	case types.BeaconNode_Lodestar:
		configPage.layout.addFormItemsWithCommonParams(configPage.localBnItems, configPage.lodestarItems, nil)
	case types.BeaconNode_Nimbus:
		configPage.layout.addFormItemsWithCommonParams(configPage.localBnItems, configPage.nimbusItems, nil)
	case types.BeaconNode_Prysm:
		configPage.layout.addFormItemsWithCommonParams(configPage.localBnItems, configPage.prysmItems, nil)
	case types.BeaconNode_Teku:
		configPage.layout.addFormItemsWithCommonParams(configPage.localBnItems, configPage.tekuItems, nil)
	}

	configPage.layout.refresh()
}

// Handle all of the form changes when the BN has changed (external mode)
func (configPage *BeaconConfigPage) handleExternalBnChanged() {
	configPage.layout.form.Clear(true)
	configPage.layout.form.AddFormItem(configPage.clientModeDropdown.item)
	configPage.layout.form.AddFormItem(configPage.externalBnDropdown.item)

	// Split into Prysm and non-Prysm
	commonSettings := []*parameterizedFormItem{}
	prysmSettings := []*parameterizedFormItem{}
	for _, item := range configPage.externalBnItems {
		if item.parameter.GetCommon().ID == config.PrysmRpcUrlID {
			prysmSettings = append(prysmSettings, item)
		} else {
			commonSettings = append(commonSettings, item)
		}
	}

	// Show items based on the client selection
	configPage.layout.addFormItems(commonSettings)
	if configPage.masterConfig.ExternalBeaconConfig.BeaconNode.Value == types.BeaconNode_Prysm {
		configPage.layout.addFormItems(prysmSettings)
	}

	configPage.layout.refresh()
}

// Handle a bulk redraw request
func (configPage *BeaconConfigPage) handleLayoutChanged() {
	configPage.handleClientModeChanged()
}
