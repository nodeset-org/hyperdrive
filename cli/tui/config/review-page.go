package config

import (
	"context"
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	modconfig "github.com/nodeset-org/hyperdrive/modules/config"
	"github.com/rivo/tview"
)

// Constants
const reviewPageID string = "review-settings"

type changedSetting struct {
	name     string
	oldValue string
	newValue string
}

type changedSection struct {
	header  string
	changes []changedSetting
}

// The changed settings review page
type ReviewPage struct {
	md      *MainDisplay
	changes []*changedSection
	page    *page
}

// Create a page to review any changes
func NewReviewPage(md *MainDisplay) *ReviewPage {
	//var changes []*changedSection
	//var totalAffectedContainers map[config.ContainerID]bool
	var containersToRestart []string
	var message string

	// Create the visual list for all of the changed settings
	changeBox := tview.NewTextView().
		SetDynamicColors(true).
		SetWordWrap(true)
	changeBox.SetBorder(true)
	changeBox.SetBackgroundColor(BackgroundColor)
	changeBox.SetBorderPadding(0, 0, 1, 1)

	// Update the base settings
	err := md.newInstance.ConvertToKnownType(md.NewSettings)
	if err != nil {
		panic(fmt.Errorf("error converting updated base settings to known type: %w", err))
	}

	// Update the module settings
	for _, modulePage := range md.settingsHome.modulesPage.moduleSubpages {
		modInstance := modulePage.instance
		modInstance.SetSettings(modulePage.settings)
	}

	// Validate the new configuration locally
	fullErrorStringBuilder := strings.Builder{}
	localValidationResult := modconfig.Validate(md.Config, md.newInstance)
	errorString := processSectionValidationResult("Hyperdrive", localValidationResult.ContainerResult)
	if errorString != "" {
		fullErrorStringBuilder.WriteString(errorString)
	}
	for _, modulePage := range md.settingsHome.modulesPage.moduleSubpages {
		if !modulePage.instance.Enabled {
			continue
		}
		modValidationResult := modconfig.Validate(modulePage.info.Configuration, modulePage.settings)
		moduleErrorString := processModuleValidationResult(modulePage.info, modValidationResult)
		if moduleErrorString != "" {
			fullErrorStringBuilder.WriteString(moduleErrorString)
		}
	}

	// If there aren't any local issues, process the settings with the modules
	modulePrettyNames := map[string]string{}
	openPortMap := map[string]map[string]uint16{}
	if fullErrorStringBuilder.Len() == 0 {
		for _, modulePage := range md.settingsHome.modulesPage.moduleSubpages {
			// Check if it's enabled
			if !modulePage.instance.Enabled {
				continue
			}

			// Process the settings
			fqmn := modulePage.info.Descriptor.GetFullyQualifiedModuleName()
			md.NewSettings.Modules[fqmn].Restart = nil // Clear the service restart cache
			modulePrettyNames[fqmn] = fmt.Sprintf("%s / %s", modulePage.info.Descriptor.Author, modulePage.info.Descriptor.Name)
			gac, err := md.moduleManager.GetGlobalAdapterClient(fqmn)
			if err != nil {
				fullErrorStringBuilder.WriteString(fmt.Sprintf("Module %s:\nGlobal adapter client error: %s\n\n", modulePrettyNames[fqmn], err))
				continue
			}
			response, err := gac.ProcessSettings(context.Background(), md.PreviousSettings, md.NewSettings)
			if err != nil {
				fullErrorStringBuilder.WriteString(fmt.Sprintf("Module %s:\nSettings validation error: %s\n\n", modulePrettyNames[fqmn], err))
				continue
			}
			if len(response.Errors) > 0 {
				fullErrorStringBuilder.WriteString(fmt.Sprintf("Module %s:\n", modulePrettyNames[fqmn]))
				for _, err := range response.Errors {
					fullErrorStringBuilder.WriteString(err)
					fullErrorStringBuilder.WriteString("\n")
				}
				fullErrorStringBuilder.WriteString("\n")
				continue
			}

			// Map the open ports
			openPortMap[fqmn] = response.Ports

			// Add the services / containers to restart
			for _, service := range response.ServicesToRestart {
				containerName := md.PreviousSettings.ProjectName + "-" + string(modulePage.info.Descriptor.Shortcut) + "_" + service
				containersToRestart = append(containersToRestart, containerName)
				md.NewSettings.Modules[fqmn].Restart = append(md.NewSettings.Modules[fqmn].Restart, service)
			}
		}
	}

	if fullErrorStringBuilder.Len() > 0 {
		// Create the error message if there are any issues
		message = "[orange]WARNING: Your configuration encountered errors. You must correct the following in order to save it:\n\n"
		message += fullErrorStringBuilder.String()
		changeBox.SetText(message)
	} else {
		// Print the config changes
		fullChangesBuilder := strings.Builder{}
		diff := modconfig.CompareSettings(md.Config, md.previousInstance, md.newInstance)
		diffString := processSectionDifference("Hyperdrive", diff.ContainerDifference)
		if diffString != "" {
			fullChangesBuilder.WriteString(diffString)
		}
		for _, modulePage := range md.settingsHome.modulesPage.moduleSubpages {
			if !modulePage.instance.Enabled {
				if modulePage.previousSettings == nil {
					continue
				}
				if modulePage.previousInstance.Enabled {
					fullChangesBuilder.WriteString(fmt.Sprintf("Module %s: Enabled -> Disabled\n", modulePage.info.Descriptor.Name))
					continue
				}
			} else {
				if modulePage.previousSettings == nil {
					fullChangesBuilder.WriteString(fmt.Sprintf("Module %s: New module\n", modulePage.info.Descriptor.Name))
					continue
				}
				if !modulePage.previousInstance.Enabled {
					fullChangesBuilder.WriteString(fmt.Sprintf("Module %s: Disabled -> Enabled\n", modulePage.info.Descriptor.Name))
					continue
				}
			}
			modDiff := modconfig.CompareSettings(modulePage.info.Configuration, modulePage.previousSettings, modulePage.settings)
			moduleDiffString := processModuleDifferences(modulePage.info, modDiff)
			if moduleDiffString != "" {
				fullChangesBuilder.WriteString(moduleDiffString)
			}
		}

		// Print the list of containers to restart
		if fullChangesBuilder.String() == "" {
			fullChangesBuilder.WriteString("<No changes>")
		} else if len(containersToRestart) > 0 {
			fullChangesBuilder.WriteString("The following containers must be restarted for these changes to take effect:")
			for _, container := range containersToRestart {
				//containerName := oldConfig.Hyperdrive.GetDockerArtifactName(string(container))
				//builder.WriteString(fmt.Sprintf("\n\t%s", containerName))
				//containersToRestart = append(containersToRestart, container)
				fullChangesBuilder.WriteString("\n\t" + container)
			}
		}

		changeBox.SetText(fullChangesBuilder.String())
	}

	// Create the layout
	width := 86

	// Create the main text view
	descriptionText := "Please review your changes below.\nScroll through them using the arrow keys, and press Enter when you're ready to save them."
	lines := tview.WordWrap(descriptionText, width-4)
	textViewHeight := len(lines) + 1
	textView := tview.NewTextView().
		SetText(descriptionText).
		SetTextAlign(tview.AlignCenter).
		SetWordWrap(true).
		SetTextColor(tview.Styles.PrimaryTextColor)
	textView.SetBackgroundColor(BackgroundColor)
	textView.SetBorderPadding(0, 0, 1, 1)

	var buttonGrid *tview.Flex

	if fullErrorStringBuilder.Len() > 0 {
		buttonGrid = tview.NewFlex().
			SetDirection(tview.FlexColumn).
			AddItem(tview.NewBox().
				SetBackgroundColor(BackgroundColor), 0, 1, false)
	} else {
		// Create the save button
		saveButton := tview.NewButton("Save Settings")
		saveButton.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyUp || event.Key() == tcell.KeyDown {
				changeBox.InputHandler()(event, nil)
				return nil
			}
			return event
		})
		// Save when selected
		saveButton.SetSelectedFunc(func() {
			md.ShouldSave = true
			//if changeNetworks && !md.isNew {
			//	md.ChangeNetworks = true
			//}
			md.app.Stop()
		})
		saveButton.SetBackgroundColorActivated(tcell.Color46)
		saveButton.SetLabelColorActivated(tcell.ColorBlack)

		buttonGrid = tview.NewFlex().
			SetDirection(tview.FlexColumn).
			AddItem(tview.NewBox().
				SetBackgroundColor(BackgroundColor), 0, 1, false).
			AddItem(saveButton, len(saveButton.GetLabel())+2, 0, true).
			AddItem(tview.NewBox().
				SetBackgroundColor(BackgroundColor), 0, 1, false)
	}

	// Row spacers with the correct background color
	spacer1 := tview.NewBox().
		SetBackgroundColor(BackgroundColor)
	spacer2 := tview.NewBox().
		SetBackgroundColor(BackgroundColor)
	spacer3 := tview.NewBox().
		SetBackgroundColor(BackgroundColor)
	spacer4 := tview.NewBox().
		SetBackgroundColor(BackgroundColor)
	spacerL := tview.NewBox().
		SetBackgroundColor(BackgroundColor)
	spacerR := tview.NewBox().
		SetBackgroundColor(BackgroundColor)

	// The main content grid
	contentGrid := tview.NewGrid().
		SetRows(1, textViewHeight, 1, 0, 1, 1, 1).
		SetColumns(1, 0, 1).
		AddItem(spacer1, 0, 1, 1, 1, 0, 0, false).
		AddItem(textView, 1, 1, 1, 1, 0, 0, false).
		AddItem(spacer2, 2, 1, 1, 1, 0, 0, false).
		AddItem(changeBox, 3, 1, 1, 1, 0, 0, false).
		AddItem(spacer3, 4, 1, 1, 1, 0, 0, false).
		AddItem(buttonGrid, 5, 1, 1, 1, 0, 0, true).
		AddItem(spacer4, 6, 1, 1, 1, 0, 0, false).
		AddItem(spacerL, 0, 0, 7, 1, 0, 0, false).
		AddItem(spacerR, 0, 2, 7, 1, 0, 0, false)
	contentGrid.
		SetBackgroundColor(BackgroundColor).
		SetBorder(true).
		SetTitle(" Review Changes ")
	contentGrid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			md.setPage(md.settingsHome.homePage)
			return nil
		default:
			return event
		}
	})

	// A grid with variable spaced borders that surrounds the fixed-size content grid
	borderGrid := tview.NewGrid().
		SetColumns(0, width, 0)
	borderGrid.AddItem(contentGrid, 1, 1, 1, 1, 0, 0, true)

	// Get the total content height, including spacers and borders
	borderGrid.SetRows(1, 0, 1, 1, 1)

	// Create the nav footer text view
	navString1 := "Arrow keys: Navigate     Space/Enter: Select"
	navTextView1 := tview.NewTextView().
		SetDynamicColors(false).
		SetRegions(false).
		SetWrap(false)
	fmt.Fprint(navTextView1, navString1)

	navString2 := "Esc: Go Back     Ctrl+C: Quit without Saving"
	navTextView2 := tview.NewTextView().
		SetDynamicColors(false).
		SetRegions(false).
		SetWrap(false)
	fmt.Fprint(navTextView2, navString2)

	// Create the nav footer
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
	borderGrid.AddItem(navBar, 3, 1, 1, 1, 0, 0, true)

	page := newPage(nil, reviewPageID, "Review Settings", "", borderGrid)

	return &ReviewPage{
		md: md,
		//changedSettings: changes,
		page: page,
	}
}

// Process the validation result for a module, printing any errors to a formatted string
func processModuleValidationResult(info *modconfig.ModuleInfo, result *modconfig.ModuleValidationResult) string {
	return processSectionValidationResult(string(info.Descriptor.Name), result.ContainerResult)
}

// Process the validation result for a module section, printing any errors to a formatted string
func processSectionValidationResult(header string, result *modconfig.ContainerResult) string {
	builder := strings.Builder{}

	// Check each param
	for _, paramResult := range result.ParameterResults {
		if len(paramResult.Errors) == 0 {
			continue
		}

		lead := header
		if header != "" {
			lead += "/"
		}
		lead += paramResult.Parameter.GetName() + ": "
		for _, err := range paramResult.Errors {
			builder.WriteString(lead + err.Error() + "\n\n")
		}
	}

	// Check each section
	for _, sectionResult := range result.SectionResults {
		lead := header
		if header != "" {
			lead += "/"
		}
		lead += sectionResult.Section.GetName()
		if sectionResult.Error != nil {
			builder.WriteString(lead + ": " + sectionResult.Error.Error() + "\n\n")
		} else {
			builder.WriteString(processSectionValidationResult(lead, sectionResult.ContainerResult))
		}
	}

	return builder.String()
}

// Process the difference between two modules, printing any changes to a formatted string
func processModuleDifferences(info *modconfig.ModuleInfo, diff *modconfig.ModuleDifference) string {
	return processSectionDifference(string(info.Descriptor.Name), diff.ContainerDifference)
}

// Process the difference between two sections, printing any changes to a formatted string
func processSectionDifference(header string, result *modconfig.ContainerDifference) string {
	builder := strings.Builder{}

	// Check each param
	for _, paramDiff := range result.ParameterDifferences {
		lead := header
		if header != "" {
			lead += "/"
		}
		lead += paramDiff.Parameter.GetName() + ": "
		builder.WriteString(lead + paramDiff.OldValue.String() + " -> " + paramDiff.NewValue.String() + "\n")
	}

	// Check each section
	for _, sectionDiff := range result.SectionDifferences {
		lead := header
		if header != "" {
			lead += "/"
		}
		lead += sectionDiff.Section.GetName()
		builder.WriteString(processSectionDifference(lead, sectionDiff.ContainerDifference))
	}

	return builder.String()
}
