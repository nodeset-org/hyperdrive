package config

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/nodeset-org/hyperdrive/modules/config"
	modconfig "github.com/nodeset-org/hyperdrive/modules/config"
)

type iSectionPage interface {
	getMainDisplay() *MainDisplay
	getPage() *page
	handleLayoutChanged()
}

type SectionPage struct {
	md     *MainDisplay
	parent iSectionPage
	page   *page
	layout *standardLayout

	section  config.ISection
	settings *modconfig.SettingsSection
	params   []*parameterizedFormItem
	subPages []iSectionPage
}

func NewSectionPage(md *MainDisplay, parent iSectionPage, section config.ISection, settings *modconfig.SettingsSection) *SectionPage {
	sectionPage := &SectionPage{
		md:       md,
		parent:   parent,
		section:  section,
		settings: settings,
	}
	sectionPage.createContent()

	// Create the section page
	sectionPage.page = newPage(
		parent.getPage(),
		parent.getPage().id+"/"+string(section.GetID()),
		string(section.GetName()),
		string(section.GetDescription().Default), // TEMPLATE
		sectionPage.layout.grid,
	)
	sectionPage.setupSubpages()

	// Do the initial draw
	sectionPage.handleLayoutChanged()
	return sectionPage
}

func (p *SectionPage) createContent() {
	// Create the layout
	p.layout = newStandardLayout()
	p.layout.createForm(string(p.section.GetName()))
	p.layout.setupEscapeReturnHomeHandler(p.md, p.parent.getPage())

	// Set up the params
	sectionCfg := p.section
	params := sectionCfg.GetParameters()
	settings := []modconfig.IParameterSetting{}
	for _, param := range params {
		id := param.GetID()
		setting, err := p.settings.GetParameter(id)
		if err != nil {
			panic(fmt.Errorf("error getting [%s] parameter setting [%s]: %w", p.section.GetName(), id, err)) // TODO: better logging, like FQMN
		}
		settings = append(settings, setting)
	}
	p.params = createParameterizedFormItems(settings, p.layout.descriptionBox)
	p.layout.mapParameterizedFormItems(p.params...)
}

// Set up the subpages
func (p *SectionPage) setupSubpages() {
	sectionCfg := p.section
	subsections := sectionCfg.GetSections()
	for _, section := range subsections {
		id := section.GetID()
		setting, err := p.settings.GetSection(id)
		if err != nil {
			panic(fmt.Errorf("error getting [%s] section setting [%s]: %w", p.section.GetName(), id, err)) // TODO: better logging, like FQMN
		}
		subPage := NewSectionPage(p.md, p, section, setting)
		p.subPages = append(p.subPages, subPage)

		// Map the description to the section label for button shifting later
		label := section.GetName()
		p.layout.mapButtonDescription(label, section.GetDescription())
	}
	for _, subpage := range p.subPages {
		p.md.pages.AddPage(subpage.getPage().id, subpage.getPage().content, true, false)
	}
}

// Get the main display
func (p *SectionPage) getMainDisplay() *MainDisplay {
	return p.md
}

// Get the underlying page
func (p *SectionPage) getPage() *page {
	return p.page
}

// Handle a bulk redraw request
func (p *SectionPage) handleLayoutChanged() {
	p.layout.form.Clear(true)
	p.layout.form.ClearButtons()

	p.layout.addFormItems(p.params)
	subsections := p.section.GetSections()
	for i, section := range subsections {
		subPage := p.subPages[i]
		p.layout.form.AddButton(section.GetName(), func() {
			subPage.handleLayoutChanged()
			p.md.setPage(subPage.getPage())
		})
		button := p.layout.form.GetButton(i)
		button.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyDown, tcell.KeyTab:
				return tcell.NewEventKey(tcell.KeyTab, 0, 0)
			case tcell.KeyUp, tcell.KeyBacktab:
				return tcell.NewEventKey(tcell.KeyBacktab, 0, 0)
			default:
				return event
			}
		})
	}

	p.layout.refresh()
}
