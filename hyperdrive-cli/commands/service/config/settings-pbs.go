package config

import (
	"github.com/nodeset-org/hyperdrive-daemon/shared/config/pbs"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/rivo/tview"
	"github.com/rocket-pool/node-manager-core/config"
)

const (
	mevWarning  string = "[orange]NOTE: You have externally-managed client mode selected and PBS enabled. You must have PBS enabled in your externally-managed Beacon Node's configuration for this to function properly - otherwise you may not be able to publish blocks and will miss significant rewards!"
	pbsDisabled string = "PBS clients are not available on this network yet, so PBS will be disabled. Check back at a later date once they've added support for it."
)

// The page wrapper for the MEV-boost config
type PbsConfigPage struct {
	home                  *settingsHome
	page                  *page
	layout                *standardLayout
	masterConfig          *client.GlobalConfig
	enableBox             *parameterizedFormItem
	modeBox               *parameterizedFormItem
	clientBox             *parameterizedFormItem
	selectionModeBox      *parameterizedFormItem
	localItems            []*parameterizedFormItem
	commitBoostItems      []*parameterizedFormItem
	mevBoostItems         []*parameterizedFormItem
	externalItems         []*parameterizedFormItem
	flashbotsBox          *parameterizedFormItem
	bloxrouteMaxProfitBox *parameterizedFormItem
	bloxrouteRegulatedBox *parameterizedFormItem
	titanRegionalBox      *parameterizedFormItem
}

// Creates a new page for the MEV-Boost settings
func NewPbsConfigPage(home *settingsHome) *PbsConfigPage {
	configPage := &PbsConfigPage{
		home:         home,
		masterConfig: home.md.Config,
	}
	configPage.createContent()

	configPage.page = newPage(
		home.homePage,
		"settings-pbs-client",
		"PBS Client",
		"Select this to configure the settings for your PBS Client, the source of blocks with MEV rewards for your validators.\n\nFor more information on PBS and PBS clients, please see https://writings.flashbots.net/writings/why-run-mevboost/",
		configPage.layout.grid,
	)

	return configPage
}

// Get the underlying page
func (configPage *PbsConfigPage) getPage() *page {
	return configPage.page
}

// Creates the content for the MEV-Boost settings page
func (configPage *PbsConfigPage) createContent() {
	// Create the layout
	configPage.layout = newStandardLayout()
	configPage.layout.createForm(&configPage.masterConfig.Hyperdrive.Network, "PBS Client Settings")
	configPage.layout.setupEscapeReturnHomeHandler(configPage.home.md, configPage.home.homePage)

	// Set up the form items
	configPage.enableBox = createParameterizedCheckbox(&configPage.masterConfig.Hyperdrive.Pbs.Enable)
	configPage.modeBox = createParameterizedDropDown(&configPage.masterConfig.Hyperdrive.Pbs.Mode, configPage.layout.descriptionBox)
	configPage.clientBox = createParameterizedDropDown(&configPage.masterConfig.Hyperdrive.Pbs.LocalPbsClientConfig.Client, configPage.layout.descriptionBox)
	configPage.selectionModeBox = createParameterizedDropDown(&configPage.masterConfig.Hyperdrive.Pbs.LocalPbsClientConfig.RelaySelectionMode, configPage.layout.descriptionBox)

	localParams := []config.IParameter{
		&configPage.masterConfig.Hyperdrive.Pbs.LocalPbsClientConfig.CustomRelays,
		&configPage.masterConfig.Hyperdrive.Pbs.LocalPbsClientConfig.Port,
		&configPage.masterConfig.Hyperdrive.Pbs.LocalPbsClientConfig.OpenRpcPort,
	}
	commitBoostParams := []config.IParameter{
		&configPage.masterConfig.Hyperdrive.Pbs.LocalPbsClientConfig.CommitBoostPbsConfig.ContainerTag,
		&configPage.masterConfig.Hyperdrive.Pbs.LocalPbsClientConfig.CommitBoostPbsConfig.AdditionalFlags,
	}
	mevBoostParams := []config.IParameter{
		&configPage.masterConfig.Hyperdrive.Pbs.LocalPbsClientConfig.MevBoostConfig.ContainerTag,
		&configPage.masterConfig.Hyperdrive.Pbs.LocalPbsClientConfig.MevBoostConfig.AdditionalFlags,
	}
	externalParams := []config.IParameter{&configPage.masterConfig.Hyperdrive.Pbs.ExternalPbsClientConfig.ExternalUrl}

	configPage.localItems = createParameterizedFormItems(localParams, configPage.layout.descriptionBox)
	configPage.commitBoostItems = createParameterizedFormItems(commitBoostParams, configPage.layout.descriptionBox)
	configPage.mevBoostItems = createParameterizedFormItems(mevBoostParams, configPage.layout.descriptionBox)
	configPage.externalItems = createParameterizedFormItems(externalParams, configPage.layout.descriptionBox)

	configPage.flashbotsBox = createParameterizedCheckbox(&configPage.masterConfig.Hyperdrive.Pbs.LocalPbsClientConfig.FlashbotsRelay)
	configPage.bloxrouteMaxProfitBox = createParameterizedCheckbox(&configPage.masterConfig.Hyperdrive.Pbs.LocalPbsClientConfig.BloxRouteMaxProfitRelay)
	configPage.bloxrouteRegulatedBox = createParameterizedCheckbox(&configPage.masterConfig.Hyperdrive.Pbs.LocalPbsClientConfig.BloxRouteRegulatedRelay)
	configPage.titanRegionalBox = createParameterizedCheckbox(&configPage.masterConfig.Hyperdrive.Pbs.LocalPbsClientConfig.TitanRegionalRelay)

	// Map the parameters to the form items in the layout
	configPage.layout.mapParameterizedFormItems(configPage.enableBox, configPage.modeBox, configPage.clientBox, configPage.selectionModeBox)
	configPage.layout.mapParameterizedFormItems(configPage.flashbotsBox, configPage.bloxrouteMaxProfitBox, configPage.bloxrouteRegulatedBox, configPage.titanRegionalBox)
	configPage.layout.mapParameterizedFormItems(configPage.localItems...)
	configPage.layout.mapParameterizedFormItems(configPage.commitBoostItems...)
	configPage.layout.mapParameterizedFormItems(configPage.mevBoostItems...)
	configPage.layout.mapParameterizedFormItems(configPage.externalItems...)

	// Set up the setting callbacks
	configPage.enableBox.item.(*tview.Checkbox).SetChangedFunc(func(checked bool) {
		if configPage.masterConfig.Hyperdrive.Pbs.Enable.Value == checked {
			return
		}
		configPage.masterConfig.Hyperdrive.Pbs.Enable.Value = checked
		configPage.handleLayoutChanged()
	})
	configPage.modeBox.item.(*DropDown).SetSelectedFunc(func(text string, index int) {
		if configPage.masterConfig.Hyperdrive.Pbs.Mode.Value == configPage.masterConfig.Hyperdrive.Pbs.Mode.Options[index].Value {
			return
		}
		configPage.masterConfig.Hyperdrive.Pbs.Mode.Value = configPage.masterConfig.Hyperdrive.Pbs.Mode.Options[index].Value
		configPage.handleModeChanged()
	})
	configPage.clientBox.item.(*DropDown).SetSelectedFunc(func(text string, index int) {
		if configPage.masterConfig.Hyperdrive.Pbs.LocalPbsClientConfig.Client.Value == configPage.masterConfig.Hyperdrive.Pbs.LocalPbsClientConfig.Client.Options[index].Value {
			return
		}
		configPage.masterConfig.Hyperdrive.Pbs.LocalPbsClientConfig.Client.Value = configPage.masterConfig.Hyperdrive.Pbs.LocalPbsClientConfig.Client.Options[index].Value
		configPage.handleRelaySelectionModeChanged()
	})
	configPage.selectionModeBox.item.(*DropDown).SetSelectedFunc(func(text string, index int) {
		if configPage.masterConfig.Hyperdrive.Pbs.LocalPbsClientConfig.RelaySelectionMode.Value == configPage.masterConfig.Hyperdrive.Pbs.LocalPbsClientConfig.RelaySelectionMode.Options[index].Value {
			return
		}
		configPage.masterConfig.Hyperdrive.Pbs.LocalPbsClientConfig.RelaySelectionMode.Value = configPage.masterConfig.Hyperdrive.Pbs.LocalPbsClientConfig.RelaySelectionMode.Options[index].Value
		configPage.handleRelaySelectionModeChanged()
	})

	// Do the initial draw
	configPage.handleLayoutChanged()
}

// Handle all of the form changes when the MEV-Boost mode has changed
func (configPage *PbsConfigPage) handleModeChanged() {
	configPage.layout.form.Clear(true)
	configPage.layout.form.AddFormItem(configPage.enableBox.item)
	if configPage.masterConfig.Hyperdrive.Pbs.Enable.Value {
		configPage.layout.form.AddFormItem(configPage.modeBox.item)

		selectedMode := configPage.masterConfig.Hyperdrive.Pbs.Mode.Value
		switch selectedMode {
		case config.ClientMode_Local:
			configPage.handleRelaySelectionModeChanged()
		case config.ClientMode_External:
			if configPage.masterConfig.Hyperdrive.IsLocalMode() {
				// Only show these to Docker users, not Hybrid users
				configPage.layout.addFormItems(configPage.externalItems)
			}
			configPage.layout.refresh()
		}
	}
}

// Handle all of the form changes when the relay selection mode has changed
func (configPage *PbsConfigPage) handleRelaySelectionModeChanged() {
	configPage.layout.form.Clear(true)
	configPage.layout.form.AddFormItem(configPage.enableBox.item)
	configPage.layout.form.AddFormItem(configPage.modeBox.item)
	configPage.layout.form.AddFormItem(configPage.clientBox.item)

	configPage.layout.form.AddFormItem(configPage.selectionModeBox.item)
	selectedMode := configPage.masterConfig.Hyperdrive.Pbs.LocalPbsClientConfig.RelaySelectionMode.Value
	switch selectedMode {
	case pbs.PbsRelaySelectionMode_Manual:
		relays := configPage.masterConfig.Hyperdrive.Pbs.LocalPbsClientConfig.GetAvailableRelays(configPage.masterConfig.Hyperdrive.GetEthNetworkName())
		for _, relay := range relays {
			switch relay.ID {
			case pbs.PbsRelayID_Flashbots:
				configPage.layout.form.AddFormItem(configPage.flashbotsBox.item)
			case pbs.PbsRelayID_BloxrouteMaxProfit:
				configPage.layout.form.AddFormItem(configPage.bloxrouteMaxProfitBox.item)
			case pbs.PbsRelayID_BloxrouteRegulated:
				configPage.layout.form.AddFormItem(configPage.bloxrouteRegulatedBox.item)
			case pbs.PbsRelayID_TitanRegional:
				configPage.layout.form.AddFormItem(configPage.titanRegionalBox.item)
			}
		}
	}

	configPage.layout.addFormItems(configPage.localItems)

	switch configPage.masterConfig.Hyperdrive.Pbs.LocalPbsClientConfig.Client.Value {
	case pbs.PbsClient_CommitBoost:
		configPage.layout.addFormItems(configPage.commitBoostItems)
	case pbs.PbsClient_MevBoost:
		configPage.layout.addFormItems(configPage.mevBoostItems)
	}

	configPage.layout.refresh()
}

// Handle a bulk redraw request
func (configPage *PbsConfigPage) handleLayoutChanged() {
	// Patch to add a PBS warning if in hybrid mode
	enableParam := configPage.masterConfig.Hyperdrive.Pbs.Enable
	if enableParam.DescriptionsByNetwork == nil {
		enableParam.DescriptionsByNetwork = map[config.Network]string{}
	}
	description := enableParam.Description
	augmentedDescription := description + "\n\n" + mevWarning
	for _, option := range configPage.masterConfig.Hyperdrive.Network.Options {
		network := option.Value
		if configPage.masterConfig.Hyperdrive.ClientMode.Value == config.ClientMode_External {
			enableParam.DescriptionsByNetwork[network] = augmentedDescription
		} else {
			enableParam.DescriptionsByNetwork[network] = description
		}
	}

	// Rebuild the parameter maps based on the selected network
	configPage.handleModeChanged()
}
