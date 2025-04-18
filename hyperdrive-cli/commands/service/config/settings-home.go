package config

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rocket-pool/node-manager-core/config"
)

const settingsHomeID string = "settings-home"

// This is a container for the primary settings category selection home screen.
type settingsHome struct {
	homePage         *page
	saveButton       *tview.Button
	wizardButton     *tview.Button
	hyperdrivePage   *HyperdriveConfigPage
	loggingPage      *LoggingConfigPage
	ecPage           *ExecutionConfigPage
	fallbackPage     *FallbackConfigPage
	bnPage           *BeaconConfigPage
	mevBoostPage     *MevBoostConfigPage
	metricsPage      *MetricsConfigPage
	modulesPage      *ModulesPage
	categoryList     *tview.List
	settingsSubpages []settingsPage
	content          tview.Primitive
	md               *mainDisplay
	warningModal     *choiceModalLayout
}

// Creates a new SettingsHome instance and adds (and its subpages) it to the main display.
func newSettingsHome(md *mainDisplay) *settingsHome {
	homePage := newPage(nil, settingsHomeID, "Categories", "", nil)

	// Create the page and return it
	home := &settingsHome{
		md:       md,
		homePage: homePage,
	}

	// Create the settings subpages
	home.hyperdrivePage = NewHyperdriveConfigPage(home)
	home.loggingPage = NewLoggingConfigPage(home)
	home.ecPage = NewExecutionConfigPage(home)
	home.bnPage = NewBeaconConfigPage(home)
	home.fallbackPage = NewFallbackConfigPage(home)
	home.mevBoostPage = NewMevBoostConfigPage(home)
	home.metricsPage = NewMetricsConfigPage(home)
	home.modulesPage = NewModulesPage(home)
	settingsSubpages := []settingsPage{
		home.hyperdrivePage,
		home.loggingPage,
		home.ecPage,
		home.bnPage,
		home.fallbackPage,
		home.mevBoostPage,
		home.metricsPage,
		home.modulesPage,
	}
	home.settingsSubpages = settingsSubpages

	// Add the subpages to the main display
	for _, subpage := range settingsSubpages {
		md.pages.AddPage(subpage.getPage().id, subpage.getPage().content, true, false)
	}
	home.createContent()
	homePage.content = home.content
	md.pages.AddPage(homePage.id, home.content, true, false)

	// Make the MEV-Boost warning
	home.createMevWarningModal()
	return home
}

// Create the content for this page
func (home *settingsHome) createContent() {

	layout := newStandardLayout()

	// Create the category list
	categoryList := tview.NewList().
		SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
			// Disable MEV-Boost on unsupported networks
			if mainText == home.mevBoostPage.page.title {
				if !home.md.Config.Hyperdrive.MevBoost.HasRelays() {
					layout.descriptionBox.SetText("MEV-Boost is not available on this network.")
					return
				}
			}

			// Disable Modules on unsupported networks
			modulesDisabled := false
			// modulesDisabled = (home.md.Config.Hyperdrive.Network.Value == config.Network_Hoodi)
			if mainText == home.modulesPage.page.title {
				if modulesDisabled {
					layout.descriptionBox.SetText("Modules are not currently available on this network.")
					return
				}
			}

			layout.descriptionBox.SetText(home.settingsSubpages[index].getPage().description)
		})
	categoryList.SetBackgroundColor(BackgroundColor)
	categoryList.SetBorderPadding(0, 0, 1, 1)

	home.categoryList = categoryList

	// Set tab to switch to the save and quit buttons
	categoryList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab || event.Key() == tcell.KeyBacktab {
			home.md.app.SetFocus(home.saveButton)
			return nil
		}
		return event
	})

	// Add all of the subpages to the list
	for _, subpage := range home.settingsSubpages {
		categoryList.AddItem(subpage.getPage().title, "", 0, nil)
	}
	categoryList.SetSelectedFunc(func(i int, s1, s2 string, r rune) {
		// Disable MEV-Boost on unsupported networks
		if home.settingsSubpages[i].getPage().title == home.mevBoostPage.page.title {
			if !home.md.Config.Hyperdrive.MevBoost.HasRelays() {
				return
			}
		}

		// Disable Modules on unsupported networks
		modulesDisabled := false
		// modulesDisabled = (home.md.Config.Hyperdrive.Network.Value == config.Network_Hoodi)
		if home.settingsSubpages[i].getPage().title == home.modulesPage.page.title {
			if modulesDisabled {
				return
			}
		}

		home.settingsSubpages[i].handleLayoutChanged()
		home.md.setPage(home.settingsSubpages[i].getPage())
	})

	// Make it the content of the layout and set the default description text
	layout.setContent(categoryList, categoryList.Box, "Select a Category")
	layout.descriptionBox.SetText(home.settingsSubpages[0].getPage().description)

	// Make the footer
	footer, height := home.createFooter()
	layout.setFooter(footer, height)

	// Set the home page's content to be the standard layout's grid
	home.content = layout.grid

}

// Create the footer, including the nav bar and the save / quit buttons
func (home *settingsHome) createFooter() (tview.Primitive, int) {

	// Nav bar
	navString1 := "Arrow keys: Navigate             Space/Enter: Select"
	navTextView1 := tview.NewTextView().
		SetDynamicColors(false).
		SetRegions(false).
		SetWrap(false)
	navBar1 := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(navTextView1, len(navString1), 1, false).
		AddItem(nil, 0, 1, false)
	fmt.Fprint(navTextView1, navString1)

	navString2 := "Tab: Go to the Buttons   Ctrl+C: Quit without Saving"
	navTextView2 := tview.NewTextView().
		SetDynamicColors(false).
		SetRegions(false).
		SetWrap(false)
	navBar2 := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(navTextView2, len(navString2), 1, false).
		AddItem(nil, 0, 1, false)
	fmt.Fprint(navTextView2, navString2)

	// Save and Quit buttons
	saveButton := tview.NewButton("Review Changes and Save")
	wizardButton := tview.NewButton("Open the Config Wizard")
	home.saveButton = saveButton
	home.wizardButton = wizardButton

	saveButton.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab || event.Key() == tcell.KeyBacktab {
			home.md.app.SetFocus(home.categoryList)
			return nil
		} else if event.Key() == tcell.KeyRight ||
			event.Key() == tcell.KeyLeft ||
			event.Key() == tcell.KeyUp ||
			event.Key() == tcell.KeyDown {
			home.md.app.SetFocus(wizardButton)
			return nil
		}
		return event
	})
	saveButton.SetSelectedFunc(func() {
		// Show the MEV-Boost warning modal first if applicable
		if home.md.Config.Hyperdrive.MevBoost.Enable.Value && home.md.Config.Hyperdrive.ClientMode.Value == config.ClientMode_External {
			home.showMevWarningModal()
		} else {
			home.showReviewPage()
		}
	})
	saveButton.SetStyle(tcell.StyleDefault.Background(HomeButtonUnfocusedBackgroundColor).Foreground(HomeButtonUnfocusedTextColor))
	//saveButton.SetDisabledStyle(tcell.StyleDefault.Background(HomeButtonUnfocusedBackgroundColor).Foreground(HomeButtonUnfocusedTextColor))
	saveButton.SetActivatedStyle(tcell.StyleDefault.Background(ButtonFocusedBackgroundColor).Foreground(ButtonFocusedTextColor))

	wizardButton.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab || event.Key() == tcell.KeyBacktab {
			home.md.app.SetFocus(home.categoryList)
			return nil
		} else if event.Key() == tcell.KeyRight ||
			event.Key() == tcell.KeyLeft ||
			event.Key() == tcell.KeyUp ||
			event.Key() == tcell.KeyDown {
			home.md.app.SetFocus(saveButton)
			return nil
		}
		return event
	})
	wizardButton.SetSelectedFunc(func() {
		home.md.wizard.welcomeModal.show()
	})
	wizardButton.SetStyle(tcell.StyleDefault.Background(HomeButtonUnfocusedBackgroundColor).Foreground(HomeButtonUnfocusedTextColor))
	//wizardButton.SetDisabledStyle(tcell.StyleDefault.Background(HomeButtonUnfocusedBackgroundColor).Foreground(HomeButtonUnfocusedTextColor))
	wizardButton.SetActivatedStyle(tcell.StyleDefault.Background(ButtonFocusedBackgroundColor).Foreground(ButtonFocusedTextColor))

	// Create overall layout for the footer
	buttonBar := tview.NewFlex().
		AddItem(nil, 0, 3, false).
		AddItem(saveButton, 25, 1, false).
		AddItem(nil, 0, 1, false).
		AddItem(wizardButton, 24, 1, false).
		AddItem(nil, 0, 3, false)

	footer := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(buttonBar, 1, 1, false).
		AddItem(nil, 1, 1, false).
		AddItem(navBar1, 1, 1, false).
		AddItem(navBar2, 1, 1, false)

	return footer, footer.GetItemCount()

}

// Refreshes the settings on all of the config pages to match the config's values
func (home *settingsHome) refresh() {
	/*
		if home.hyperdrivePage != nil {
			home.hyperdrivePage.layout.refresh()
		}*/

	if home.loggingPage != nil {
		home.loggingPage.layout.refresh()
	}

	if home.ecPage != nil {
		home.ecPage.layout.refresh()
	}

	if home.bnPage != nil {
		home.bnPage.layout.refresh()
	}

	if home.fallbackPage != nil {
		home.fallbackPage.layout.refresh()
	}

	if home.fallbackPage != nil {
		home.fallbackPage.layout.refresh()
	}

	if home.mevBoostPage != nil {
		home.mevBoostPage.layout.refresh()
	}

	if home.metricsPage != nil {
		home.metricsPage.layout.refresh()
	}
}

// Shows the review page
func (home *settingsHome) showReviewPage() {
	home.md.pages.RemovePage(reviewPageID)
	reviewPage := NewReviewPage(home.md, home.md.PreviousConfig, home.md.Config)
	home.md.pages.AddPage(reviewPage.page.id, reviewPage.page.content, true, true)
	home.md.setPage(reviewPage.page)
}

func (home *settingsHome) createMevWarningModal() {
	// Create the modal
	modal := newChoiceModalLayout(
		home.md.app,
		"MEV-Boost",
		76,
		mevWarning,
		[]string{"Ok"},
		[]string{},
		DirectionalModalHorizontal,
	)

	back := func() {
		home.showReviewPage()
	}

	done := func(buttonIndex int, buttonLabel string) {
		back()
	}

	modal.done = done
	modal.back = back

	page := newPage(
		home.homePage,
		"mev-warning-modal",
		"",
		"",
		modal.borderGrid,
	)
	home.md.pages.AddPage(page.id, page.content, true, false)
	modal.page = page
	home.warningModal = modal
}

func (home *settingsHome) showMevWarningModal() {
	home.md.setPage(home.warningModal.page)
	home.warningModal.focus(0)
}
