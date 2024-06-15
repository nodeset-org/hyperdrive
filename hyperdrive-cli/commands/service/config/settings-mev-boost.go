package config

import (
	hdconfig "github.com/nodeset-org/hyperdrive-daemon/shared/config"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/rivo/tview"
	"github.com/rocket-pool/node-manager-core/config"
)

// The page wrapper for the MEV-boost config
type MevBoostConfigPage struct {
	home                  *settingsHome
	page                  *page
	layout                *standardLayout
	masterConfig          *client.GlobalConfig
	enableBox             *parameterizedFormItem
	modeBox               *parameterizedFormItem
	selectionModeBox      *parameterizedFormItem
	localItems            []*parameterizedFormItem
	externalItems         []*parameterizedFormItem
	flashbotsBox          *parameterizedFormItem
	bloxrouteRegulatedBox *parameterizedFormItem
	edenBox               *parameterizedFormItem
	titanRegionalBox      *parameterizedFormItem
	customRelaysBox       *parameterizedFormItem
}

// Creates a new page for the MEV-Boost settings
func NewMevBoostConfigPage(home *settingsHome) *MevBoostConfigPage {
	configPage := &MevBoostConfigPage{
		home:         home,
		masterConfig: home.md.Config,
	}
	configPage.createContent()

	configPage.page = newPage(
		home.homePage,
		"settings-mev-boost",
		"MEV-Boost",
		"Select this to configure the settings for the Flashbots MEV-Boost client, the source of blocks with MEV rewards for your minipools.\n\nFor more information on Flashbots, MEV, and MEV-Boost, please see https://writings.flashbots.net/writings/why-run-mevboost/",
		configPage.layout.grid,
	)

	return configPage
}

// Get the underlying page
func (configPage *MevBoostConfigPage) getPage() *page {
	return configPage.page
}

// Creates the content for the MEV-Boost settings page
func (configPage *MevBoostConfigPage) createContent() {
	// Create the layout
	configPage.layout = newStandardLayout()
	configPage.layout.createForm(&configPage.masterConfig.Hyperdrive.Network, "MEV-Boost Settings")
	configPage.layout.setupEscapeReturnHomeHandler(configPage.home.md, configPage.home.homePage)

	// Set up the form items
	configPage.enableBox = createParameterizedCheckbox(&configPage.masterConfig.Hyperdrive.MevBoost.Enable)
	configPage.modeBox = createParameterizedDropDown(&configPage.masterConfig.Hyperdrive.MevBoost.Mode, configPage.layout.descriptionBox)
	configPage.selectionModeBox = createParameterizedDropDown(&configPage.masterConfig.Hyperdrive.MevBoost.SelectionMode, configPage.layout.descriptionBox)

	localParams := []config.IParameter{
		&configPage.masterConfig.Hyperdrive.MevBoost.CustomRelays,
		&configPage.masterConfig.Hyperdrive.MevBoost.Port,
		&configPage.masterConfig.Hyperdrive.MevBoost.OpenRpcPort,
		&configPage.masterConfig.Hyperdrive.MevBoost.ContainerTag,
		&configPage.masterConfig.Hyperdrive.MevBoost.AdditionalFlags,
	}
	externalParams := []config.IParameter{&configPage.masterConfig.Hyperdrive.MevBoost.ExternalUrl}

	configPage.localItems = createParameterizedFormItems(localParams, configPage.layout.descriptionBox)
	configPage.externalItems = createParameterizedFormItems(externalParams, configPage.layout.descriptionBox)

	configPage.flashbotsBox = createParameterizedCheckbox(&configPage.masterConfig.Hyperdrive.MevBoost.FlashbotsRelay)
	configPage.bloxrouteRegulatedBox = createParameterizedCheckbox(&configPage.masterConfig.Hyperdrive.MevBoost.BloxRouteRegulatedRelay)
	configPage.edenBox = createParameterizedCheckbox(&configPage.masterConfig.Hyperdrive.MevBoost.EdenRelay)
	configPage.titanRegionalBox = createParameterizedCheckbox(&configPage.masterConfig.Hyperdrive.MevBoost.TitanRegionalRelay)

	// Map the parameters to the form items in the layout
	configPage.layout.mapParameterizedFormItems(configPage.enableBox, configPage.modeBox, configPage.selectionModeBox)
	configPage.layout.mapParameterizedFormItems(configPage.flashbotsBox, configPage.bloxrouteRegulatedBox, configPage.edenBox, configPage.titanRegionalBox)
	configPage.layout.mapParameterizedFormItems(configPage.localItems...)
	configPage.layout.mapParameterizedFormItems(configPage.externalItems...)

	// Set up the setting callbacks
	configPage.enableBox.item.(*tview.Checkbox).SetChangedFunc(func(checked bool) {
		if configPage.masterConfig.Hyperdrive.MevBoost.Enable.Value == checked {
			return
		}
		configPage.masterConfig.Hyperdrive.MevBoost.Enable.Value = checked
		configPage.handleLayoutChanged()
	})
	configPage.modeBox.item.(*DropDown).SetSelectedFunc(func(text string, index int) {
		if configPage.masterConfig.Hyperdrive.MevBoost.Mode.Value == configPage.masterConfig.Hyperdrive.MevBoost.Mode.Options[index].Value {
			return
		}
		configPage.masterConfig.Hyperdrive.MevBoost.Mode.Value = configPage.masterConfig.Hyperdrive.MevBoost.Mode.Options[index].Value
		configPage.handleModeChanged()
	})
	configPage.selectionModeBox.item.(*DropDown).SetSelectedFunc(func(text string, index int) {
		if configPage.masterConfig.Hyperdrive.MevBoost.SelectionMode.Value == configPage.masterConfig.Hyperdrive.MevBoost.SelectionMode.Options[index].Value {
			return
		}
		configPage.masterConfig.Hyperdrive.MevBoost.SelectionMode.Value = configPage.masterConfig.Hyperdrive.MevBoost.SelectionMode.Options[index].Value
		configPage.handleSelectionModeChanged()
	})

	// Do the initial draw
	configPage.handleLayoutChanged()
}

// Handle all of the form changes when the MEV-Boost mode has changed
func (configPage *MevBoostConfigPage) handleModeChanged() {
	configPage.layout.form.Clear(true)
	configPage.layout.form.AddFormItem(configPage.enableBox.item)
	if configPage.masterConfig.Hyperdrive.MevBoost.Enable.Value {
		configPage.layout.form.AddFormItem(configPage.modeBox.item)

		selectedMode := configPage.masterConfig.Hyperdrive.MevBoost.Mode.Value
		switch selectedMode {
		case config.ClientMode_Local:
			configPage.handleSelectionModeChanged()
		case config.ClientMode_External:
			if configPage.masterConfig.Hyperdrive.IsLocalMode() {
				// Only show these to Docker users, not Hybrid users
				configPage.layout.addFormItems(configPage.externalItems)
			}
		}
	}

	configPage.layout.refresh()
}

// Handle all of the form changes when the relay selection mode has changed
func (configPage *MevBoostConfigPage) handleSelectionModeChanged() {
	configPage.layout.form.Clear(true)
	configPage.layout.form.AddFormItem(configPage.enableBox.item)
	configPage.layout.form.AddFormItem(configPage.modeBox.item)

	configPage.layout.form.AddFormItem(configPage.selectionModeBox.item)
	selectedMode := configPage.masterConfig.Hyperdrive.MevBoost.SelectionMode.Value
	switch selectedMode {
	case hdconfig.MevSelectionMode_Manual:
		relays := configPage.masterConfig.Hyperdrive.MevBoost.GetAvailableRelays()
		for _, relay := range relays {
			switch relay.ID {
			case hdconfig.MevRelayID_Flashbots:
				configPage.layout.form.AddFormItem(configPage.flashbotsBox.item)
			case hdconfig.MevRelayID_BloxrouteRegulated:
				configPage.layout.form.AddFormItem(configPage.bloxrouteRegulatedBox.item)
			case hdconfig.MevRelayID_Eden:
				configPage.layout.form.AddFormItem(configPage.edenBox.item)
			case hdconfig.MevRelayID_TitanRegional:
				configPage.layout.form.AddFormItem(configPage.titanRegionalBox.item)
			}
		}
	}

	configPage.layout.addFormItems(configPage.localItems)
}

// Handle a bulk redraw request
func (configPage *MevBoostConfigPage) handleLayoutChanged() {
	// Rebuild the parameter maps based on the selected network
	configPage.handleModeChanged()
}
