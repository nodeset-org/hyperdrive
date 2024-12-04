# Hyperdrive Module Configuration Specification

One of the most important aspects of a Hyperdrive module is its **configuration**. This refers to a series of named settings or values that the user can modify (or in some cases, are modified behind-the-scenes based on the configuration of other Hyperdrive modules, or even the hardware / software configuration of the node machine *itself*). As with all configuration systems, these settings will inform behavioral changes to your module's services during runtime.

In the context of Hyperdrive, there are two entities related to configuring your module: **metadata** and **instances**.


## Metadata

**Configuration Metadata** refers to information *about* the configuration. Things such as parameter names, types, and descriptions all fall under metadata. This contains everything necessary to dynamically generate a UI for configuring your module interactively. Hyperdrive will request this from your module's adapter at various stages. The metadata doesn't need to be the same every time, especially if it leverages the [templating feature](./templates.md) to allow dynamic changes during runtime (described later in this section).

A configuration metadata object consists of the following two properties, each of which can have zero or more of the following entities:

- `parameters` ([Metadata Parameter\[\]](#parameters), required): an array of metadata for individual settings that will ultimately be treated as key-value pair entries when the configuration is [instantiated](#instances).

- `sections` ([Metadata Section\[\]](#section)): an array of groups of parameters used purely for organizational purposes. They can also have their own section children underneath them.


### Section

A **Section** adheres to an individual section of the configuration (a grouping of parameters and subsections). It should have the following properties:

- `name` ([Identifier](./types.md#identifier), required): a human-readable name for the section.

- `description` ([Dynamic string](./types.md#dynamic-properties) - [Description](./types.md#description), required): the human-readable description of the section. This will correspond to the value of the description box when the individual property is being configured in Hyperdrive's interactive configurator. It can be a string with any characters.

- `parameters` ([Metadata Parameter\[\]](#parameters), required): an array of metadata for individual settings that will ultimately be treated as key-value pair entries when the configuration is [instantiated](#instances).

- `sections` ([Metadata Section\[\]](#section), required): an array of groups of parameters used purely for organizational purposes. They can also have their own section children underneath them.

- `disabled` ([Dynamic bool](./types.md#dynamic-properties), optional): indicates whether this section should remain visible, but not interactive ("grayed out") in the configurator UI. If this is true, we suggest updating the description to indicate *why* it's disabled for the user's own knowledge.

- `hidden` ([Dynamic bool](./types.md#dynamic-properties), optional): prevents the section from appearing in the configurator UI at all. This is helpful if you have parameters or internal system settings that need to be changed depending on various other configuration settings, but should not be user-accessible.


### Parameters

A **Parameter** is an individual configuration setting that will be represented as a key-value pair. Parameters can have many different value types, provide their own (self-contained) validation routines, and even depend on other Parameters.

The type of a parameter's value is dictated by its `type` property. Some Parameters have different fields than others based on this `type`. These type-specific properties are defined after this section.


#### Common Properties

Below is a list of the common properties that all parameters share:

- `id` ([Identifier](./types.md#identifier), required): a unique identifier for the parameter. This is *not* presented to the user; it is only used internally as a way to reference the parameter.

- `name` (string, required): the human-readable name for the property. This will correspond to the label of the individual property field that will be presented to the user during interactive configuration. It can be a string with any characters, though we recommend it consists only of printable ones for legibility's sake.

- `description` ([Dynamic string](./types.md#dynamic-properties) - [Description](./types.md#description), required): the human-readable description of the property. This will correspond to the value of the description box when the individual property is being configured in Hyperdrive's interactive configurator. It can be a string with any characters.

- `type` (enum, required): defines how to interpret the value of the parameter, and the kind of UI element that will represent it in the configurator. This also has some bearing on how the parameter will be serialized, and thus what behavior should be taken to deserialize it. See the [Parameter Types](#parameter-types) section below for more information on the various allowed types.

- `default` (Parameter Type, required): the default value to assign to this parameter. The type of the value must correspond to the Parameter Type used above.

- `value` (Parameter Type, required): the current value assigned to this parameter according to the latest saved configuration. The type of the value must correspond to the Parameter Type used above.

- `advanced` (bool, optional): indicates whether this parameter should be hidden from the user unless they've entered "advanced mode" during service configuration, where all options are present. If not provided, defaults to `false`.

- `disabled` ([Dynamic bool](./types.md#dynamic-properties), optional): indicates whether this parameter should remain visible, but not interactive ("grayed out") in the configurator UI. If this is true, we suggest updating the description to indicate *why* it's disabled for the user's own knowledge.

- `hidden` ([Dynamic bool](./types.md#dynamic-properties), optional): prevents the parameter from appearing in the configurator UI at all. This is helpful if you have parameters or internal system settings that need to be changed depending on various other configuration settings, but should not be user-accessible.

- `overwriteOnUpgrade` (bool, required): causes the parameter's current value to be replaced with the default value when Hyperdrive detects a version upgrade has been installed, but not yet applied (the services haven't been started with the new version yet). This is helpful for things that routinely change with new versions, such as Docker container tags.

- `affectsContainers` (string[], required): an array of strings, each of which is the base name of a Docker container your module owns (without the project name prefix that Hyperdrive appends) that will be affected by changing this parameter. Affected containers will be restarted after the new configuration is saved.


#### Parameter Types

The following types of parameters are allowed:

- `bool` indicates a simple boolean value. Valid values are `true` and `false`. Boolean parameters will be represented in the configurator UI with a checkbox.

- `int` is the type for signed integer values. Valid values can be any number that fits within a signed 64-bit integer. If you need to limit the size of the value, such as to a 16-bit signed integer, use the `minValue` and `maxValue` properties of the parameter. These parameters will be represented in the configurator UI with a textbox.

    Additional properties:

    - `minValue` (optional): an `int` that represents the minimum value the parameter can have. If the parameter is less than this, it will fail the internal validation check (prior to calling your module adapter's validation routine).

        To specify that there is no `minValue`, exclude this property when serializing the configuration.

    - `maxValue` (optional): an `int` that represents the maximum value the parameter can have. If the parameter is greater than this, it will fail the internal validation check (prior to calling your module adapter's validation routine).

        To specify that there is no `maxValue`, exclude this property when serializing the configuration.

- `uint` is the type for unsigned integer values. Valid values can be any number that fits within an unsigned 64-bit integer. If you need to limit the size of the value, such as to a 16-bit unsigned integer for describing a network TCP/UDP port, use the `minValue` and `maxValue` properties of the parameter. These parameters will be represented in the configurator UI with a textbox.

    Additional properties:

    - `minValue` (optional): a `uint` that represents the minimum value the parameter can have. If the parameter is less than this, it will fail the internal validation check (prior to calling your module adapter's validation routine).

        To specify that there is no `minValue`, exclude this property when serializing the configuration.

    - `maxValue` (optional): a `uint` that represents the maximum value the parameter can have. If the parameter is greater than this, it will fail the internal validation check (prior to calling your module adapter's validation routine).

        To specify that there is no `maxValue`, exclude this property when serializing the configuration.

- `float` is the type for a double-precision (64-bit), signed floating-point number in standard IEEE 754 format. These parameters will be represented in the configurator UI with a textbox.

    Additional properties:

    - `minValue` (optional): a `float` that represents the minimum value the parameter can have. If the parameter is less than this, it will fail the internal validation check (prior to calling your module adapter's validation routine).

        To specify that there is no `minValue`, exclude this property when serializing the configuration.

    - `maxValue` (optional): a `float` that represents the maximum value the parameter can have. If the parameter is greater than this, it will fail the internal validation check (prior to calling your module adapter's validation routine).

        To specify that there is no `maxValue`, exclude this property when serializing the configuration.

- `string` is the type for general-purpose strings. The input can be anything; if you want to specify its format, use the `regex` and/or `maxLength` properties. These parameters will be represented in the configurator UI with a textbox.

    Additional properties:

    - `maxLength` (optional): a `uint` that dictates that maximum number of characters (in UTF-8 format) the value is allowed to have. If the parameter's length is greater than this, it will fail the internal validation check (prior to calling your module adapter's validation routine).

        To specify that there is no `maxLength`, exclude this property when serializing the configuration.

    - `regex` (optional) a `string` that provides a standard regular expression to be used when validating the parameter. the maximum value the parameter can have. If the parameter is greater than this, it will fail the internal validation check (prior to calling your module adapter's validation routine).

        To specify that there is no `regex`, exclude this property when serializing the configuration.

- `choice` represents a parameter that's a choice between multiple fixed options, such as an enum. These parameters will be represented in the configurator UI with a dropdown selection box. Choice parameters are variably typed, with a value of type `ChoiceType`. This must be a simple type - any of the other property types above can be used as a `ChoiceType`.

    Additional properties:

    - `options` ([Option\[\]](#parameter-options), required): an array of items that can be valid values for the parameter. Any value that does not correspond to one of these will be considered invalid. Each of these will be included in the property's dropdown unless explicitly hidden (see the Option type below for more details).


#### Parameter Options

`choice` parameters provide a set of metadata for their options, each of which corresponds to a value they are allowed to have. As a reminder, they are variably typed (using `ChoiceType`, which can be any of the other property types). All options must have values of that same type. The Option type has the following properties:

- `name` (string, required): the human-readable name for this option. This will correspond to the value of the option that will be presented to the user within the property's dropdown list. It can be a string with any characters, though we recommend it consists only of printable ones for legibility's sake.

- `description` ([Dynamic string](./types.md#dynamic-properties), required): the human-readable description of the option. This will correspond to the value of the description box when the option is highlighted in the dropdown list during selection of the property's value. It can be a string with any characters. For formatting, see the common [description](#common-properties) property details.

- `value` (`ChoiceType`, required): the value to assign to the property when this option is selected.

- `disabled` ([Dynamic bool](./types.md#dynamic-properties), optional): indicates whether this option should remain visible, but not interactive ("grayed out") in the configurator UI. If this is true, we suggest updating the description to indicate *why* it's disabled for the user's own knowledge.
  
- `hidden` ([Dynamic bool](./types.md#dynamic-properties), optional): when `true`, the option will be hidden from the dropdown list in the configurator UI. This is helpful if you have parameters or internal system settings that need to be changed depending on various other configuration settings, but should not be user-accessible. It's also helpful if you want to conditionally hide this option based on the other configuration settings.


## Instances

A **Configuration Instance** is an object that corresponds to a "filled out" version of a Configuration Metadata object. It has an identical Section and Parameter structure, but instead of Parameters it has **Parameter Instances**. which are simply key-value pairs. These indicate the value assigned to each parameter without any of the extraneous metadata. The type of value in each Parameter Instance is the same as the type of the corresponding Parameter.

A configuration instance object consists of the following properties, which must align exactly with the corresponding Metadata object:

- `parameters` ([Parameter Instance\[\]](#parameter-instance), required): an array of key-value pairs, where each one's name corresponds to the ID of a [Parameter Metadata](#parameters) object and the value is the value of that corresponding parameter.

- `sections` ([Instance Section\[\]](#section-1)): an array of groups of parameter instances used purely for organizational purposes. They can also have their own section children underneath them.


### Section

A **Section** in an instance is hierarchically identical to a [metadata section]() 

- `name` ([Identifier](./types.md#identifier), required): the name of the corresponding section in the Configuration Metadata object.

- `parameters` ([Parameter Instance\[\]](#parameter-instance), required): an array of key-value pairs, where each one's name corresponds to the ID of a [Parameter Metadata](#parameters) object and the value is the value of that corresponding parameter.

- `sections` ([Instance Section\[\]](#section-1)): an array of groups of parameter instances used purely for organizational purposes. They can also have their own section children underneath them.


### Parameter Instance

Parameter instances are simple key-value pair properties. The key corresponds to the ID of the parameter metadata it represents. The value is the value assigned to the instance of that parameter metadata. The type of the value is dictated by the parameter metadata.
