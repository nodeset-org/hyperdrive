package config

import (
	"fmt"
	"strconv"
	"strings"
	"text/template"

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
	fqmn   string

	section   config.ISection
	settings  *modconfig.SettingsSection
	params    []*parameterizedFormItem
	subPages  []iSectionPage
	redrawing bool
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

	// Create the section page
	sectionPage.page = newPage(
		parent.getPage(),
		parent.getPage().id+"/"+string(section.GetID()),
		string(section.GetName()),
		"", // string(section.GetDescription().Default), // TEMPLATE
		sectionPage.layout.grid,
	)
	sectionPage.setupSubpages()

	// Do the initial draw
	sectionPage.handleLayoutChanged()
	return sectionPage
}

func (p *SectionPage) createContent() {
	// Create the layout
	p.layout = newStandardLayout(p.getMainDisplay())
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
	p.params = createParameterizedFormItems(settings, p.layout.descriptionBox, p.handleLayoutChanged)
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
		subPage := NewSectionPage(p.md, p, section, setting, p.fqmn)
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
	if p.redrawing {
		return
	}
	p.redrawing = true
	defer func() {
		p.redrawing = false
	}()

	p.layout.form.Clear(true)

	params := []*parameterizedFormItem{}
	md := p.getMainDisplay()
	for _, param := range p.params {
		metadata := param.parameter.GetMetadata()
		hidden := metadata.GetHidden()

		// Handle parameters that don't have a hidden template
		if hidden.Template == "" {
			if !hidden.Default {
				params = append(params, param)
			}
			continue
		}

		// Generate a template source for the parameter
		templateSource := parameterTemplateSource{
			configurationTemplateSource: configurationTemplateSource{
				fqmn:              p.fqmn,
				hdSettings:        md.newInstance,
				moduleSettingsMap: md.moduleSettingsMap,
			},
			parameter: param.parameter.GetMetadata(),
		}

		// Update the hidden status
		template, err := template.New(string(metadata.GetID())).Parse(hidden.Template)
		if err != nil {
			fqmn := p.fqmn
			panic(fmt.Errorf("error parsing hidden template for parameter [%s:%s]: %w", fqmn, metadata.GetID(), err))
		}
		result := &strings.Builder{}
		err = template.Execute(result, templateSource)
		if err != nil {
			fqmn := p.fqmn
			panic(fmt.Errorf("error executing hidden template for parameter [%s:%s]: %w", fqmn, metadata.GetID(), err))
		}

		hiddenResult, err := strconv.ParseBool(result.String())
		if err != nil {
			fqmn := p.fqmn
			panic(fmt.Errorf("error parsing hidden template result for parameter [%s:%s]: %w", fqmn, metadata.GetID(), err))
		}
		if !hiddenResult {
			params = append(params, param)
		}
	}
	p.layout.addFormItems(params)

	subsections := p.section.GetSections()
	for i, section := range subsections {
		subPage := p.subPages[i]

		hidden := section.GetHidden()

		// Handle sections that don't have a hidden template
		if hidden.Template == "" {
			if !hidden.Default {
				addSubsectionButton(section, subPage, md, p.layout.form)
			}
			continue
		}

		// Generate a template source for the section
		templateSource := configurationTemplateSource{
			fqmn:              p.fqmn,
			hdSettings:        md.newInstance,
			moduleSettingsMap: md.moduleSettingsMap,
		}

		// Update the hidden status
		template, err := template.New(string(section.GetID())).Parse(hidden.Template)
		if err != nil {
			fqmn := p.fqmn
			panic(fmt.Errorf("error parsing hidden template for section [%s:%s]: %w", fqmn, section.GetID(), err))
		}
		result := &strings.Builder{}
		err = template.Execute(result, templateSource)
		if err != nil {
			fqmn := p.fqmn
			panic(fmt.Errorf("error executing hidden template for section [%s:%s]: %w", fqmn, section.GetID(), err))
		}

		hiddenResult, err := strconv.ParseBool(result.String())
		if err != nil {
			fqmn := p.fqmn
			panic(fmt.Errorf("error parsing hidden template result for section [%s:%s]: %w", fqmn, section.GetID(), err))
		}
		if !hiddenResult {
			addSubsectionButton(section, subPage, md, p.layout.form)
		}
	}

	p.layout.refresh()
}
