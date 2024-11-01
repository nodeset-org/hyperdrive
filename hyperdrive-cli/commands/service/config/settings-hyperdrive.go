package config

import (
	hdconfig "github.com/nodeset-org/hyperdrive-daemon/shared/config"
	"github.com/nodeset-org/hyperdrive-daemon/shared/config/ids"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/rocket-pool/node-manager-core/config"
)

// The page wrapper for the Hyperdrive config
type HyperdriveConfigPage struct {
	home           *settingsHome
	page           *page
	layout         *standardLayout
	masterConfig   *client.GlobalConfig
	beforeItems    []*parameterizedFormItem
	txModeBox      *parameterizedFormItem
	txCustomUrlBox *parameterizedFormItem
	afterItems     []*parameterizedFormItem
}

// Creates a new page for the Hyperdrive settings
func NewHyperdriveConfigPage(home *settingsHome) *HyperdriveConfigPage {
	configPage := &HyperdriveConfigPage{
		home:         home,
		masterConfig: home.md.Config,
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
	beforeFormItems := []config.IParameter{}
	afterFormItems := []config.IParameter{}
	seenTxEnable := false
	for _, parameter := range masterConfig.Hyperdrive.GetParameters() {
		switch parameter.GetCommon().ID {
		case ids.ClientModeID:
			// Ignore the client mode item since it's presented in the EC / BN sections
			continue
		case ids.TxEndpointModeID:
			seenTxEnable = true
			continue
		case ids.TxCustomRpcUrlID:
			continue
		default:
			if !seenTxEnable {
				beforeFormItems = append(beforeFormItems, parameter)
			} else {
				afterFormItems = append(afterFormItems, parameter)
			}
		}
	}
	configPage.txCustomUrlBox = createParameterizedStringField(&masterConfig.Hyperdrive.TxCustomRpcUrl)
	configPage.beforeItems = createParameterizedFormItems(beforeFormItems, layout.descriptionBox)
	configPage.afterItems = createParameterizedFormItems(afterFormItems, layout.descriptionBox)

	// Map the parameters to the form items in the layout
	configPage.layout.mapParameterizedFormItems(configPage.beforeItems...)
	configPage.layout.mapParameterizedFormItems(configPage.afterItems...)
	configPage.layout.mapParameterizedFormItems(configPage.txCustomUrlBox)

	// Generate the TX mode dropdown box for the current network
	configPage.generateTxModeBox(masterConfig.Hyperdrive.Network.Value)

	// Set up the setting callbacks
	for _, item := range configPage.beforeItems {
		configPage.processNetworkParameter(item)
	}
	for _, item := range configPage.afterItems {
		configPage.processNetworkParameter(item)
	}

	// Do the initial draw
	configPage.handleLayoutChanged()
}

// Set up the handler for the Network dropdown
func (configPage *HyperdriveConfigPage) processNetworkParameter(item *parameterizedFormItem) {
	if item.parameter.GetCommon().ID != ids.NetworkID {
		return
	}
	dropDown := item.item.(*DropDown)
	dropDown.SetSelectedFunc(func(text string, index int) {
		newNetwork := configPage.home.md.Config.Hyperdrive.Network.Options[index].Value
		if newNetwork == configPage.home.md.Config.Hyperdrive.Network.Value {
			return
		}
		configPage.generateTxModeBox(newNetwork)
		configPage.txModeBox.item.(*DropDown).SetCurrentOption(0) // Reset the TX mode to ensure it's valid for the new network
		configPage.home.md.Config.ChangeNetwork(newNetwork)
		configPage.home.refresh()
	})
}

// Generate the TX mode dropdown box, filtering out the options that don't have URLs on the provided network
// The TUI wasn't really designed to handle dropdows with dynamic options, so we have to do some manual work here
func (configPage *HyperdriveConfigPage) generateTxModeBox(network config.Network) {
	res := configPage.masterConfig.AllNetworkSettings.Hyperdrive[network].NetworkResources

	// Filter out the options that don't have URLs on the provided network
	origOptions := configPage.masterConfig.Hyperdrive.TxEndpointMode.Options
	options := []*config.ParameterOption[hdconfig.TxEndpointMode]{}
	for _, option := range origOptions {
		switch option.Value {
		case hdconfig.TxEndpointMode_FlashbotsProtect:
			if res.FlashbotsProtectUrl == "" {
				continue
			}
		case hdconfig.TxEndpointMode_MevBlocker:
			if res.MevBlockerUrl == "" {
				continue
			}
		}
		options = append(options, option)
	}

	// Make a clone of the param with the new options
	paramCopy := config.Parameter[hdconfig.TxEndpointMode]{
		ParameterCommon: configPage.masterConfig.Hyperdrive.TxEndpointMode.GetCommon(),
		Value:           configPage.masterConfig.Hyperdrive.TxEndpointMode.Value,
		Default:         configPage.masterConfig.Hyperdrive.TxEndpointMode.Default,
		Options:         options,
	}

	// Create the dropdown box using the cloned parameter with the filtered options
	box := createParameterizedDropDown(&paramCopy, configPage.layout.descriptionBox)
	box.item.(*DropDown).SetSelectedFunc(func(text string, index int) {
		selection := options[index].GetValueAsAny().(hdconfig.TxEndpointMode)
		if configPage.masterConfig.Hyperdrive.TxEndpointMode.Value == selection {
			return
		}
		// Update both the real parameter and the copy so it displays correctly
		configPage.masterConfig.Hyperdrive.TxEndpointMode.Value = selection
		paramCopy.Value = selection
		configPage.handleModeChanged()
	})
	configPage.layout.mapParameterizedFormItems(box)
	configPage.txModeBox = box
}

// Handle all of the form changes when the TX mode has changed
func (configPage *HyperdriveConfigPage) handleModeChanged() {
	configPage.layout.form.Clear(true)

	// Add the stuff before the TX mode
	for _, formItem := range configPage.beforeItems {
		configPage.layout.form.AddFormItem(formItem.item)
	}

	// Add the TX mode and custom URL if applicable
	configPage.layout.form.AddFormItem(configPage.txModeBox.item)
	if configPage.masterConfig.Hyperdrive.TxEndpointMode.Value == hdconfig.TxEndpointMode_Custom {
		configPage.layout.form.AddFormItem(configPage.txCustomUrlBox.item)
	}

	// Add the stuff after the TX mode
	for _, formItem := range configPage.afterItems {
		configPage.layout.form.AddFormItem(formItem.item)
	}

	configPage.layout.refresh()
}

// Handle a bulk redraw request
func (configPage *HyperdriveConfigPage) handleLayoutChanged() {
	configPage.handleModeChanged()
}
