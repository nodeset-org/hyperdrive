package config

import (
	"github.com/gdamore/tcell/v2"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
)

// The page wrapper for the logging config
type LoggingConfigPage struct {
	home         *settingsHome
	page         *page
	layout       *standardLayout
	masterConfig *client.GlobalConfig
	loggingItems []*parameterizedFormItem
}

// Creates a new page for the logging settings
func NewLoggingConfigPage(home *settingsHome) *LoggingConfigPage {
	configPage := &LoggingConfigPage{
		home:         home,
		masterConfig: home.md.Config,
	}
	configPage.createContent()

	configPage.page = newPage(
		home.homePage,
		"settings-logging",
		"Logging",
		"Configure Hyperdrive's daemon and module logs.",
		configPage.layout.grid,
	)

	return configPage
}

// Get the underlying page
func (configPage *LoggingConfigPage) getPage() *page {
	return configPage.page
}

// Creates the content for the logging settings page
func (configPage *LoggingConfigPage) createContent() {
	// Create the layout
	configPage.layout = newStandardLayout()
	configPage.layout.createForm(&configPage.masterConfig.Hyperdrive.Network, "Logging Settings")

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
	configPage.loggingItems = createParameterizedFormItems(configPage.masterConfig.Hyperdrive.Logging.GetParameters(), configPage.layout.descriptionBox)

	// Map the parameters to the form items in the layout
	configPage.layout.mapParameterizedFormItems(configPage.loggingItems...)

	// Do the initial draw
	configPage.handleLayoutChanged()
}

// Handle a bulk redraw request
func (configPage *LoggingConfigPage) handleLayoutChanged() {
	configPage.layout.form.Clear(true)
	configPage.layout.addFormItems(configPage.loggingItems)
	configPage.layout.refresh()
}
