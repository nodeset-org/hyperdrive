package config

// Interface for elements that expose an identifier
type IIdentifiable interface {
	GetID() Identifier
}

// Interface for elements that have a dynamic description
type IDescribable interface {
	GetDescription() DynamicProperty[string]
}

// Interface for elements that have a dynamic hidden flag
type IHideable interface {
	GetHidden() DynamicProperty[bool]
}

// Interface for elements that have a dynamic disabled flag
type IDisableable interface {
	GetDisabled() DynamicProperty[bool]
}
