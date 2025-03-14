# Template Support

---

## Disclaimer

The following is a **preliminary** specification and is actively under development. The functions, types, or behaviors described here are subject to change as Hyperdrive matures, and as formal SDKs for different languages are built.

If you are interested in building a module for Hyperdrive, please contact NodeSet (info@nodeset.org) or join our Discord Server (https://discord.gg/dNshadxVkg) to engage with our development team.

---

One of the most important parts of the Hyperdrive ecosystem is its unified, modular configuration. The configuration is a single entity that encompasses the configuration for the entire Hyperdrive instance, including for all modules. This is passed between portions at various stages of operation, so modules can determine what behavior to take when starting - even if their behavior is predicated on settings housed in another module.

To provide maximum flexibility, Hyperdrive supports the standard [Go `text/template` system](https://pkg.go.dev/text/template). This is a highly customizable data-driven substitution and logic invocation system that can support anything from simply retrieving a configuration parameter to returning a custom value based on some arithmetic calculated on the fly.

Hyperdrive uses Go templates in several locations:

- The Docker Compose file template for a module's [adapter](./adapter.md)
- The [configuration metadata](./config.md#metadata) returned by the module's adapter
- Docker Compose file templates for module services

Each of these use cases has a different way of utilizing templates, which are covered in the sections below.


## Module Adapter Docker Compose Template

Templates can be used anywhere inside a Docker Compose file template for your adapter. The following methods are supported:


### AdapterContainerName

`{{.AdapterContainerName}}` has the name of the Docker container that the adapter should be created with. Your adapter should have a `container_name` property in its service with this value.


### AdapterEnvironmentVariables

`{{.AdapterEnvironmentVariables}}` is an array of `NAME=value` strings, each of which represents a line that can go under the `environment` property of your template.

For [global adapters](./adapter.md#execution-mode), the only value in this array will be `HD_ADAPTER_MODE=global`.

For [project adapters](./adapter.md#execution-mode), this will include `HD_ADAPTER_MODE=project` as well as a value for each of the [adapter project variables](./adapter.md#environment-variables). They are all provided explicitly below as well for posterity.


### AdapterMode

`{{.AdapterMode}}` is the [execution mode](./adapter.md#execution-mode) that the adapter instance will run in. This corresponds to the `HD_ADAPTER_MODE` [environment variable](./adapter.md#environment-variables).


### ModuleConfigDir

`{{.ModuleConfigDir}}` is the full path of the directory that your module should use for storing its configuration settings. It corresponds to the `HD_CONFIG_DIR` [environment variable](./adapter.md#environment-variables).


### ModuleLogDir

`{{.ModuleLogDir}}` is the full path of the directory that your module should use for storing its log files. It corresponds to the `HD_LOG_DIR` [environment variable](./adapter.md#environment-variables).


### AdapterKeyFile

`{{.AdapterKeyFile}}` is the full path to the file that contains the key to use for authenticating adapter commands while in Project mode. It corresponds to the `HD_KEY_FILE` [environment variable](./adapter.md#environment-variables).


### ModuleComposeDir

`{{.ModuleComposeDir}}` is the full path of the director that will contain all of your module's [service compose templates](./adapter.md#docker-compose-file) after they've been instantiated and saved as completed YAML files, ready for Docker Compose consumption.


### AdapterVolumes

`{{.AdapterVolumes}}` is an array of `path:path` strings, each of which represents a line that can go under the `volumes` property of your template. Your adapter will need to mount each of these for the adapter to work properly in project mode.


### ModuleNetwork

`{{.ModuleNetwork}}` is the name of the Docker Compose network that will be created for your adapter, your module services, and the artifacts from any other modules while running in project mode. In global mode this will be empty. You should use this for the `networks` property of your template, both within its `services` property for your adapter and for the top-level `networks` property.


## Configuration Metadata

Templates can be used in configuration metadata to control [dynamic properties](./types.md#dynamic-properties). Within a dynamic property template, the following methods are available:


### GetValue

`{{.GetValue <FQPN>}}` retrieves the value of the provided parameter, which is specified by its [Fully Qualified Parameter Name](./types.md#fully-qualified-parameter-name).

For example, say you want to hide a parameter using the [Hidden](./config.md#common-properties) property based on whether or not Hyperdrive had IPv6 support enabled (`hyperdrive:enableIPv6`). You could template this in your module's configuration metadata with the following:

```json
{
    ...
    "parameters": [{
        "id": "myParameter",
        ...
        "hidden": {
            "default": false,
            "template": "{{if eq .GetValue(\"hyperdrive:enableIPv6\") true}}false{{else}}true{{end}}"
        }
    }]
}
```

Any dynamic property templates will be run dynamically whenever one of the following events occur:

- The user has entered your module's configuration section
- One of the parameters in your own module's config has its value changed


### GetValueArray

`{{.GetValueArray <FQPN> <delimiter>}}` takes the value of the provided parameter, which is specified by its [Fully Qualified Parameter Name](), splits it according to the `delimiter` string, and returns the resulting array. This is useful to take parameters that represent multiple values separated by a comma, semicolon, or other delimiter, and split them into an explicit array so they can be iterated on in a template via the `for` keyword.


### UseDefault

`{{.UseDefault}}` can be used within the `template` of a dynamic property to indicate that Hyperdrive should use its `default` value. This is helpful to reduce duplication if your default value ever changes, so you don't have to set it multiple times. 


## Service Docker Compose Templates

Templates can be used anywhere inside a Docker Compose file template for service definitions. The following methods are supported:


### ModuleComposeProject

`{{.ModuleComposeProject}}` provides the Docker Compose project name for the project your service belongs to. Use this as the prefix for your service `container_name` properties to specify the name of any of your service containers for clarity.


### ModuleNetwork

`{{.ModuleNetwork}}` is the name of the Docker network that your service can use to connect to other service modules, the project adapter, the Hyperdrive daemon, and the services of any other modules.


### GetValue

`{{.GetValue <FQPN>}}` retrieves the value of the provided parameter, which is specified by its [Fully Qualified Parameter Name](./types.md#fully-qualified-parameter-name).

For example, say you want to retrieve the currently selected network (the `network` Parameter) from the standard Ethereum Node module (`nodeset/ethereumNode`) along with a custom address in your own module's config (`me/my-module:some-address`). You could template this in your module's artifacts with the following rules (using a Docker Compose YAML template excerpt in this example):

```
NETWORK: {{.GetValue nodeset/ethereumNode:network}}
FEE_RECIPIENT: {{.GetValue some-address}}
```

Whenever Hyperdrive starts its services (including the modules), the template file is run through Go's templating engine and stored within the module's configuration directory. It is then fed into the Docker Compose engine to start the container.


### GetValueArray

`{{.GetValueArray <FQPN> <delimiter>}}` takes the value of the provided parameter, which is specified by its [Fully Qualified Parameter Name](), splits it according to the `delimiter` string, and returns the resulting array. This is useful to take parameters that represent multiple values separated by a comma, semicolon, or other delimiter, and split them into an explicit array so they can be iterated on in a template via the `for` keyword.


### ModuleConfigDir

` {{.ModuleConfigDir}}` returns the full path of the directory that stores your service's configuration, as provided by the adapter during a `set-settings` operation.


### ModuleLogDir

` {{.ModuleLogDir}}` returns the full path of the directory that stores your module's log files.


### ModuleDataDir

`{{.ModuleDataDir}}` returns the full path of the directory meant for your module to use as its primary data directory. This will be a subdirectory within the user's specified data directory. Anything that needs to live persistently on the filesystem (such as supplemental user settings) that will survive after a Docker service termination should go here; large things like chain data, which do not need to survive a Docker service termination, should go into a named Docker volume instead.


### HyperdriveDaemonUrl

`{{.HyperdriveDaemonUrl}}` returns the full URL, including scheme and port, for the Hyperdrive daemon's API endpoint. Your service can use this to send HTTP API requests to Hyperdrive.


### HyperdriveJwtKeyFile

`{{.HyperdriveJwtKeyFile}}` returns the path of the file on-disk that your daemon must use as the JWT authentication secret key when sending requests to the Hyperdrive API.
