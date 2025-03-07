package config

import (
	"fmt"
	"strconv"
	"strings"
	"text/template"

	"github.com/gdamore/tcell/v2"
	"github.com/nodeset-org/hyperdrive/modules"
	"github.com/nodeset-org/hyperdrive/modules/config"
	"github.com/rivo/tview"
)

// A layout container with the standard elements and design
type standardLayout struct {
	md             *MainDisplay
	grid           *tview.Grid
	content        tview.Primitive
	descriptionBox *tview.TextView
	footer         tview.Primitive
	form           *Form
	formItemMap    map[tview.FormItem]*parameterizedFormItem
	buttonMap      map[*tview.Button]*metadataButton
	redrawing      bool
}

// Creates a new StandardLayout instance, which includes the grid and description box preconstructed.
func newStandardLayout(md *MainDisplay) *standardLayout {
	// Create the main display grid
	grid := tview.NewGrid().
		SetColumns(-5, 2, -3).
		SetRows(0, 1, 0).
		SetBorders(false)

	// Create the description box
	descriptionBox := tview.NewTextView()
	descriptionBox.SetBorder(true)
	descriptionBox.SetBorderPadding(0, 0, 1, 1)
	descriptionBox.SetTitle(" Description ")
	descriptionBox.SetWordWrap(true)
	descriptionBox.SetBackgroundColor(BackgroundColor)
	descriptionBox.SetDynamicColors(true)

	grid.AddItem(descriptionBox, 0, 2, 1, 1, 0, 0, false)

	return &standardLayout{
		md:             md,
		grid:           grid,
		descriptionBox: descriptionBox,
	}
}

// Sets the main content (the box on the left side of the screen) for this layout,
// applying the default styles to it.
func (layout *standardLayout) setContent(content tview.Primitive, contentBox *tview.Box, title string) {
	// Set the standard properties for the content (border and title)
	contentBox.SetBorder(true)
	contentBox.SetBorderPadding(1, 1, 1, 1)
	contentBox.SetTitle(fmt.Sprintf(" %s ", title))

	// Add the content to the grid
	layout.content = content
	layout.grid.AddItem(content, 0, 0, 1, 1, 0, 0, true)
}

// Sets the footer for this layout.
func (layout *standardLayout) setFooter(footer tview.Primitive, height int) {
	if footer == nil {
		layout.grid.SetRows(0, 1)
	} else {
		// Add the footer to the grid
		layout.footer = footer
		layout.grid.SetRows(0, 1, height)
		layout.grid.AddItem(footer, 2, 0, 1, 3, 0, 0, false)
	}
}

// Create a standard form for this layout (for settings pages)
func (layout *standardLayout) createForm(title string) {
	// Initialize the form item and button slices and maps
	layout.formItemMap = map[tview.FormItem]*parameterizedFormItem{}
	layout.buttonMap = map[*tview.Button]*metadataButton{}

	// Create the form
	form := NewForm().
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetHorizontal(false)
	form.
		SetButtonBackgroundColor(form.fieldBackgroundColor).
		SetBackgroundColor(BackgroundColor).
		SetBorderPadding(0, 0, 0, 0)

	// Set up the selected parameter change callback to update the description box
	form.SetChangedFunc(func(index int) {
		itemCount := form.GetFormItemCount()
		buttonCount := form.GetButtonCount()
		if index < itemCount {
			formItem := form.GetFormItem(index)
			paramItem, exists := layout.formItemMap[formItem]
			if !exists {
				// Handle form items that were added out-of-band and aren't part of the config
				return
			}

			param := paramItem.parameter
			metadata := param.GetMetadata()
			defaultValue := metadata.GetDefault()

			description := metadata.GetDescription() // TEMPLATE!
			var descriptionText string
			if description.Template == "" {
				descriptionText = description.Default
			} else {
				// Generate a template source for the parameter
				templateSource := parameterTemplateSource{
					configurationTemplateSource: configurationTemplateSource{
						fqmn:              modules.HyperdriveFqmn,
						hdSettings:        layout.md.newInstance,
						moduleSettingsMap: layout.md.moduleSettingsMap,
					},
					parameter: metadata,
				}

				// Execute the description template
				template, err := template.New(string(metadata.GetID())).Parse(description.Template)
				if err != nil {
					fqmn := modules.HyperdriveFqmn
					panic(fmt.Errorf("error parsing description template for parameter [%s:%s]: %w", fqmn, metadata.GetID(), err))
				}
				result := &strings.Builder{}
				err = template.Execute(result, templateSource)
				if err != nil {
					fqmn := modules.HyperdriveFqmn
					panic(fmt.Errorf("error executing description template for parameter [%s:%s]: %w", fqmn, metadata.GetID(), err))
				}
				descriptionText = result.String()
			}

			descriptionText = fmt.Sprintf("Default: %v\n\n%s", defaultValue, descriptionText)
			layout.descriptionBox.SetText(descriptionText)
			layout.descriptionBox.ScrollToBeginning()
		} else if index < itemCount+buttonCount {
			// This is a button
			button := form.GetButton(index - itemCount)
			metadataButton, exists := layout.buttonMap[button]
			if !exists {
				// Handle buttons that were added out-of-band and aren't part of the config
				return
			}

			section := metadataButton.section
			description := section.GetDescription()
			var descriptionText string
			if description.Template == "" {
				descriptionText = description.Default
			} else {
				// Generate a template source for the parameter
				templateSource := configurationTemplateSource{
					fqmn:              modules.HyperdriveFqmn,
					hdSettings:        layout.md.newInstance,
					moduleSettingsMap: layout.md.moduleSettingsMap,
				}

				// Execute the description template
				id := string(section.GetID())
				template, err := template.New(id).Parse(description.Template)
				if err != nil {
					fqmn := modules.HyperdriveFqmn
					panic(fmt.Errorf("error parsing description template for section [%s:%s]: %w", fqmn, id, err))
				}
				result := &strings.Builder{}
				err = template.Execute(result, templateSource)
				if err != nil {
					fqmn := modules.HyperdriveFqmn
					panic(fmt.Errorf("error executing description template for section [%s:%s]: %w", fqmn, id, err))
				}
				descriptionText = result.String()
			}
			layout.descriptionBox.SetText(descriptionText)
		}
	})

	layout.form = form
	layout.setContent(form, form.Box, title)
	layout.createSettingFooter()
}

// Refreshes all of the form items to show the current configured values
func (layout *standardLayout) refresh() {
	for i := 0; i < layout.form.GetFormItemCount(); i++ {
		formItem := layout.form.GetFormItem(i)
		paramItem, exists := layout.formItemMap[formItem]
		if !exists {
			// Handle form items that were added out-of-band and aren't part of the config
			continue
		}
		param := paramItem.parameter
		metadata := param.GetMetadata()

		// Set the form item to the current value
		if metadata.GetType() == config.ParameterType_Bool {
			// Bool
			formItem.(*tview.Checkbox).SetChecked(param.GetValue().(bool))
		} else if choiceParam, ok := metadata.(config.IChoiceParameter); ok {
			// Choice
			for i, option := range choiceParam.GetOptions() {
				if option.GetValue() == param.String() {
					formItem.(*DropDown).SetCurrentOption(i)
				}
			}
		} else {
			// Everything else
			inputField, ok := formItem.(*tview.InputField)
			if !ok {
				panic(fmt.Errorf("form item [%s] is not an input field, it's a %T", metadata.GetID(), formItem))
			}
			inputField.SetText(param.String())
		}
	}

	// Focus the first element
	layout.form.SetFocus(0)
}

// Create the footer, including the nav bar
func (layout *standardLayout) createSettingFooter() {
	// Nav bar
	navString1 := "Arrow keys: Navigate   Space/Enter: Change Setting"
	navTextView1 := tview.NewTextView().
		SetDynamicColors(false).
		SetRegions(false).
		SetWrap(false)
	fmt.Fprint(navTextView1, navString1)

	navString2 := "Esc: Return to Previous Page"
	navTextView2 := tview.NewTextView().
		SetDynamicColors(false).
		SetRegions(false).
		SetWrap(false)
	fmt.Fprint(navTextView2, navString2)

	navBar := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().
			AddItem(tview.NewBox(), 0, 1, false).
			AddItem(navTextView1, len(navString1), 1, false).
			AddItem(tview.NewBox(), 0, 1, false),
			1, 1, false).
		AddItem(tview.NewFlex().
			AddItem(tview.NewBox(), 0, 1, false).
			AddItem(navTextView2, len(navString2), 1, false).
			AddItem(tview.NewBox(), 0, 1, false),
			1, 1, false)

	layout.setFooter(navBar, 2)
}

// Register a collection of form items with this layout
func (layout *standardLayout) registerFormItems(params ...*parameterizedFormItem) {
	for _, param := range params {
		layout.formItemMap[param.item] = param
	}
}

// Register a collection of buttons with this layout
func (layout *standardLayout) registerButtons(buttons ...*metadataButton) {
	for _, button := range buttons {
		layout.buttonMap[button.button] = button
	}
}

// Sets up a handler to return to the specified homePage when the user presses escape on the layout.
func (layout *standardLayout) setupEscapeReturnHomeHandler(md *MainDisplay, homePage *page) {
	layout.grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Return to the home page
		if event.Key() == tcell.KeyEsc {
			// Close all dropdowns and break if one was open
			for _, param := range layout.formItemMap {
				dropDown, ok := param.item.(*DropDown)
				if ok && dropDown.open {
					dropDown.CloseList(md.app)
					return nil
				}
			}
			md.setPage(homePage)
			return nil
		}
		return event
	})
}

// Redraw the layout's form
// formInit: A callback that runs after the form has been cleared, before the items are added back in. Return false if the redraw should end after this step.
func (layout *standardLayout) redrawForm(
	fqmn string,
	formItems []*parameterizedFormItem,
	buttons []*metadataButton,
	formInit func() bool,
) {
	// Prevent re-entry if we're already redrawing
	if layout.redrawing {
		return
	}
	layout.redrawing = true
	defer func() {
		layout.redrawing = false
	}()

	// Get the item that's currently selected, if there is one
	var itemToFocus tview.FormItem = nil
	focusedItemIndex, focusedButtonIndex := layout.form.GetFocusedItemIndex()
	if focusedItemIndex != -1 {
		focusedItem := layout.form.GetFormItem(focusedItemIndex)
		for _, pfi := range formItems {
			if pfi.item == focusedItem {
				itemToFocus = focusedItem
				break
			}
		}
	}

	// Get the button that's currently selected, if there is one
	var buttonToFocus *tview.Button = nil
	if focusedButtonIndex != -1 {
		item := layout.form.GetButton(focusedButtonIndex)
		for _, button := range buttons {
			if button.button == item {
				buttonToFocus = item
				break
			}
		}
	}

	// Clear the form and run the initializer callback
	layout.form.Clear(true)
	if formInit != nil {
		shouldStop := formInit()
		if shouldStop {
			layout.refresh()
			return
		}
	}

	// Add the parameter form items back in
	md := layout.md
	visibleParams := []*parameterizedFormItem{}
	for _, pfi := range formItems {
		metadata := pfi.parameter.GetMetadata()
		hidden := metadata.GetHidden()

		// Handle parameters that don't have a hidden template
		if hidden.Template == "" {
			if !hidden.Default {
				visibleParams = append(visibleParams, pfi)
			}
			continue
		}

		// Generate a template source for the parameter
		templateSource := parameterTemplateSource{
			configurationTemplateSource: configurationTemplateSource{
				fqmn:              fqmn,
				hdSettings:        md.newInstance,
				moduleSettingsMap: md.moduleSettingsMap,
			},
			parameter: metadata,
		}

		// Execute the hidden template
		id := string(metadata.GetID())
		template, err := template.New(id).Parse(hidden.Template)
		if err != nil {
			panic(fmt.Errorf("error parsing hidden template for parameter [%s:%s]: %w", fqmn, id, err))
		}
		result := &strings.Builder{}
		err = template.Execute(result, templateSource)
		if err != nil {
			panic(fmt.Errorf("error executing hidden template for parameter [%s:%s]: %w", fqmn, id, err))
		}

		hiddenResult, err := strconv.ParseBool(result.String())
		if err != nil {
			panic(fmt.Errorf("error parsing hidden template result for parameter [%s:%s]: %w", fqmn, id, err))
		}
		if !hiddenResult {
			visibleParams = append(visibleParams, pfi)
		}
	}
	for _, pfi := range visibleParams {
		layout.form.AddFormItem(pfi.item)
	}

	// Add the subsection buttons back in
	visibleButtons := []*metadataButton{}
	for _, button := range buttons {
		section := button.section
		hidden := section.GetHidden()

		// Handle sections that don't have a hidden template
		if hidden.Template == "" {
			if !hidden.Default {
				visibleButtons = append(visibleButtons, button)
			}
			continue
		}

		// Generate a template source for the section
		templateSource := configurationTemplateSource{
			fqmn:              fqmn,
			hdSettings:        md.newInstance,
			moduleSettingsMap: md.moduleSettingsMap,
		}

		// Update the hidden status
		template, err := template.New(string(section.GetID())).Parse(hidden.Template)
		if err != nil {
			panic(fmt.Errorf("error parsing hidden template for section [%s:%s]: %w", fqmn, section.GetID(), err))
		}
		result := &strings.Builder{}
		err = template.Execute(result, templateSource)
		if err != nil {
			panic(fmt.Errorf("error executing hidden template for section [%s:%s]: %w", fqmn, section.GetID(), err))
		}

		hiddenResult, err := strconv.ParseBool(result.String())
		if err != nil {
			panic(fmt.Errorf("error parsing hidden template result for section [%s:%s]: %w", fqmn, section.GetID(), err))
		}
		if !hiddenResult {
			visibleButtons = append(visibleButtons, button)
		}
	}
	for _, button := range visibleButtons {
		layout.form.AddButton(button.button)
	}

	// Redraw the layout
	layout.refresh()

	// Reselect the item that was previously selected if possible, otherwise focus the enable button
	if itemToFocus != nil {
		for _, param := range visibleParams {
			if param.item != itemToFocus {
				continue
			}

			label := param.parameter.GetMetadata().GetName()
			index := layout.form.GetFormItemIndex(label)
			if index != -1 {
				layout.form.SetFocus(index)
			}
			break
		}
	} else if buttonToFocus != nil {
		for _, button := range visibleButtons {
			if button.button != buttonToFocus {
				continue
			}

			label := button.section.GetName()
			index := layout.form.GetButtonIndex(label)
			if index != -1 {
				layout.form.SetFocus(len(visibleParams) + index)
			}
			break
		}
	} else {
		layout.form.SetFocus(0)
	}
}
