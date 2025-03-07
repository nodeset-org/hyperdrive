package config

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/modules"
)

// The page wrapper for the Hyperdrive config
type HyperdriveConfigPage struct {
	home      *settingsHome
	page      *page
	layout    *standardLayout
	formItems []*parameterizedFormItem
}

// Creates a new page for the Hyperdrive settings
func NewHyperdriveConfigPage(home *settingsHome) *HyperdriveConfigPage {
	configPage := &HyperdriveConfigPage{
		home: home,
	}

	configPage.createContent()
	configPage.page = newPage(
		home.homePage,
		"settings-hyperdrive",
		"Hyperdrive",
		"Select this to configure the settings for Hyperdrive itself.",
		configPage.layout.grid,
	)

	// Do the initial draw
	configPage.handleLayoutChanged()
	return configPage
}

// Get the underlying page
func (p *HyperdriveConfigPage) getPage() *page {
	return p.page
}

// Creates the content for the Hyperdrive settings page
func (p *HyperdriveConfigPage) createContent() {
	// Create the layout
	p.layout = newStandardLayout(p.home.md, modules.HyperdriveFqmn)
	p.layout.createForm("Hyperdrive Settings")
	p.layout.setupEscapeReturnHomeHandler(p.home.md, p.home.homePage)

	params := p.home.md.Config.GetParameters()
	newInstance := p.home.md.newInstance
	for _, param := range params {
		id := param.GetID()
		paramSetting, err := newInstance.GetParameter(id)
		if err != nil {
			panic(fmt.Errorf("error getting base parameter setting [%s]: %w", id, err))
		}

		// Create the form item for the parameter
		pfi := createParameterizedFormItem(paramSetting, p.layout, p.handleLayoutChanged)
		p.layout.registerFormItems(pfi)
		p.formItems = append(p.formItems, pfi)
	}
}

// Handle a bulk redraw request
func (p *HyperdriveConfigPage) handleLayoutChanged() {
	p.layout.redrawForm(p.formItems, nil, nil)
}
