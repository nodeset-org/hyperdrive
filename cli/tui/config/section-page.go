package config

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/modules/config"
	modconfig "github.com/nodeset-org/hyperdrive/modules/config"
)

type iSectionPage interface {
	getPage() *page
	handleLayoutChanged()
}

type SectionPage struct {
	md        *MainDisplay
	parent    iSectionPage
	page      *page
	layout    *standardLayout
	fqmn      string
	subPages  []iSectionPage
	section   config.ISection
	settings  *modconfig.SettingsSection
	formItems []*parameterizedFormItem
	buttons   []*metadataButton
}

func NewSectionPage(md *MainDisplay, parent iSectionPage, section config.ISection, settings *modconfig.SettingsSection, fqmn string) *SectionPage {
	sectionPage := &SectionPage{
		md:       md,
		parent:   parent,
		section:  section,
		settings: settings,
		fqmn:     fqmn,
	}
	sectionPage.createContent()

	// Do the initial draw
	sectionPage.handleLayoutChanged()
	return sectionPage
}

func (p *SectionPage) createContent() {
	// Create the layout
	md := p.md
	p.layout = newStandardLayout(md, p.fqmn)
	p.layout.createForm(string(p.section.GetName()))
	p.layout.setupEscapeReturnHomeHandler(p.md, p.parent.getPage())

	// Create the section page
	p.page = newPage(
		p.parent.getPage(),
		p.parent.getPage().id+"/"+string(p.section.GetID()),
		string(p.section.GetName()),
		"", // Don't need a description since it's handled by the button that opens this page
		p.layout.grid,
	)

	// Set up the form items
	sectionCfg := p.section
	params := sectionCfg.GetParameters()
	for _, param := range params {
		id := param.GetID()
		paramSetting, err := p.settings.GetParameter(id)
		if err != nil {
			panic(fmt.Errorf("error getting \"%s\" parameter setting \"%s\": %w", p.section.GetName(), id, err)) // TODO: better logging, like FQMN
		}

		// Create the form item for the parameter
		pfi := createParameterizedFormItem(paramSetting, p.layout, p.handleLayoutChanged)
		p.layout.registerFormItems(pfi)
		p.formItems = append(p.formItems, pfi)
	}

	// Set up the section subpages
	subsections := sectionCfg.GetSections()
	for _, section := range subsections {
		id := section.GetID()
		settingsSection, err := p.settings.GetSection(id)
		if err != nil {
			panic(fmt.Errorf("error getting \"%s\" section setting \"%s\": %w", p.section.GetName(), id, err)) // TODO: better logging, like FQMN
		}

		// Create the subpage
		subPage := NewSectionPage(p.md, p, section, settingsSection, p.fqmn)
		p.subPages = append(p.subPages, subPage)
		md.pages.AddPage(subPage.getPage().id, subPage.getPage().content, true, false)

		// Create the metadata button for the section
		button := createMetadataButton(section, subPage, md)
		p.layout.registerButtons(button)
		p.buttons = append(p.buttons, button)
	}
}

// Get the underlying page
func (p *SectionPage) getPage() *page {
	return p.page
}

// Handle a bulk redraw request
func (p *SectionPage) handleLayoutChanged() {
	p.layout.redrawForm(p.formItems, p.buttons, nil)
}
