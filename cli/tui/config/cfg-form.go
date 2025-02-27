package config

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/nodeset-org/hyperdrive/modules/config"
	"github.com/rivo/tview"
)

// A form item linked to a Parameter
type parameterizedFormItem struct {
	parameter config.IParameterSetting
	item      tview.FormItem
}

/*
func registerEnableCheckbox(param config.Parameter[bool], checkbox *tview.Checkbox, form *Form, items []*parameterizedFormItem) {
	checkbox.SetChangedFunc(func(checked bool) {
		param.Value = checked
		if !checked {
			form.Clear(true)
			form.AddFormItem(checkbox)
		} else {
			for _, item := range items {
				form.AddFormItem(item.item)
			}
		}
	})
}
*/

// Create a list of form items based on a set of parameters
func createParameterizedFormItems(params []config.IParameterSetting, descriptionBox *tview.TextView) []*parameterizedFormItem {
	formItems := []*parameterizedFormItem{}
	for _, param := range params {
		item := getTypedFormItem(param, descriptionBox)
		formItems = append(formItems, item)
	}
	return formItems
}

// Create a form item binding for a parameter based on its type
func getTypedFormItem(param config.IParameterSetting, descriptionBox *tview.TextView) *parameterizedFormItem {
	metadata := param.GetMetadata()
	switch metadata.GetType() {
	case config.ParameterType_Choice:
		return createParameterizedDropDown(param, descriptionBox)
	case config.ParameterType_Bool:
		return createParameterizedCheckbox(param)
	case config.ParameterType_Int:
		return createParameterizedIntField(param)
	case config.ParameterType_Uint:
		return createParameterizedUintField(param)
	case config.ParameterType_Float:
		return createParameterizedFloatField(param)
	case config.ParameterType_String:
		return createParameterizedStringField(param)
	default:
		panic(fmt.Sprintf("param [%s] is not a supported type for form item binding", metadata.GetName()))
	}
}

// Create a standard form checkbox
func createParameterizedCheckbox(param config.IParameterSetting) *parameterizedFormItem {
	metadata := param.GetMetadata()
	item := tview.NewCheckbox().
		SetLabel(metadata.GetName()).
		SetChecked(param.GetValue().(bool)).
		SetChangedFunc(func(checked bool) {
			err := param.SetValue(checked)
			if err != nil {
				panic(fmt.Sprintf("error setting checkbox value for parameter [%s]: %s", metadata.GetID(), err.Error()))
			}
		})
	item.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyDown, tcell.KeyTab:
			return tcell.NewEventKey(tcell.KeyTab, 0, 0)
		case tcell.KeyUp, tcell.KeyBacktab:
			return tcell.NewEventKey(tcell.KeyBacktab, 0, 0)
		default:
			return event
		}
	})

	return &parameterizedFormItem{
		parameter: param,
		item:      item,
	}
}

// Create a standard int field
func createParameterizedIntField(param config.IParameterSetting) *parameterizedFormItem {
	metadata := param.GetMetadata()
	item := tview.NewInputField().
		SetLabel(metadata.GetName()).
		SetAcceptanceFunc(tview.InputFieldInteger)
	item.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			item.SetText("")
		} else {
			value, err := strconv.ParseInt(item.GetText(), 10, 64)
			if err != nil {
				// TODO: show error modal?
				item.SetText("")
			} else {
				err := param.SetValue(value)
				if err != nil {
					panic(fmt.Sprintf("error setting int value for parameter [%s]: %s", metadata.GetID(), err.Error()))
				}
			}
		}
	})
	item.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyDown, tcell.KeyTab:
			return tcell.NewEventKey(tcell.KeyTab, 0, 0)
		case tcell.KeyUp, tcell.KeyBacktab:
			return tcell.NewEventKey(tcell.KeyBacktab, 0, 0)
		default:
			return event
		}
	})

	return &parameterizedFormItem{
		parameter: param,
		item:      item,
	}
}

// Create a standard uint field
func createParameterizedUintField(param config.IParameterSetting) *parameterizedFormItem {
	metadata := param.GetMetadata()
	item := tview.NewInputField().
		SetLabel(metadata.GetName()).
		SetAcceptanceFunc(tview.InputFieldInteger)
	item.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			item.SetText("")
		} else {
			value, err := strconv.ParseUint(item.GetText(), 10, 64)
			if err != nil {
				// TODO: show error modal?
				item.SetText("")
			} else {
				err := param.SetValue(value)
				if err != nil {
					panic(fmt.Sprintf("error setting uint value for parameter [%s]: %s", metadata.GetID(), err.Error()))
				}
			}
		}
	})
	item.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyDown, tcell.KeyTab:
			return tcell.NewEventKey(tcell.KeyTab, 0, 0)
		case tcell.KeyUp, tcell.KeyBacktab:
			return tcell.NewEventKey(tcell.KeyBacktab, 0, 0)
		default:
			return event
		}
	})

	return &parameterizedFormItem{
		parameter: param,
		item:      item,
	}
}

// Create a standard string field
func createParameterizedStringField(param config.IParameterSetting) *parameterizedFormItem {
	metadata := param.GetMetadata().(*config.StringParameter)
	item := tview.NewInputField().
		SetLabel(metadata.GetName())
	item.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			item.SetText("")
		} else {
			err := param.SetValue(strings.TrimSpace(item.GetText()))
			if err != nil {
				panic(fmt.Sprintf("error setting string value for parameter [%s]: %s", metadata.GetID(), err.Error()))
			}
		}
	})
	item.SetAcceptanceFunc(func(textToCheck string, lastChar rune) bool {
		if metadata.MaxLength > 0 {
			if uint64(len(textToCheck)) > metadata.MaxLength {
				return false
			}
		}
		// TODO: regex support
		return true
	})
	item.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyDown, tcell.KeyTab:
			return tcell.NewEventKey(tcell.KeyTab, 0, 0)
		case tcell.KeyUp, tcell.KeyBacktab:
			return tcell.NewEventKey(tcell.KeyBacktab, 0, 0)
		default:
			return event
		}
	})

	return &parameterizedFormItem{
		parameter: param,
		item:      item,
	}
}

// Create a standard float field
func createParameterizedFloatField(param config.IParameterSetting) *parameterizedFormItem {
	metadata := param.GetMetadata().(*config.FloatParameter)
	item := tview.NewInputField().
		SetLabel(metadata.GetName()).
		SetAcceptanceFunc(tview.InputFieldFloat)
	item.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			item.SetText("")
		} else {
			value, err := strconv.ParseFloat(item.GetText(), 64)
			if err != nil {
				// TODO: show error modal?
				item.SetText("")
			} else {
				err := param.SetValue(value)
				if err != nil {
					panic(fmt.Sprintf("error setting float value for parameter [%s]: %s", metadata.GetID(), err.Error()))
				}
			}
		}
	})
	item.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyDown, tcell.KeyTab:
			return tcell.NewEventKey(tcell.KeyTab, 0, 0)
		case tcell.KeyUp, tcell.KeyBacktab:
			return tcell.NewEventKey(tcell.KeyBacktab, 0, 0)
		default:
			return event
		}
	})

	return &parameterizedFormItem{
		parameter: param,
		item:      item,
	}
}

// Create a standard choice field
func createParameterizedDropDown(param config.IParameterSetting, descriptionBox *tview.TextView) *parameterizedFormItem {
	metadata := param.GetMetadata().(config.IChoiceParameter)
	// Create the list of options
	options := []string{}
	descriptions := []string{}
	values := []any{}
	for _, option := range metadata.GetOptions() {
		options = append(options, option.GetName())
		descriptions = append(descriptions, option.GetDescription().Default) // TEMPLATE!
		values = append(values, option.GetValue())
	}
	item := NewDropDown().
		SetLabel(metadata.GetName()).
		SetOptions(options, func(text string, index int) {
			err := param.SetValue(values[index])
			if err != nil {
				panic(fmt.Sprintf("error setting choice value for parameter [%s]: %s", metadata.GetID(), err.Error()))
			}
		}).
		SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
			descriptionBox.SetText(descriptions[index])
		})
	item.SetTextOptions(" ", " ", "", "", "")
	item.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyDown, tcell.KeyTab:
			return tcell.NewEventKey(tcell.KeyTab, 0, 0)
		case tcell.KeyUp, tcell.KeyBacktab:
			return tcell.NewEventKey(tcell.KeyBacktab, 0, 0)
		default:
			return event
		}
	})
	list := item.GetList()
	list.SetSelectedBackgroundColor(tcell.Color46)
	list.SetSelectedTextColor(tcell.ColorBlack)
	list.SetBackgroundColor(tcell.ColorBlack)
	list.SetMainTextColor(tcell.ColorLightGray)

	return &parameterizedFormItem{
		parameter: param,
		item:      item,
	}
}
