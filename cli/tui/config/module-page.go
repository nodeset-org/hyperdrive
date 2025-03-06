package config

import (
	"fmt"
	"strconv"
	"strings"

	"text/template"

	"github.com/gdamore/tcell/v2"
	"github.com/nodeset-org/hyperdrive/modules/config"
	modconfig "github.com/nodeset-org/hyperdrive/modules/config"
	"github.com/rivo/tview"
)

type ModulePage struct {
	modulesPage      *ModulesPage
	page             *page
	layout           *standardLayout
	info             *config.ModuleInfo
	fqmn             string
	previousInstance *modconfig.ModuleInstance
	instance         *modconfig.ModuleInstance
	previousSettings *modconfig.ModuleSettings
	settings         *modconfig.ModuleSettings
	subPages         []iSectionPage
	redrawing        bool

	// Fields
	enableBox *parameterizedFormItem
	params    []*parameterizedFormItem
}

// Create a new module page
func NewModulePage(modulesPage *ModulesPage, info *config.ModuleInfo, previousInstance *modconfig.ModuleInstance, instance *modconfig.ModuleInstance) *ModulePage {
	fqmn := info.Descriptor.GetFullyQualifiedModuleName()
	previousSettings := modconfig.CreateModuleSettings(info.Configuration)
	err := previousSettings.DeserializeFromMap(previousInstance.Settings)
	if err != nil {
		panic(fmt.Errorf("error deserializing previous module [%s] settings: %w", fqmn, err))
	}
	settings := modulesPage.home.md.moduleSettingsMap[fqmn]

	moduleDescString := strings.Builder{}
	moduleDescString.WriteString("Name:\n" + string(info.Descriptor.Name) + "\n\n")
	moduleDescString.WriteString("Author:\n" + string(info.Descriptor.Author) + "\n\n")
	moduleDescString.WriteString("Version:\n" + info.Descriptor.Version.String() + "\n\n")
	moduleDescString.WriteString("Description:\n" + string(info.Descriptor.Description))

	modulePage := &ModulePage{
		modulesPage:      modulesPage,
		info:             info,
		fqmn:             fqmn,
		previousInstance: previousInstance,
		instance:         instance,
		previousSettings: previousSettings,
		settings:         settings,
	}
	modulePage.createContent()

	modulePage.page = newPage(
		modulesPage.page,
		modulePrefix+fqmn,
		string(info.Descriptor.Name),
		moduleDescString.String(),
		modulePage.layout.grid,
	)
	modulePage.setupSubpages()

	// Do the initial draw
	modulePage.handleLayoutChanged()
	return modulePage
}

// Creates the content for the Constellation settings page
func (p *ModulePage) createContent() {
	// Create the layout
	p.layout = newStandardLayout(p.getMainDisplay())
	p.layout.createForm(string(p.info.Descriptor.Name) + " Settings")
	p.layout.setupEscapeReturnHomeHandler(p.modulesPage.home.md, p.modulesPage.page)
	p.layout.form.SetButtonBackgroundColor(p.layout.form.fieldBackgroundColor)

	// Set up the enable box
	enableParam := NewEnableParamInstance(p.info, p.instance)
	p.enableBox = createParameterizedCheckbox(enableParam, p.handleLayoutChanged)

	// Set up the form items
	moduleCfg := p.info.Configuration
	params := moduleCfg.GetParameters()
	settings := []modconfig.IParameterSetting{}
	for _, param := range params {
		id := param.GetID()
		setting, err := p.settings.GetParameter(id)
		if err != nil {
			panic(fmt.Errorf("error getting logging parameter setting [%s]: %w", id, err))
		}
		settings = append(settings, setting)
	}
	p.params = createParameterizedFormItems(settings, p.layout.descriptionBox, p.handleLayoutChanged)

	// Map the parameters to the form items in the layout
	p.layout.mapParameterizedFormItems(p.enableBox)
	p.layout.mapParameterizedFormItems(p.params...)

	// Set up the setting callbacks
	p.enableBox.item.(*tview.Checkbox).SetChangedFunc(func(checked bool) {
		if p.instance.Enabled == checked {
			return
		}
		p.instance.Enabled = checked
		p.handleLayoutChanged()
	})
}

func (p *ModulePage) setupSubpages() {
	moduleCfg := p.info.Configuration
	subsections := moduleCfg.GetSections()
	for _, section := range subsections {
		id := section.GetID()
		setting, err := p.settings.GetSection(id)
		if err != nil {
			panic(fmt.Errorf("error getting [%s] section setting [%s]: %w", p.fqmn, id, err))
		}
		subPage := NewSectionPage(p.getMainDisplay(), p, section, setting, p.fqmn)
		p.subPages = append(p.subPages, subPage)

		// Map the description to the section label for button shifting later
		label := section.GetName()
		p.layout.mapButtonDescription(label, section.GetDescription())
	}
	for _, subpage := range p.subPages {
		p.getMainDisplay().pages.AddPage(subpage.getPage().id, subpage.getPage().content, true, false)
	}
}

// Get the main display
func (p *ModulePage) getMainDisplay() *MainDisplay {
	return p.modulesPage.home.md
}

// Get the underlying page
func (p *ModulePage) getPage() *page {
	return p.page
}

// Handle a bulk redraw request
func (p *ModulePage) handleLayoutChanged() {
	if p.redrawing {
		return
	}
	p.redrawing = true
	defer func() {
		p.redrawing = false
	}()

	p.layout.form.Clear(true)
	p.layout.form.AddFormItem(p.enableBox.item)

	// Break if the module is disabled
	if !p.instance.Enabled {
		p.layout.refresh()
		return
	}

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
				fqmn:              p.info.Descriptor.GetFullyQualifiedModuleName(),
				hdSettings:        md.newInstance,
				moduleSettingsMap: md.moduleSettingsMap,
			},
			parameter: param.parameter.GetMetadata(),
		}

		// Update the hidden status
		template, err := template.New(string(metadata.GetID())).Parse(hidden.Template)
		if err != nil {
			fqmn := p.info.Descriptor.GetFullyQualifiedModuleName()
			panic(fmt.Errorf("error parsing hidden template for parameter [%s:%s]: %w", fqmn, metadata.GetID(), err))
		}
		result := &strings.Builder{}
		err = template.Execute(result, templateSource)
		if err != nil {
			fqmn := p.info.Descriptor.GetFullyQualifiedModuleName()
			panic(fmt.Errorf("error executing hidden template for parameter [%s:%s]: %w", fqmn, metadata.GetID(), err))
		}

		hiddenResult, err := strconv.ParseBool(result.String())
		if err != nil {
			fqmn := p.info.Descriptor.GetFullyQualifiedModuleName()
			panic(fmt.Errorf("error parsing hidden template result for parameter [%s:%s]: %w", fqmn, metadata.GetID(), err))
		}
		if !hiddenResult {
			params = append(params, param)
		}
	}
	p.layout.addFormItems(params)

	subsections := p.info.Configuration.GetSections()
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
			fqmn:              p.info.Descriptor.GetFullyQualifiedModuleName(),
			hdSettings:        md.newInstance,
			moduleSettingsMap: md.moduleSettingsMap,
		}

		// Update the hidden status
		template, err := template.New(string(section.GetID())).Parse(hidden.Template)
		if err != nil {
			fqmn := p.info.Descriptor.GetFullyQualifiedModuleName()
			panic(fmt.Errorf("error parsing hidden template for section [%s:%s]: %w", fqmn, section.GetID(), err))
		}
		result := &strings.Builder{}
		err = template.Execute(result, templateSource)
		if err != nil {
			fqmn := p.info.Descriptor.GetFullyQualifiedModuleName()
			panic(fmt.Errorf("error executing hidden template for section [%s:%s]: %w", fqmn, section.GetID(), err))
		}

		hiddenResult, err := strconv.ParseBool(result.String())
		if err != nil {
			fqmn := p.info.Descriptor.GetFullyQualifiedModuleName()
			panic(fmt.Errorf("error parsing hidden template result for section [%s:%s]: %w", fqmn, section.GetID(), err))
		}
		if !hiddenResult {
			addSubsectionButton(section, subPage, md, p.layout.form)
		}
	}

	p.layout.refresh()
}

func addSubsectionButton(section modconfig.ISection, subPage iSectionPage, md *MainDisplay, form *Form) {
	form.AddButton(section.GetName(), func() {
		subPage.handleLayoutChanged()
		md.setPage(subPage.getPage())
	})
	button := form.GetButton(form.GetButtonCount() - 1)
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
