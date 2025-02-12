package config

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	modconfig "github.com/nodeset-org/hyperdrive/modules/config"
	"github.com/nodeset-org/hyperdrive/shared"
	"github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/nodeset-org/hyperdrive/shared/utils"
	"github.com/rivo/tview"
)

// This represents the primary TUI for the configuration command
type MainDisplay struct {
	navHeader *tview.TextView
	pages     *tview.Pages
	app       *tview.Application
	content   *tview.Box
	mainGrid  *tview.Grid
	//wizard              *wizard
	settingsHome        *settingsHome
	isNew               bool
	isUpdate            bool
	previousWidth       int
	previousHeight      int
	Config              *config.HyperdriveConfig
	PreviousSettings    *config.HyperdriveSettings
	NewSettings         *config.HyperdriveSettings
	ShouldSave          bool
	ContainersToRestart []config.ContainerID
	ChangeNetworks      bool

	// Private fields
	previousInstance *modconfig.ModuleSettings
	newInstance      *modconfig.ModuleSettings
	moduleManager    *utils.ModuleManager
}

// Creates a new MainDisplay instance.
func NewMainDisplay(
	app *tview.Application,
	moduleManager *utils.ModuleManager,
	config *config.HyperdriveConfig,
	previousSettings *config.HyperdriveSettings,
	newSettings *config.HyperdriveSettings,
	isNew bool,
	isUpdate bool,
) *MainDisplay {
	// Create a copy of the original settings for comparison purposes
	if previousSettings == nil {
		previousSettings = newSettings.CreateCopy()
	}

	// Create Hyperdrive settings instances
	previousInstance := modconfig.CreateModuleSettings(config)
	err := previousInstance.CopySettingsFromKnownType(previousSettings)
	if err != nil {
		panic(fmt.Errorf("error copying previous settings to HD config instance: %w", err))
	}
	newInstance := modconfig.CreateModuleSettings(config)
	err = newInstance.CopySettingsFromKnownType(newSettings)
	if err != nil {
		panic(fmt.Errorf("error copying new settings to HD config instance: %w", err))
	}

	// Create the main grid
	grid := tview.NewGrid().
		SetColumns(1, 0, 1).   // 1-unit border
		SetRows(1, 1, 1, 0, 1) // Also 1-unit border

	grid.SetBackgroundColor(NonInteractiveBackgroundColor)

	grid.SetBorder(true).
		SetTitle(fmt.Sprintf(" Hyperdrive %s Configuration ", shared.HyperdriveVersion)).
		SetBorderColor(BorderColor).
		SetTitleColor(BorderColor).
		SetBackgroundColor(NonInteractiveBackgroundColor)

	// Padding
	grid.AddItem(tview.NewBox().SetBackgroundColor(NonInteractiveBackgroundColor), 0, 0, 1, 1, 0, 0, false)
	grid.AddItem(tview.NewBox().SetBackgroundColor(NonInteractiveBackgroundColor), 0, 1, 1, 1, 0, 0, false)
	grid.AddItem(tview.NewBox().SetBackgroundColor(NonInteractiveBackgroundColor), 0, 2, 1, 1, 0, 0, false)

	// Create the navigation header
	navHeader := tview.NewTextView().
		SetDynamicColors(false).
		SetRegions(false).
		SetWrap(false)
	grid.AddItem(navHeader, 1, 1, 1, 1, 0, 0, false)
	grid.AddItem(tview.NewBox().SetBackgroundColor(NonInteractiveBackgroundColor), 1, 0, 1, 1, 0, 0, false)
	grid.AddItem(tview.NewBox().SetBackgroundColor(NonInteractiveBackgroundColor), 1, 2, 1, 1, 0, 0, false)

	// Padding
	grid.AddItem(tview.NewBox().SetBackgroundColor(NonInteractiveBackgroundColor), 2, 0, 1, 1, 0, 0, false)
	grid.AddItem(tview.NewBox().SetBackgroundColor(NonInteractiveBackgroundColor), 2, 1, 1, 1, 0, 0, false)
	grid.AddItem(tview.NewBox().SetBackgroundColor(NonInteractiveBackgroundColor), 2, 2, 1, 1, 0, 0, false)

	// Create the page collection
	pages := tview.NewPages()
	grid.AddItem(pages, 3, 1, 1, 1, 0, 0, true)
	grid.AddItem(tview.NewBox().SetBackgroundColor(NonInteractiveBackgroundColor), 3, 0, 1, 1, 0, 0, false)
	grid.AddItem(tview.NewBox().SetBackgroundColor(NonInteractiveBackgroundColor), 3, 2, 1, 1, 0, 0, false)

	// Padding
	grid.AddItem(tview.NewBox().SetBackgroundColor(NonInteractiveBackgroundColor), 4, 0, 1, 1, 0, 0, false)
	grid.AddItem(tview.NewBox().SetBackgroundColor(NonInteractiveBackgroundColor), 4, 1, 1, 1, 0, 0, false)
	grid.AddItem(tview.NewBox().SetBackgroundColor(NonInteractiveBackgroundColor), 4, 2, 1, 1, 0, 0, false)

	// Create the resize warning
	resizeWarning := tview.NewTextView().
		SetText("Your terminal is too small to run the service configuration app.\n\nPlease resize your terminal window and make it larger to see the app properly.").
		SetTextAlign(tview.AlignCenter).
		SetWordWrap(true).
		SetTextColor(tview.Styles.PrimaryTextColor)
	resizeWarning.SetBackgroundColor(BackgroundColor)
	resizeWarning.SetBorderPadding(0, 0, 1, 1)

	// Create the main display object
	md := &MainDisplay{
		navHeader:        navHeader,
		pages:            pages,
		app:              app,
		content:          grid.Box,
		mainGrid:         grid,
		isNew:            isNew,
		isUpdate:         isUpdate,
		moduleManager:    moduleManager,
		Config:           config,
		PreviousSettings: previousSettings,
		NewSettings:      newSettings,
		previousInstance: previousInstance,
		newInstance:      newInstance,
	}

	// Create all of the child elements
	md.settingsHome = newSettingsHome(md)
	//md.wizard = newWizard(md)

	// Set up the resize warning
	md.app.SetBeforeDrawFunc(func(screen tcell.Screen) bool {
		x, y := screen.Size()
		if x == md.previousWidth && y == md.previousHeight {
			return false
		}
		if x < 112 || y < 32 {
			grid.RemoveItem(pages)
			grid.AddItem(resizeWarning, 3, 1, 1, 1, 0, 0, false)
		} else {
			grid.RemoveItem(resizeWarning)
			grid.AddItem(pages, 3, 1, 1, 1, 0, 0, true)
		}
		md.previousWidth = x
		md.previousHeight = y

		return false
	})

	md.setPage(md.settingsHome.homePage)
	app.SetRoot(grid, true)
	return md
}

// Sets the current page that is on display.
func (md *MainDisplay) setPage(page *page) {
	md.navHeader.SetText(page.getHeader())
	md.pages.SwitchToPage(page.id)
}

func (md *MainDisplay) UpdateSettingsFromTuiSelections() error {
	// Copy the base settings
	err := md.newInstance.ConvertToKnownType(md.NewSettings)
	if err != nil {
		return fmt.Errorf("error converting updated base settings to known type: %w", err)
	}

	// Copy the module settings
	for _, modulePage := range md.settingsHome.modulesPage.moduleSubpages {
		modInstance := modulePage.instance
		modInstance.SetSettings(modulePage.settings)
	}
	return nil
}
