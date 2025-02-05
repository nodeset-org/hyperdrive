package config

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/nodeset-org/hyperdrive/modules/config"
	modconfig "github.com/nodeset-org/hyperdrive/modules/config"
	"github.com/rivo/tview"
)

type ModulePage struct {
	modulesPage *ModulesPage
	page        *page
	layout      *standardLayout
	info        *config.ModuleInfo
	instance    *modconfig.ModuleInstance
	settings    *modconfig.ModuleSettings
	subPages    []iSectionPage

	// Fields
	enableBox *parameterizedFormItem
	params    []*parameterizedFormItem
}

// Create a new module page
func NewModulePage(modulesPage *ModulesPage, info *config.ModuleInfo, instance *modconfig.ModuleInstance) *ModulePage {
	fqmn := info.Descriptor.GetFullyQualifiedModuleName()
	settings := modconfig.CreateModuleSettings(info.Configuration)
	err := settings.DeserializeFromMap(instance.Settings)
	if err != nil {
		panic(fmt.Errorf("error deserializing module [%s] settings: %w", fqmn, err))
	}

	moduleDescString := strings.Builder{}
	moduleDescString.WriteString("Name:\n" + string(info.Descriptor.Name) + "\n\n")
	moduleDescString.WriteString("Author:\n" + string(info.Descriptor.Author) + "\n\n")
	moduleDescString.WriteString("Version:\n" + info.Descriptor.Version.String() + "\n\n")
	moduleDescString.WriteString("Description:\n" + string(info.Descriptor.Description))

	modulePage := &ModulePage{
		modulesPage: modulesPage,
		info:        info,
		instance:    instance,
		settings:    settings,
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
	return modulePage
}

// Creates the content for the Constellation settings page
func (configPage *ModulePage) createContent() {
	// Create the layout
	configPage.layout = newStandardLayout()
	configPage.layout.createForm(string(configPage.info.Descriptor.Name) + " Settings")
	configPage.layout.setupEscapeReturnHomeHandler(configPage.modulesPage.home.md, configPage.modulesPage.page)
	configPage.layout.form.SetButtonBackgroundColor(configPage.layout.form.fieldBackgroundColor)

	// Set up the enable box
	enableParam := NewEnableParamInstance(configPage.info, configPage.instance)
	configPage.enableBox = createParameterizedCheckbox(enableParam)

	// Set up the form items
	moduleCfg := configPage.info.Configuration
	params := moduleCfg.GetParameters()
	settings := []modconfig.IParameterSetting{}
	for _, param := range params {
		id := param.GetID()
		setting, err := configPage.settings.GetParameter(id)
		if err != nil {
			panic(fmt.Errorf("error getting logging parameter setting [%s]: %w", id, err))
		}
		settings = append(settings, setting)
	}
	configPage.params = createParameterizedFormItems(settings, configPage.layout.descriptionBox)

	// Map the parameters to the form items in the layout
	configPage.layout.mapParameterizedFormItems(configPage.enableBox)
	configPage.layout.mapParameterizedFormItems(configPage.params...)

	// Set up the setting callbacks
	configPage.enableBox.item.(*tview.Checkbox).SetChangedFunc(func(checked bool) {
		if configPage.instance.Enabled == checked {
			return
		}
		configPage.instance.Enabled = checked
		configPage.handleLayoutChanged()
	})

	// Do the initial draw
	configPage.handleLayoutChanged()
}

func (p *ModulePage) setupSubpages() {
	moduleCfg := p.info.Configuration
	subsections := moduleCfg.GetSections()
	for _, section := range subsections {
		id := section.GetID()
		setting, err := p.settings.GetSection(id)
		if err != nil {
			fqmn := p.info.Descriptor.GetFullyQualifiedModuleName()
			panic(fmt.Errorf("error getting [%s] section setting [%s]: %w", fqmn, id, err))
		}
		subPage := NewSectionPage(p.getMainDisplay(), p, section, setting)
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
func (p *ModulePage) getMainDisplay() *mainDisplay {
	return p.modulesPage.home.md
}

// Get the underlying page
func (p *ModulePage) getPage() *page {
	return p.page
}

// Handle a bulk redraw request
func (p *ModulePage) handleLayoutChanged() {
	p.layout.form.Clear(true)
	p.layout.form.AddFormItem(p.enableBox.item)

	if p.instance.Enabled {
		p.layout.addFormItems(p.params)

		subsections := p.info.Configuration.GetSections()
		for i, section := range subsections {
			subPage := p.subPages[i]
			p.layout.form.AddButton(section.GetName(), func() {
				subPage.handleLayoutChanged()
				p.getMainDisplay().setPage(subPage.getPage())
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
	}

	p.layout.refresh()
}
