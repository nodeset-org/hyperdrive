package config

import (
	"fmt"
	"strconv"
	"strings"
	"text/template"

	"github.com/nodeset-org/hyperdrive/modules/config"
	modconfig "github.com/nodeset-org/hyperdrive/modules/config"
	"github.com/rivo/tview"
)

type iSectionPage interface {
	getMainDisplay() *MainDisplay
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
	redrawing bool

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
	md := p.getMainDisplay()
	p.layout = newStandardLayout(p.getMainDisplay())
	p.layout.createForm(string(p.section.GetName()))
	p.layout.setupEscapeReturnHomeHandler(p.md, p.parent.getPage())

	// Create the section page
	p.page = newPage(
		p.parent.getPage(),
		p.parent.getPage().id+"/"+string(p.section.GetID()),
		string(p.section.GetName()),
		"", // string(section.GetDescription().Default), // TEMPLATE
		p.layout.grid,
	)

	// Set up the form items
	sectionCfg := p.section
	params := sectionCfg.GetParameters()
	for _, param := range params {
		id := param.GetID()
		paramSetting, err := p.settings.GetParameter(id)
		if err != nil {
			panic(fmt.Errorf("error getting [%s] parameter setting [%s]: %w", p.section.GetName(), id, err)) // TODO: better logging, like FQMN
		}

		// Create the form item for the parameter
		pfi := createParameterizedFormItem(paramSetting, p.layout.descriptionBox, p.handleLayoutChanged)
		p.layout.mapParameterizedFormItems(pfi)
		p.formItems = append(p.formItems, pfi)
	}

	// Set up the section subpages
	subsections := sectionCfg.GetSections()
	for _, section := range subsections {
		id := section.GetID()
		settingsSection, err := p.settings.GetSection(id)
		if err != nil {
			panic(fmt.Errorf("error getting [%s] section setting [%s]: %w", p.section.GetName(), id, err)) // TODO: better logging, like FQMN
		}

		// Create the subpage
		subPage := NewSectionPage(p.md, p, section, settingsSection, p.fqmn)
		p.subPages = append(p.subPages, subPage)
		md.pages.AddPage(subPage.getPage().id, subPage.getPage().content, true, false)

		// Create the metadata button for the section
		button := createMetadataButton(section, subPage, md)
		p.layout.mapMetadataButton(button)
		p.buttons = append(p.buttons, button)
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
	// Prevent re-entry if we're already redrawing
	if p.redrawing {
		return
	}
	p.redrawing = true
	defer func() {
		p.redrawing = false
	}()

	// Get the item that's currently selected, if there is one
	var itemToFocus tview.FormItem = nil
	focusedItemIndex, focusedButtonIndex := p.layout.form.GetFocusedItemIndex()
	if focusedItemIndex != -1 {
		focusedItem := p.layout.form.GetFormItem(focusedItemIndex)
		for _, pfi := range p.formItems {
			if pfi.item == focusedItem {
				itemToFocus = focusedItem
				break
			}
		}
	}

	// Get the button that's currently selected, if there is one
	var buttonToFocus *tview.Button = nil
	if focusedButtonIndex != -1 {
		item := p.layout.form.GetButton(focusedButtonIndex)
		for _, button := range p.buttons {
			if button.button == item {
				buttonToFocus = item
				break
			}
		}
	}

	// Clear the form
	p.layout.form.Clear(true)

	// Add the parameter form items back in
	md := p.getMainDisplay()
	params := []*parameterizedFormItem{}
	for _, pfi := range p.formItems {
		metadata := pfi.parameter.GetMetadata()
		hidden := metadata.GetHidden()

		// Handle parameters that don't have a hidden template
		if hidden.Template == "" {
			if !hidden.Default {
				params = append(params, pfi)
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
			parameter: pfi.parameter.GetMetadata(),
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
			params = append(params, pfi)
		}
	}
	p.layout.addFormItems(params)

	// Add the subsection buttons back in
	buttons := []*metadataButton{}
	for _, button := range p.buttons {
		section := button.section
		hidden := section.GetHidden()

		// Handle sections that don't have a hidden template
		if hidden.Template == "" {
			if !hidden.Default {
				buttons = append(buttons, button)
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
			buttons = append(buttons, button)
		}
	}
	p.layout.addButtons(buttons)

	// Redraw the layout
	p.layout.refresh()

	// Reselect the item that was previously selected if possible, otherwise focus the enable button
	if itemToFocus != nil {
		for _, param := range params {
			if param.item != itemToFocus {
				continue
			}

			label := param.parameter.GetMetadata().GetName()
			index := p.layout.form.GetFormItemIndex(label)
			if index != -1 {
				p.layout.form.SetFocus(index)
			}
			break
		}
	} else if buttonToFocus != nil {
		for _, button := range buttons {
			if button.button != buttonToFocus {
				continue
			}

			label := button.section.GetName()
			index := p.layout.form.GetButtonIndex(label)
			if index != -1 {
				p.layout.form.SetFocus(len(params) + index)
			}
			break
		}
	} else {
		p.layout.form.SetFocus(0)
	}
}
