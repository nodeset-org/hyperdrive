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
	formItems []*parameterizedFormItem
	buttons   []*metadataButton
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
	md := p.getMainDisplay()
	p.layout = newStandardLayout(md)
	p.layout.createForm(string(p.info.Descriptor.Name) + " Settings")
	p.layout.setupEscapeReturnHomeHandler(p.modulesPage.home.md, p.modulesPage.page)
	p.layout.form.SetButtonBackgroundColor(p.layout.form.fieldBackgroundColor)

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
	p.layout.mapParameterizedFormItems(p.enableBox)

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
		pfi := createParameterizedFormItem(paramSetting, p.layout.descriptionBox, p.handleLayoutChanged)
		p.layout.mapParameterizedFormItems(pfi)
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
		p.layout.mapMetadataButton(button)
		p.buttons = append(p.buttons, button)
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

	// Clear the form, but keep the enable box since that's always here
	p.layout.form.Clear(true)
	p.layout.form.AddFormItem(p.enableBox.item)

	// Break if the module is disabled
	if !p.instance.Enabled {
		p.layout.refresh()
		return
	}

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
				fqmn:              p.info.Descriptor.GetFullyQualifiedModuleName(),
				hdSettings:        md.newInstance,
				moduleSettingsMap: md.moduleSettingsMap,
			},
			parameter: pfi.parameter.GetMetadata(),
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
