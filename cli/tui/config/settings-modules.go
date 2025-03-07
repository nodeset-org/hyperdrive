package config

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/nodeset-org/hyperdrive/config"
	"github.com/rivo/tview"
)

// Constants
const (
	modulesPageID string = "modules"
	modulePrefix  string = "module-"
)

// The addons page
type ModulesPage struct {
	home           *settingsHome
	page           *page
	layout         *standardLayout
	masterConfig   *config.HyperdriveConfig
	categoryList   *tview.List
	moduleSubpages []*ModulePage
}

// Create a new modules page
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

	// Create the subpages for each addon
	moduleSubpages := []*ModulePage{}
	cfg := home.md.Config
	for _, module := range cfg.Modules {
		fqmn := module.Descriptor.GetFullyQualifiedModuleName()
		previousInstance := home.md.PreviousSettings.Modules[fqmn]
		instance, exists := home.md.NewSettings.Modules[fqmn]
		if !exists {
			panic(fmt.Errorf("module instance [%s] not found", fqmn))
		}
		modulePage := NewModulePage(modulesPage, module, previousInstance, instance)
		moduleSubpages = append(moduleSubpages, modulePage)
	}
	modulesPage.moduleSubpages = moduleSubpages

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
	p.layout = newStandardLayout(p.home.md, "")
	p.layout.createSettingFooter()

	// Create the category list
	categoryList := tview.NewList().
		SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
			p.layout.descriptionBox.SetText(p.moduleSubpages[index].getPage().description)
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
	for _, subpage := range p.moduleSubpages {
		categoryList.AddItem("  "+subpage.getPage().title+"  ", "", 0, nil)
	}
	categoryList.SetSelectedFunc(func(i int, s1, s2 string, r rune) {
		p.moduleSubpages[i].handleLayoutChanged()
		p.home.md.setPage(p.moduleSubpages[i].getPage())
	})

	// Make it the content of the layout and set the default description text
	p.layout.setContent(categoryList, categoryList.Box, "Select a Module")
	if len(p.moduleSubpages) > 0 {
		p.layout.descriptionBox.SetText(p.moduleSubpages[0].getPage().description)
	}

	// Make the footer
	//footer, height := addonsPage.createFooter()
	//layout.setFooter(footer, height)

	// Set the home page's content to be the standard layout's grid
	//addonsPage.content = layout.grid
}

// Handle a bulk redraw request
func (p *ModulesPage) handleLayoutChanged() {

}
