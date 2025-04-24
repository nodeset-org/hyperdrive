package config

import (
	"strings"

	modconfig "github.com/nodeset-org/hyperdrive/modules/config"
)

// Constants
const reviewPageID string = "review-settings"

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
