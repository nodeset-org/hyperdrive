package config

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/nodeset-org/hyperdrive-daemon/shared"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/rivo/tview"
	"github.com/rocket-pool/node-manager-core/config"
)

// This represents the primary TUI for the configuration command
type mainDisplay struct {
	navHeader           *tview.TextView
	pages               *tview.Pages
	app                 *tview.Application
	content             *tview.Box
	mainGrid            *tview.Grid
	wizard              *wizard
	settingsHome        *settingsHome
	isNew               bool
	isUpdate            bool
	previousWidth       int
	previousHeight      int
	PreviousConfig      *client.GlobalConfig
	Config              *client.GlobalConfig
	ShouldSave          bool
	ContainersToRestart []config.ContainerID
	ChangeNetworks      bool
}

// Creates a new MainDisplay instance.
func NewMainDisplay(app *tview.Application, previousConfig *client.GlobalConfig, config *client.GlobalConfig, isNew bool, isUpdate bool) *mainDisplay {
	// Create a copy of the original config for comparison purposes
	if previousConfig == nil {
		previousConfig = config.CreateCopy()
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
	md := &mainDisplay{
		navHeader:      navHeader,
		pages:          pages,
		app:            app,
		content:        grid.Box,
		mainGrid:       grid,
		isNew:          isNew,
		isUpdate:       isUpdate,
		PreviousConfig: previousConfig,
		Config:         config,
	}

	// Create all of the child elements
	md.settingsHome = newSettingsHome(md)
	md.wizard = newWizard(md)

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

	if isNew {
		md.wizard.welcomeModal.show()
	} else {
		md.setPage(md.settingsHome.homePage)
	}
	app.SetRoot(grid, true)
	return md
}

// Sets the current page that is on display.
func (md *mainDisplay) setPage(page *page) {
	md.navHeader.SetText(page.getHeader())
	md.pages.SwitchToPage(page.id)
}
