package config

import (
	"github.com/gdamore/tcell/v2"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/rivo/tview"
)

// Constants
const modulesPageID string = "modules"

// The addons page
type ModulesPage struct {
	home          *settingsHome
	page          *page
	layout        *standardLayout
	masterConfig  *client.GlobalConfig
	stakewisePage *StakewiseConfigPage
	categoryList  *tview.List
	addonSubpages []settingsPage
}

// Create a new addons page
func NewModulesPage(home *settingsHome) *ModulesPage {

	modulesPage := &ModulesPage{
		home:         home,
		masterConfig: home.md.Config,
	}
	modulesPage.page = newPage(
		home.homePage,
		modulesPageID,
		"Modules",
		"Manage the different Hyperdrive modules, enabling and configuring the ones you want to use.",
		nil,
	)

	// Create the addon subpages
	modulesPage.stakewisePage = NewStakewiseConfigPage(modulesPage)
	moduleSubpages := []settingsPage{
		modulesPage.stakewisePage,
	}
	modulesPage.addonSubpages = moduleSubpages

	// Add the subpages to the main display
	for _, subpage := range moduleSubpages {
		home.md.pages.AddPage(subpage.getPage().id, subpage.getPage().content, true, false)
	}

	modulesPage.createContent()
	modulesPage.page.content = modulesPage.layout.grid
	return modulesPage

}

// Get the underlying page
func (p *ModulesPage) getPage() *page {
	return p.page
}

// Creates the content for the fallback client settings page
func (p *ModulesPage) createContent() {

	p.layout = newStandardLayout()
	p.layout.createSettingFooter()

	// Create the category list
	categoryList := tview.NewList().
		SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
			p.layout.descriptionBox.SetText(p.addonSubpages[index].getPage().description)
		})
	categoryList.SetBackgroundColor(BackgroundColor)
	categoryList.SetBorderPadding(0, 0, 1, 1)
	p.categoryList = categoryList

	// Set tab to switch to the save and quit buttons
	categoryList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			// Return to the home page
			p.home.md.setPage(p.home.homePage)
			return nil
		}
		return event
	})

	// Add all of the subpages to the list
	for _, subpage := range p.addonSubpages {
		categoryList.AddItem(subpage.getPage().title, "", 0, nil)
	}
	categoryList.SetSelectedFunc(func(i int, s1, s2 string, r rune) {
		p.addonSubpages[i].handleLayoutChanged()
		p.home.md.setPage(p.addonSubpages[i].getPage())
	})

	// Make it the content of the layout and set the default description text
	p.layout.setContent(categoryList, categoryList.Box, "Select an Addon")
	p.layout.descriptionBox.SetText(p.addonSubpages[0].getPage().description)

	// Make the footer
	//footer, height := addonsPage.createFooter()
	//layout.setFooter(footer, height)

	// Set the home page's content to be the standard layout's grid
	//addonsPage.content = layout.grid
}

// Handle a bulk redraw request
func (p *ModulesPage) handleLayoutChanged() {

}
