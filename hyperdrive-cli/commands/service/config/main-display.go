package config

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/nodeset-org/hyperdrive/shared"
	"github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/nodeset-org/hyperdrive/shared/types"
	"github.com/rivo/tview"
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
	PreviousConfig      *config.HyperdriveConfig
	Config              *config.HyperdriveConfig
	ShouldSave          bool
	ContainersToRestart []types.ContainerID
	ChangeNetworks      bool
}

// Creates a new MainDisplay instance.
func NewMainDisplay(app *tview.Application, previousConfig *config.HyperdriveConfig, config *config.HyperdriveConfig, isNew bool, isUpdate bool) *mainDisplay {
	// Create a copy of the original config for comparison purposes
	if previousConfig == nil {
		previousConfig = config.CreateCopy()
	}

	// Create the main grid
	grid := tview.NewGrid().
		SetColumns(1, 0, 1).   // 1-unit border
		SetRows(1, 1, 1, 0, 1) // Also 1-unit border

	grid.SetBorder(true).
		SetTitle(fmt.Sprintf(" Hyperdrive %s Configuration ", shared.HyperdriveVersion)).
		SetBorderColor(BorderColor).
		SetTitleColor(BorderColor).
		SetBackgroundColor(tcell.ColorBlack)

	// Create the navigation header
	navHeader := tview.NewTextView().
		SetDynamicColors(false).
		SetRegions(false).
		SetWrap(false)
	grid.AddItem(navHeader, 1, 1, 1, 1, 0, 0, false)

	// Create the page collection
	pages := tview.NewPages()
	grid.AddItem(pages, 3, 1, 1, 1, 0, 0, true)

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
