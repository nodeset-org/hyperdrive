package config

import (
	"fmt"
	"strings"

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
	return modulePage
}

// Creates the content for the Constellation settings page
func (configPage *ModulePage) createContent() {

	// Create the layout
	configPage.layout = newStandardLayout()
	configPage.layout.createForm(string(configPage.info.Descriptor.Name) + " Settings")
	configPage.layout.setupEscapeReturnHomeHandler(configPage.modulesPage.home.md, configPage.modulesPage.page)

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
	}

	p.layout.refresh()
}
