package config

import (
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/rocket-pool/node-manager-core/config"
	"github.com/rocket-pool/node-manager-core/config/ids"
)

// The page wrapper for the BN configs
type BeaconConfigPage struct {
	home               *settingsHome
	page               *page
	layout             *standardLayout
	masterConfig       *client.GlobalConfig
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
		"settings-beacon",
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
	configPage.layout.createForm(&configPage.masterConfig.Hyperdrive.Network, "Beacon Node Settings")
	configPage.layout.setupEscapeReturnHomeHandler(configPage.home.md, configPage.home.homePage)

	// Set up the form items
	configPage.clientModeDropdown = createParameterizedDropDown(&configPage.masterConfig.Hyperdrive.ClientMode, configPage.layout.descriptionBox)
	configPage.localBnDropdown = createParameterizedDropDown(&configPage.masterConfig.Hyperdrive.LocalBeaconClient.BeaconNode, configPage.layout.descriptionBox)
	configPage.externalBnDropdown = createParameterizedDropDown(&configPage.masterConfig.Hyperdrive.ExternalBeaconClient.BeaconNode, configPage.layout.descriptionBox)
	configPage.localBnItems = createParameterizedFormItems(configPage.masterConfig.Hyperdrive.LocalBeaconClient.GetParameters(), configPage.layout.descriptionBox)
	configPage.lighthouseItems = createParameterizedFormItems(configPage.masterConfig.Hyperdrive.LocalBeaconClient.Lighthouse.GetParameters(), configPage.layout.descriptionBox)
	configPage.lodestarItems = createParameterizedFormItems(configPage.masterConfig.Hyperdrive.LocalBeaconClient.Lodestar.GetParameters(), configPage.layout.descriptionBox)
	configPage.nimbusItems = createParameterizedFormItems(configPage.masterConfig.Hyperdrive.LocalBeaconClient.Nimbus.GetParameters(), configPage.layout.descriptionBox)
	configPage.prysmItems = createParameterizedFormItems(configPage.masterConfig.Hyperdrive.LocalBeaconClient.Prysm.GetParameters(), configPage.layout.descriptionBox)
	configPage.tekuItems = createParameterizedFormItems(configPage.masterConfig.Hyperdrive.LocalBeaconClient.Teku.GetParameters(), configPage.layout.descriptionBox)
	configPage.externalBnItems = createParameterizedFormItems(configPage.masterConfig.Hyperdrive.ExternalBeaconClient.GetParameters(), configPage.layout.descriptionBox)

	// Take the client selections out since they're done explicitly
	localBnItems := []*parameterizedFormItem{}
	for _, item := range configPage.localBnItems {
		if item.parameter.GetCommon().ID == ids.BnID {
			continue
		}
		localBnItems = append(localBnItems, item)
	}
	configPage.localBnItems = localBnItems

	externalBnItems := []*parameterizedFormItem{}
	for _, item := range configPage.externalBnItems {
		if item.parameter.GetCommon().ID == ids.BnID {
			continue
		}
		externalBnItems = append(externalBnItems, item)
	}
	configPage.externalBnItems = externalBnItems

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
		if configPage.masterConfig.Hyperdrive.ClientMode.Value == configPage.masterConfig.Hyperdrive.ClientMode.Options[index].Value {
			return
		}
		configPage.masterConfig.Hyperdrive.ClientMode.Value = configPage.masterConfig.Hyperdrive.ClientMode.Options[index].Value
		configPage.handleClientModeChanged()
	})
	configPage.localBnDropdown.item.(*DropDown).SetSelectedFunc(func(text string, index int) {
		if configPage.masterConfig.Hyperdrive.LocalBeaconClient.BeaconNode.Value == configPage.masterConfig.Hyperdrive.LocalBeaconClient.BeaconNode.Options[index].Value {
			return
		}
		configPage.masterConfig.Hyperdrive.LocalBeaconClient.BeaconNode.Value = configPage.masterConfig.Hyperdrive.LocalBeaconClient.BeaconNode.Options[index].Value
		configPage.handleLocalBnChanged()
	})
	configPage.externalBnDropdown.item.(*DropDown).SetSelectedFunc(func(text string, index int) {
		if configPage.masterConfig.Hyperdrive.ExternalBeaconClient.BeaconNode.Value == configPage.masterConfig.Hyperdrive.ExternalBeaconClient.BeaconNode.Options[index].Value {
			return
		}
		configPage.masterConfig.Hyperdrive.ExternalBeaconClient.BeaconNode.Value = configPage.masterConfig.Hyperdrive.ExternalBeaconClient.BeaconNode.Options[index].Value
		configPage.handleExternalBnChanged()
	})

	// Do the initial draw
	configPage.handleClientModeChanged()
}

// Handle all of the form changes when the client mode has changed
func (configPage *BeaconConfigPage) handleClientModeChanged() {
	configPage.layout.form.Clear(true)
	configPage.layout.form.AddFormItem(configPage.clientModeDropdown.item)

	selectedMode := configPage.masterConfig.Hyperdrive.ClientMode.Value
	switch selectedMode {
	case config.ClientMode_Local:
		// Local (Docker mode)
		configPage.handleLocalBnChanged()

	case config.ClientMode_External:
		// External (Hybrid mode)
		configPage.handleExternalBnChanged()
	}
}

// Handle all of the form changes when the BN has changed (local mode)
func (configPage *BeaconConfigPage) handleLocalBnChanged() {
	configPage.layout.form.Clear(true)
	configPage.layout.form.AddFormItem(configPage.clientModeDropdown.item)
	configPage.layout.form.AddFormItem(configPage.localBnDropdown.item)
	selectedBn := configPage.masterConfig.Hyperdrive.LocalBeaconClient.BeaconNode.Value

	switch selectedBn {
	case config.BeaconNode_Lighthouse:
		configPage.layout.addFormItemsWithCommonParams(configPage.localBnItems, configPage.lighthouseItems, nil)
	case config.BeaconNode_Lodestar:
		configPage.layout.addFormItemsWithCommonParams(configPage.localBnItems, configPage.lodestarItems, nil)
	case config.BeaconNode_Nimbus:
		configPage.layout.addFormItemsWithCommonParams(configPage.localBnItems, configPage.nimbusItems, nil)
	case config.BeaconNode_Prysm:
		configPage.layout.addFormItemsWithCommonParams(configPage.localBnItems, configPage.prysmItems, nil)
	case config.BeaconNode_Teku:
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
		if item.parameter.GetCommon().ID == ids.PrysmRpcUrlID {
			prysmSettings = append(prysmSettings, item)
		} else {
			commonSettings = append(commonSettings, item)
		}
	}

	// Show items based on the client selection
	configPage.layout.addFormItems(commonSettings)
	if configPage.masterConfig.Hyperdrive.ExternalBeaconClient.BeaconNode.Value == config.BeaconNode_Prysm {
		configPage.layout.addFormItems(prysmSettings)
	}

	configPage.layout.refresh()
}

// Handle a bulk redraw request
func (configPage *BeaconConfigPage) handleLayoutChanged() {
	configPage.handleClientModeChanged()
}
