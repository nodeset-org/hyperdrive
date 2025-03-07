package config

import (
	"fmt"
	"strings"

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
	formItems        []*parameterizedFormItem
	buttons          []*metadataButton

	// Fields
	enableBox *parameterizedFormItem
}

// Create a new module page
func NewModulePage(modulesPage *ModulesPage, info *config.ModuleInfo, previousInstance *modconfig.ModuleInstance, instance *modconfig.ModuleInstance) *ModulePage {
	fqmn := info.Descriptor.GetFullyQualifiedModuleName()
	settings := modulesPage.home.md.moduleSettingsMap[fqmn]
	previousSettings := modconfig.CreateModuleSettings(info.Configuration)
	err := previousSettings.DeserializeFromMap(previousInstance.Settings)
	if err != nil {
		panic(fmt.Errorf("error deserializing previous module [%s] settings: %w", fqmn, err))
	}

	// Create the page
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

	// Do the initial draw
	modulePage.handleLayoutChanged()
	return modulePage
}

// Create the content for the module page
func (p *ModulePage) createContent() {
	// Create the layout
	md := p.modulesPage.home.md
	p.layout = newStandardLayout(md, p.fqmn)
	p.layout.createForm(string(p.info.Descriptor.Name) + " Settings")
	p.layout.setupEscapeReturnHomeHandler(p.modulesPage.home.md, p.modulesPage.page)

	// Create the underlying TUI page
	moduleDescString := strings.Builder{}
	moduleDescString.WriteString("Name:\n" + string(p.info.Descriptor.Name) + "\n\n")
	moduleDescString.WriteString("Author:\n" + string(p.info.Descriptor.Author) + "\n\n")
	moduleDescString.WriteString("Version:\n" + p.info.Descriptor.Version.String() + "\n\n")
	moduleDescString.WriteString("Description:\n" + string(p.info.Descriptor.Description))
	p.page = newPage(
		p.modulesPage.page,
		modulePrefix+p.fqmn,
		string(p.info.Descriptor.Name),
		moduleDescString.String(),
		p.layout.grid,
	)

	// Set up the enable box
	enableParam := NewEnableParamInstance(p.info, p.instance)
	p.enableBox = createParameterizedCheckbox(enableParam, p.handleLayoutChanged)
	p.enableBox.item.(*tview.Checkbox).SetChangedFunc(func(checked bool) {
		if p.instance.Enabled == checked {
			return
		}
		p.instance.Enabled = checked
		p.handleLayoutChanged()
	})
	p.layout.registerFormItems(p.enableBox)

	// Set up the form items
	moduleCfg := p.info.Configuration
	params := moduleCfg.GetParameters()
	for _, param := range params {
		id := param.GetID()
		paramSetting, err := p.settings.GetParameter(id)
		if err != nil {
			panic(fmt.Errorf("error getting parameter setting [%s]: %w", id, err))
		}

		// Create the form item for the parameter
		pfi := createParameterizedFormItem(paramSetting, p.layout, p.handleLayoutChanged)
		p.layout.registerFormItems(pfi)
		p.formItems = append(p.formItems, pfi)
	}

	// Set up the section subpages
	subsections := moduleCfg.GetSections()
	for _, section := range subsections {
		id := section.GetID()
		settingsSection, err := p.settings.GetSection(id)
		if err != nil {
			panic(fmt.Errorf("error getting [%s] section setting [%s]: %w", p.fqmn, id, err))
		}

		// Create the subpage
		subPage := NewSectionPage(md, p, section, settingsSection, p.fqmn)
		p.subPages = append(p.subPages, subPage)
		md.pages.AddPage(subPage.getPage().id, subPage.getPage().content, true, false)

		// Create the metadata button for the section
		button := createMetadataButton(section, subPage, md)
		p.layout.registerButtons(button)
		p.buttons = append(p.buttons, button)
	}
}

// Get the underlying page
func (p *ModulePage) getPage() *page {
	return p.page
}

// Handle a bulk redraw request
func (p *ModulePage) handleLayoutChanged() {
	p.layout.redrawForm(p.formItems, p.buttons, func() bool {
		p.layout.form.AddFormItem(p.enableBox.item)
		return !p.instance.Enabled
	})
}

/*
// Generate the TX mode dropdown box, filtering out the options that don't have URLs on the provided network
// The TUI wasn't really designed to handle dropdows with dynamic options, so we have to do some manual work here
func (configPage *HyperdriveConfigPage) generateTxModeBox(network config.Network) {
	res := configPage.masterConfig.AllNetworkSettings.Hyperdrive[network].NetworkResources

	// Filter out the options that don't have URLs on the provided network
	origOptions := configPage.masterConfig.Hyperdrive.TxEndpointMode.Options
	options := []*config.ParameterOption[hdconfig.TxEndpointMode]{}
	for _, option := range origOptions {
		switch option.Value {
		case hdconfig.TxEndpointMode_FlashbotsProtect:
			if res.FlashbotsProtectUrl == "" {
				continue
			}
		case hdconfig.TxEndpointMode_MevBlocker:
			if res.MevBlockerUrl == "" {
				continue
			}
		}
		options = append(options, option)
	}

	// Make a clone of the param with the new options
	paramCopy := config.Parameter[hdconfig.TxEndpointMode]{
		ParameterCommon: configPage.masterConfig.Hyperdrive.TxEndpointMode.GetCommon(),
		Value:           configPage.masterConfig.Hyperdrive.TxEndpointMode.Value,
		Default:         configPage.masterConfig.Hyperdrive.TxEndpointMode.Default,
		Options:         options,
	}

	// Create the dropdown box using the cloned parameter with the filtered options
	box := createParameterizedDropDown(&paramCopy, configPage.layout.descriptionBox)
	box.item.(*DropDown).SetSelectedFunc(func(text string, index int) {
		selection := options[index].GetValueAsAny().(hdconfig.TxEndpointMode)
		if configPage.masterConfig.Hyperdrive.TxEndpointMode.Value == selection {
			return
		}
		// Update both the real parameter and the copy so it displays correctly
		configPage.masterConfig.Hyperdrive.TxEndpointMode.Value = selection
		paramCopy.Value = selection
		configPage.handleModeChanged()
	})
	configPage.layout.mapParameterizedFormItems(box)
	configPage.txModeBox = box
}
*/
