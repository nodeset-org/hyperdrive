package config

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/nodeset-org/hyperdrive/modules/config"
	"github.com/rivo/tview"
)

// A layout container with the standard elements and design
type standardLayout struct {
	md             *MainDisplay
	fqmn           string
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
func newStandardLayout(md *MainDisplay, fqmn string) *standardLayout {
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
		fqmn:           fqmn,
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
		var descriptionText string
		var err error
		itemCount := form.GetFormItemCount()
		buttonCount := form.GetButtonCount()
		if index < itemCount {
			formItem := form.GetFormItem(index)
			paramItem, exists := layout.formItemMap[formItem]
			if !exists {
				// Handle form items that were added out-of-band and aren't part of the config
				return
			}

			// Get the description text
			param := paramItem.parameter.GetMetadata()
			descriptionText, err = layout.md.templateProcessor.GetEntityDescription(layout.fqmn, param)
			if err != nil {
				panic(err)
			}
			defaultValue := param.GetDefault()
			descriptionText = fmt.Sprintf("Default: %v\n\n%s", defaultValue, descriptionText)
		} else if index < itemCount+buttonCount {
			// This is a button
			button := form.GetButton(index - itemCount)
			metadataButton, exists := layout.buttonMap[button]
			if !exists {
				// Handle buttons that were added out-of-band and aren't part of the config
				return
			}

			// Get the description text
			section := metadataButton.section
			descriptionText, err = layout.md.templateProcessor.GetEntityDescription(layout.fqmn, section)
			if err != nil {
				panic(err)
			}
		} else {
			return
		}
		layout.descriptionBox.SetText(descriptionText)
		layout.descriptionBox.ScrollToBeginning()
	})

	layout.form = form
	layout.setContent(form, form.Box, title)
	layout.createSettingFooter()
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

// Redraw the layout's form, adding all of the provided form items and buttons back in.
// formInit: A callback that runs after the form has been cleared, before the items are added back in. Return false if the redraw should end after this step.
func (layout *standardLayout) redrawForm(
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

	// Add the parameter form items back in, skipping any that are hidden
	md := layout.md
	visibleParams := []*parameterizedFormItem{}
	for _, pfi := range formItems {
		isHidden, err := md.templateProcessor.IsEntityHidden(layout.fqmn, pfi.parameter.GetMetadata())
		if err != nil {
			panic(err)
		}
		if !isHidden {
			visibleParams = append(visibleParams, pfi)
		}
	}
	for _, pfi := range visibleParams {
		layout.form.AddFormItem(pfi.item)
	}

	// Add the subsection buttons back in, skipping any that are hidden
	visibleButtons := []*metadataButton{}
	for _, button := range buttons {
		hiddenResult, err := md.templateProcessor.IsEntityHidden(layout.fqmn, button.section)
		if err != nil {
			panic(err)
		}
		if !hiddenResult {
			visibleButtons = append(visibleButtons, button)
		}
	}
	for _, button := range visibleButtons {
		layout.form.AddButton(button.button)
	}

	// Refresh the values of the items to match the settings
	layout.refresh()

	// Reselect the item that was previously selected if possible
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
	}
}

// Refreshes all of the form items to show the current configured values
func (layout *standardLayout) refresh() {
	for i := range layout.form.GetFormItemCount() {
		formItem := layout.form.GetFormItem(i)
		pfi, exists := layout.formItemMap[formItem]
		if !exists {
			// Handle form items that were added out-of-band and aren't part of the config
			continue
		}
		paramSetting := pfi.parameter
		paramMetadata := paramSetting.GetMetadata()

		// Set the form item to the current value
		switch paramMetadata.GetType() {
		case config.ParameterType_Bool:
			checkbox := formItem.(*tview.Checkbox)
			value := paramSetting.GetValue().(bool)
			checkbox.SetChecked(value)

		case config.ParameterType_Choice:
			// Update the dropdown options
			choiceParam := paramMetadata.(config.IChoiceParameter)
			dropdown := formItem.(*DropDown)
			visibleValues := dropdown.ReloadDynamicOptions(paramSetting, layout)

			// Set the selected index to the current param setting
			found := false
			for i, option := range visibleValues {
				if option == paramSetting.GetValue() {
					dropdown.SetCurrentOption(i)
					found = true
					break
				}
			}
			if found {
				continue
			}

			// If the current value isn't in the visible options, set it to the default
			defaultValue := choiceParam.GetDefault()
			err := paramSetting.SetValue(defaultValue)
			if err != nil {
				panic(fmt.Errorf("error setting default value for [%s:%s]: %w", layout.fqmn, paramMetadata.GetID(), err))
			}
			for i, option := range visibleValues {
				if option == defaultValue {
					dropdown.SetCurrentOption(i)
					break
				}
			}

		default:
			// Everything else
			inputField, ok := formItem.(*tview.InputField)
			if !ok {
				panic(fmt.Errorf("form item \"%s\" is not an input field, it's a %T", paramMetadata.GetID(), formItem))
			}
			inputField.SetText(paramSetting.String())
		}
	}

	// Focus the first element
	layout.form.SetFocus(0)
}
