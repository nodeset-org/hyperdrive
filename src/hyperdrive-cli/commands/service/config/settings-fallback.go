package config

import (
	"github.com/gdamore/tcell/v2"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/rivo/tview"
	"github.com/rocket-pool/node-manager-core/config/ids"
)

// The page wrapper for the fallback config
type FallbackConfigPage struct {
	home           *settingsHome
	page           *page
	layout         *standardLayout
	masterConfig   *client.GlobalConfig
	useFallbackBox *parameterizedFormItem
	fallbackItems  []*parameterizedFormItem
}

// Creates a new page for the fallback client settings
func NewFallbackConfigPage(home *settingsHome) *FallbackConfigPage {
	configPage := &FallbackConfigPage{
		home:         home,
		masterConfig: home.md.Config,
	}
	configPage.createContent()

	configPage.page = newPage(
		home.homePage,
		"settings-fallback",
		"Fallback Clients",
		"Select this to specify a secondary, externally-managed Execution Client and Beacon Node pair.\nHyperdrive and any module's Validator Clients will use them as backups if your main Execution Client or Beacon Node ever go offline.",
		configPage.layout.grid,
	)

	return configPage
}

// Get the underlying page
func (configPage *FallbackConfigPage) getPage() *page {
	return configPage.page
}

// Creates the content for the fallback client settings page
func (configPage *FallbackConfigPage) createContent() {
	// Create the layout
	configPage.layout = newStandardLayout()
	configPage.layout.createForm(&configPage.masterConfig.Hyperdrive.Network, "Fallback Client Settings")

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
	configPage.useFallbackBox = createParameterizedCheckbox(&configPage.masterConfig.Hyperdrive.Fallback.UseFallbackClients)
	configPage.fallbackItems = createParameterizedFormItems(configPage.masterConfig.Hyperdrive.Fallback.GetParameters(), configPage.layout.descriptionBox)

	// Take the enable out since it's done explicitly
	fallbackItems := []*parameterizedFormItem{}
	for _, item := range configPage.fallbackItems {
		if item.parameter.GetCommon().ID == ids.FallbackUseFallbackClientsID {
			continue
		}
		fallbackItems = append(fallbackItems, item)
	}
	configPage.fallbackItems = fallbackItems

	// Map the parameters to the form items in the layout
	configPage.layout.mapParameterizedFormItems(configPage.useFallbackBox)
	configPage.layout.mapParameterizedFormItems(configPage.fallbackItems...)

	// Set up the setting callbacks
	configPage.useFallbackBox.item.(*tview.Checkbox).SetChangedFunc(func(checked bool) {
		if configPage.masterConfig.Hyperdrive.Fallback.UseFallbackClients.Value == checked {
			return
		}
		configPage.masterConfig.Hyperdrive.Fallback.UseFallbackClients.Value = checked
		configPage.handleUseFallbackChanged()
	})

	// Do the initial draw
	configPage.handleUseFallbackChanged()
}

// Handle all of the form changes when the Use Fallback box has changed
func (configPage *FallbackConfigPage) handleUseFallbackChanged() {
	configPage.layout.form.Clear(true)
	configPage.layout.form.AddFormItem(configPage.useFallbackBox.item)

	// Only add the supporting stuff if external clients are enabled
	if !configPage.masterConfig.Hyperdrive.Fallback.UseFallbackClients.Value {
		return
	}

	configPage.layout.addFormItems(configPage.fallbackItems)
	configPage.layout.refresh()
}

// Handle a bulk redraw request
func (configPage *FallbackConfigPage) handleLayoutChanged() {
	configPage.handleUseFallbackChanged()
}
