package types

// A configuration setting that has been changed
type ChangedSetting struct {
	// The setting name
	Name string

	// The current (old) value of the parameter
	OldValue string

	// The new (pending) value of the parameter
	NewValue string

	// A list of containers affected by this change, which will require restarting them
	AffectedContainers []ContainerID
}

// A configuration section with one or more changes
type ChangedSection struct {
	// The name of the section
	Name string

	// The list of parameters within this section that have changed
	Settings []*ChangedSetting

	// The list of subsections that may or may not have changes
	Subsections []*ChangedSection
}
