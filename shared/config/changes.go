package config

import "github.com/nodeset-org/hyperdrive/shared/types"

// Get all of the settings that have changed between the given config sections
// Assumes the config sections represent the same element, just different instances
func getChangedSettings(old IConfigSection, new IConfigSection) (*types.ChangedSection, int) {
	changedSection := &types.ChangedSection{
		Settings:    []*types.ChangedSetting{},
		Subsections: []*types.ChangedSection{},
	}

	// Go through the parameters
	totalCount := 0
	newParams := new.GetParameters()
	for i, oldParam := range old.GetParameters() {
		newParam := newParams[i]
		oldVal := oldParam.GetValueAsString()
		newVal := newParam.GetValueAsString()
		if oldVal != newVal {
			changedSection.Settings = append(changedSection.Settings, &types.ChangedSetting{
				Name:               oldParam.GetCommon().Name,
				OldValue:           oldVal,
				NewValue:           newVal,
				AffectedContainers: oldParam.GetCommon().AffectsContainers,
			})
			totalCount++
		}
	}

	// Go through the subsections
	newSubconfigs := new.GetSubconfigs()
	for name, oldSubconfig := range old.GetSubconfigs() {
		newSubconfig := newSubconfigs[name]
		subsection, subcount := getChangedSettings(oldSubconfig, newSubconfig)
		if subcount > 0 {
			subsection.Name = oldSubconfig.GetTitle()
			changedSection.Subsections = append(changedSection.Subsections, subsection)
			totalCount += subcount
		}
	}

	return changedSection, totalCount
}

// Get a list of containers that will be need to be restarted after this change is applied
func getAffectedContainers(section *types.ChangedSection, containers map[types.ContainerID]bool) {
	for _, setting := range section.Settings {
		for _, affectedContainer := range setting.AffectedContainers {
			containers[affectedContainer] = true
		}
	}

	for _, subsection := range section.Subsections {
		getAffectedContainers(subsection, containers)
	}
}
