# Template Support

One of the most important parts of the Hyperdrive ecosystem is its unified, modular configuration. The configuration is a single entity that encompasses the configuration for the entire Hyperdrive instance, including for all modules. This is passed between portions at various stages of operation, so modules can determine what behavior to take when starting - even if their behavior is predicated on settings housed in another module.

To provide maximum flexibility, Hyperdrive supports the standard [Go `text/template` system](https://pkg.go.dev/text/template). This is a highly customizable data-driven substitution and logic invocation system that can support anything from simply retrieving a configuration parameter to returning a custom value based on some arithmetic calculated on the fly.

Hyperdrive uses Go templates in several locations:

- The Docker Compose file template for a module's [adapter](./adapter.md)
- The [configuration metadata](./config.md#metadata) returned by the module's adapter
- Docker Compose file templates for module services

Each of these use cases has a different way of utilizing templates, which are covered in the sections below.l


## Module Adapter Docker Compose Template

Templates can be used anywhere inside a Docker Compose file template for your adapter. The following methods are supported:


### ModuleConfigDir

`{{.ModuleConfigDir}}` retrieves the full path to the directory that your adapter should save its configuration (and any other extraneous files) into. This will also be made available to your service Docker Compose templates so it can be mounted as a volume for file retrieval.


### ModuleSecretFile

`{{.ModuleSecretFile}}` retrieves the full path of the file that contains the secret key used to authenticate requests to your adapter. Your adapter file can mount this as a volume so it can be read and compared against the secret provided with any incoming requests.


### ModuleLogDir

`{{.ModuleLogDir}}` retrieves the full path of the standard Hyperdrive logging directory that your module should write its log files to in order for the `hyperdrive service adapter-logs` command to work properly.


### ModuleJwtKeyFile

`{{.ModuleJwtKeyFile}}` returns the path of the file on-disk that should be used as the JWT authentication secret key by your module's services. If your services require JWT authentication for any HTTP-based interactions, your adapter can use this to load the secret and authenticate with them.


## Configuration Metadata

Templates can be used in configuration metadata to control [dynamic properties](./types.md#dynamic-properties). Within a dynamic property template, the following methods are available:


### GetValue

`{{.GetValue <FQPN>}}` retrieves the value of the provided parameter, which is specified by its [Fully Qualified Parameter Name](./types.md#fully-qualified-parameter-name).

For example, say you want to retrieve the currently selected network (the `network` Parameter) from the standard Ethereum Node module (`nodeset/ethereumNode`) along with a custom address in your own module's config (`me/my-module:some-address`). You could template this in your module's artifacts with the following rules (using a Docker Compose YAML template excerpt in this example):

```
NETWORK: {{.GetValue nodeset/ethereumNode:network}}
FEE_RECIPIENT: {{.GetValue some-address}}
```

Any dynamic property templates will be run dynamically whenever one of the following events occur:

- The user has entered your module's configuration section
- One of the parameters in your own module's config has its value changed


### GetValueArray

`{{.GetValueArray <FQPN> <delimiter>}}` takes the value of the provided parameter, which is specified by its [Fully Qualified Parameter Name](), splits it according to the `delimiter` string, and returns the resulting array. This is useful to take parameters that represent multiple values separated by a comma, semicolon, or other delimiter, and split them into an explicit array so they can be iterated on in a template via the `for` keyword.


### UseDefault

`{{.UseDefault}}` indicates that your template for a property does not apply, or the conditions for it are not met, and Hyperdrive should use the default property value instead.


## Service Docker Compose Templates

Templates can be used anywhere inside a Docker Compose file template for service definitions. The following methods are supported:


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

` {{.ModuleConfigDir}}` returns the full path of the directory that stores your module's configuration, along with any other extraneous files, as saved by your module adapter.


### ModuleDataDir

`{{.ModuleDataDir}}` returns the full path of the directory meant for your module to use as its primary data directory. This will be a subdirectory within the user's specified data directory. Anything that needs to live persistently on the filesystem (such as supplemental user settings) that will survive after a Docker service termination should go here; large things like chain data, which do not need to survive a Docker service termination, should go into a named Docker volume instead.


### HyperdriveDaemonUrl

`{{.HyperdriveDaemonUrl}}` returns the full URL, including scheme and port, for the Hyperdrive daemon's API endpoint. Your service can use this to send HTTP API requests to Hyperdrive.


### ModuleJwtKeyFile

`{{.ModuleJwtKeyFile}}` returns the path of the file on-disk that should be used as the JWT authentication secret key by your service.


### HyperdriveJwtKeyFile

`{{.HyperdriveJwtKeyFile}}` returns the path of the file on-disk that your daemon must use as the JWT authentication secret key when sending requests to the Hyperdrive API.
