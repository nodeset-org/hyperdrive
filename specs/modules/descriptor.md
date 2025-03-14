# Hyperdrive Module Descriptor Specification

---

## Disclaimer

The following is a **preliminary** specification and is actively under development. The functions, types, or behaviors described here are subject to change as Hyperdrive matures, and as formal SDKs for different languages are built.

If you are interested in building a module for Hyperdrive, please contact NodeSet (info@nodeset.org) or join our Discord Server (https://discord.gg/dNshadxVkg) to engage with our development team.

---

This is the specification for a Hyperdrive module's descriptor file. This file should be bundled with a module and describe its characteristics.

A descriptor file must be written in JSON format.


## Structure

The descriptor has the following top-level properties:


- `name` ([Identifier](./types.md#identifier), required): a human-readable name for the module. It will also serve as the module's main ID (when combined with the module author during conflicts). It should be unique to your module. Its format must adhere to the [Identifier](./types.md#identifier) rules.

- `shortcut` ([Identifier](./types.md#identifier), required): a custom shortened version of the name. This will be used by the Hyperdrive CLI when the user is invoking your module, and as a prefix in your module's Docker container names. You should pick something unique to your module. Its format must adhere to the [Identifier](./types.md#identifier) rules.

- `description` (string, required): a human-readable description of what the module does, exclusively used for node operators to reference when viewing the module's metadata. It is a generic string; all characters are valid.

- `version` (string, required): the version of the module that's described by this descriptor. It must be a valid [semantic version](https://semver.org/).

- `author` (string, required): a human-readable name for the individual, company, or other entity that created / currently maintains the module. It will be displayed to the user while viewing module metadata, and combined with the name to create a unique name during module conflicts. Its format must adhere to the [Identifier](./types.md#identifier) rules.

- `url` (string, optional): an optional URL to some website that provides more information about the module. It can be a homepage, a source code repository, a link to documentation, or anything of that nature. It will be displayed to the user while viewing module metadata. It can be any valid URL string.

- `email` (string, optional): an optional e-mail address that can be used to contact the module's maintainers. It will be displayed to the user while viewing module metadata. It can be any valid e-mail address.

- `dependencies` ([Dependency\[\]](#dependencies), required): an array of dependency strings indicating which Hyperdrive modules your module requires in order to work properly. While this property itself is required, it can be an empty array if your module does not have any dependencies. Dependency strings are described in detail below.


## Dependencies

Dependency strings can have two forms:

- `<FQMN>`
- `<FQMN> <op> <version>`

where `<FQMN>` refers to the [Fully Qualified Module Name](./types.md#fully-qualified-module-name). The full version described in the second form describes the module's author, the module's name, an operator used to compare installed versions against the specified version, and the specified version.

The operator can be any of the following:

- Less than (`<`)
- Less than or equal (`<=`)
- Equal to (`=`)
- Greater than or equal to (`>=`)
- Greater than (`>`)

Version must be a valid [semantic version](https://semver.org/). The version of the dependency that's currently installed will be compared to the specified version using any one of the above operators, and upgraded (or resulting in a dependency conflict) accordingly.

The first form simply omits the operator and version.

For instance:

```
nodeset/ethereumNode
```

Indicates that any version of `ethereumNode` authored by `nodeset` is required to run your module.

```
nodeset/ethereumNode >= 1.2.3
```

Indicates version `1.2.3` or higher of `ethereumNode` is required to run your module. 


## Examples

The following is a complete example of a JSON descriptor:

```json
{
    "name": "demo-module",
    "shortcut": "dm",
    "description": "This is a simple demo of a Hyperdrive module.",
    "version": "1.0.0",
    "author": "nodeset",
    "url": "https://nodeset.io",
    "email": "info@nodeset.io",
    "dependencies": [
        "nodeset/some-dependency",
        "nodeset/another-dependency >= 1.2.3"
    ]
}
```