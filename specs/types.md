# Basic Types

This file defines the basic types and their validation routines that are used throughout Hyperdrive modules.


## Identifier

An **Identifier** is a string that only allows a certain subset of characters:
- Capital letters (`A-Z`)
- Lowercase letters (`a-z`)
- Digits (`0-9`)
- Underscores (`_`)
- Hyphens (`-`)
- Periods (`.`)

Identifiers can have one or more character.


## Fully Qualified Module Name

An **FQMN** (Fully-Qualified Module Name) is a formal name for a Hyperdrive module that can be used to disambiguate it from other modules that share the same name. They have the following format:

`<author_name>/<module_name>`

where:

- `author_name` is the complete name of the module's author, as reported in the descriptor's [Author](./descriptor.md#author-required) field.
- `module_name` is the complete name of the module itself, as reported in the descriptor's [Name](./descriptor.md#name-required) field.

For example, if the author of a module was `nodeset` and the name was `ethereumNode`, the FQMN would be `nodeset/ethereumNode`.


## Fully Qualified Parameter Name

An **FQPN** (Fully Qualified Parameter Name) is used to identify a parameter in a module's [Configuration Metadata](./config.md#metadata) or [Configuration Instance](./config.md#instance). FQPNs can take two forms:

1. `<FQMN>:<config_section_path>/.../<parameter_id>`
2. `<config_section_path>/.../<parameter_id>`

The first form is universally acceptable everywhere. It begins with the target Module's FQMN, followed by a `:`, and then the complete path to the target Parameter (described below).

The second form is shorthand for the second. It skips the module's FQMN and can be used when referencing a module's own configuration. Hyperdrive knows the context of which module's artifacts are using this notation, and can amend the FQMN to them automatically.

Pathing information related to the Sections containing a Parameter is included in these forms. For a top level parameter that doesn't reside in a Section, simply use the Parameter [ID](./config.md#common-properties) itself. Otherwise, the format is similar to a filesystem directory path; the topmost Section name is on the left, followed by a `/`, followed by a Section name, and so on all the way down to the ID of the target Parameter.


## Dynamic Properties

Dynamic properties are properties that can modify their values dynamically at runtime based on the values of the entire Hyperdrive [Configuration Instance](./config.md#instance) parameters (including all modules). They support [the standard Go `text/template` form](./templates.md) to control the property's value.

Dynamic properties are formally referred to as type `Dynamic PropertyType`, where `PropertyType` itself is a simple type (such as `bool`, `string`, `int`, and so on). They have the following properties:

- `default` (`PropertyType`, required): the default value to use for this property. If the template below is not provided, or if your template indicates that it does not apply via [`.UseDefault`](./templates.md#usedefault), then this value will be used.
- `template` (string, optional): the [template](./templates.md) to use when determining this property's value. This will be re-evaluated every time one of your module's configuration parameters has a value change, and when the user re-enters the UI for your module's configuration.