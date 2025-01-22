package config

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/nodeset-org/hyperdrive/modules/config"
	"github.com/rivo/tview"
)

// A layout container with the standard elements and design
type standardLayout struct {
	grid           *tview.Grid
	content        tview.Primitive
	descriptionBox *tview.TextView
	footer         tview.Primitive
	form           *Form
	parameters     map[tview.FormItem]*parameterizedFormItem
}

// Creates a new StandardLayout instance, which includes the grid and description box preconstructed.
func newStandardLayout() *standardLayout {
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
	layout.parameters = map[tview.FormItem]*parameterizedFormItem{}

	// Create the form
	form := NewForm().
		SetFieldBackgroundColor(tcell.ColorBlack)
	form.
		SetBackgroundColor(BackgroundColor).
		SetBorderPadding(0, 0, 0, 0)

	// Set up the selected parameter change callback to update the description box
	form.SetChangedFunc(func(index int) {
		if index < form.GetFormItemCount() {
			formItem := form.GetFormItem(index)
			param := layout.parameters[formItem].parameter
			metadata := param.GetMetadata()
			defaultValue := metadata.GetDefault()
			description := metadata.GetDescription().Default // TEMPLATE!
			descriptionText := fmt.Sprintf("Default: %v\n\n%s", defaultValue, description)
			layout.descriptionBox.SetText(descriptionText)
			layout.descriptionBox.ScrollToBeginning()
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
		param := layout.parameters[formItem].parameter
		metadata := param.GetMetadata()

		// Set the form item to the current value
		if _, ok := metadata.(*config.BoolParameter); ok {
			// Bool
			formItem.(*tview.Checkbox).SetChecked(param.GetValue().(bool))
		} else if choiceParam, ok := param.(config.IChoiceParameter); ok {
			// Choice
			for i, option := range choiceParam.GetOptions() {
				if option.GetValue() == param.String() {
					formItem.(*DropDown).SetCurrentOption(i)
				}
			}
		} else {
			// Everything else
			formItem.(*tview.InputField).SetText(param.String())
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

	navString2 := "Esc: Go Back to Categories"
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

// Add a collection of form items to this layout's form
func (layout *standardLayout) addFormItems(params []*parameterizedFormItem) {
	for _, param := range params {
		layout.form.AddFormItem(param.item)
	}
}

func (layout *standardLayout) mapParameterizedFormItems(params ...*parameterizedFormItem) {
	for _, param := range params {
		layout.parameters[param.item] = param
	}
}

// Sets up a handler to return to the specified homePage when the user presses escape on the layout.
func (layout *standardLayout) setupEscapeReturnHomeHandler(md *mainDisplay, homePage *page) {
	layout.grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Return to the home page
		if event.Key() == tcell.KeyEsc {
			// Close all dropdowns and break if one was open
			for _, param := range layout.parameters {
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
