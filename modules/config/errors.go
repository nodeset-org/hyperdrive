package config

import "fmt"

// =====================
// === ErrorNotFound ===
// =====================

type EntryType string

const (
	EntryType_Parameter EntryType = "parameter"
	EntryType_Section   EntryType = "section"
)

type ErrorNotFound struct {
	// The ID of the parameter or section that was not found
	ID Identifier

	// The type of the object that was not found
	EntryType EntryType
}

// Create a new error for a parameter not found
func NewErrorNotFound(id Identifier, entryType EntryType) ErrorNotFound {
	return ErrorNotFound{
		ID:        id,
		EntryType: entryType,
	}
}

// Get the error message for a parameter not found
func (e ErrorNotFound) Error() string {
	return fmt.Sprintf("%s [%s] not found", string(e.EntryType), string(e.ID))
}
