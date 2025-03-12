package config

import (
	"fmt"
	"regexp"
)

const (
	// Regex pattern for the identifier
	IdentifierPattern string = `^[a-zA-Z0-9_\-\.]+$`
)

var (
	// Regex for the identifier
	IdentifierRegex = regexp.MustCompile(IdentifierPattern)
)

// Implementation of the Identifier type in the spec
type Identifier string

// String representation of the Identifier
func (i Identifier) String() string {
	return string(i)
}

// Marshal the Identifier to text
func (i Identifier) MarshalText() ([]byte, error) {
	return []byte(i), nil
}

// Unmarshal the Identifier from text
func (i *Identifier) UnmarshalText(text []byte) error {
	value := string(text)
	if !IdentifierRegex.MatchString(value) {
		return fmt.Errorf("invalid identifier: %s", value)
	}

	*i = Identifier(value)
	return nil
}

// Deserialize an identifier from a map
func deserializeIdentifier(data map[string]any, propertyName string, property *Identifier, optional bool) error {
	// Get the property in string form
	var value string
	exists, err := deserializeProperty(data, propertyName, &value, optional)
	if err != nil {
		return err
	}
	if !exists && optional {
		return nil
	}

	// Make sure the string is valid
	err = property.UnmarshalText([]byte(value))
	if err != nil {
		return fmt.Errorf("invalid %s \"%s\"", propertyName, value)
	}
	return nil
}
