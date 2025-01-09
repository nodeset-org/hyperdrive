# Hyperdrive Module Specification

Hyperdrive is a software system that allows its users (Hyperdrive node operators) to quickly install, configure, execute, monitor, and otherwise participate in applications or services (called "modules") in a common, familiar environment. Many of these applications are tied to [NodeSet](https://nodeset.io), though Hyperdrive suppoirts modules that do not require it as well.

Module authors are responsible for assembling new or existing applications or services into the Hyperdrive module format. This guide covers (at a high level) the components, structure, and packaging required to meet that format. It is intended for developers that would like to become module authors, and has been written with that audience in mind.


## Docker Requirement

As Hyperdrive is intended to work across a wide gamut of Operating Systems, architectures, and Distributions, all of its running components must be containerized. Hyperdrive uses [Docker](https://docs.docker.com/get-started/docker-overview/) as its container engine and the [Docker Compose](https://docs.docker.com/compose/intro/compose-application-model/) plugin for container management. Any binary files or services you want to run must be provided in the form of one or more Docker Compose template(s).


## Components

Hyperdrive modules have the following components:

- A [descriptor](./descriptor.md) that provides information about your module, so it can be explored prior to installation.

- An [adapter](./adapter.md), in the form of a Docker Compose template file, that bridges your module's native configuration and Hyperdrive's configuration format. It must provide the configuration in Hyperdrive format, validate it prior to saving, save it so that your services can use it in their native format, and so on.

- One or more **service(s)**, in the form of Docker Compose template files. These can be as simple as a standalone binary file that responds to CLI commands, or as complex as multi-container application suites that coordinate with each other. They are essentially the components with the functionality you want your module to execute. They typically use the [template system](./templates.md) to allow for variable substitution based on the project configuration.


## Packages

Each module must be bundled inside a **Package** - a single compressed archive file in standard [ZIP](https://en.wikipedia.org/wiki/ZIP_(file_format)) format. The structure of a package must be as follows:

- `descriptor.json`: The [descriptor](./descriptor.md) file describing the module.

- `adapter.tmpl`: The Docker Compose template file for the [module adapter](./adapter.md). It will be run through Hyperdrive's [templating system](./templates.md) prior to execution.

- `templates/` A folder that contains templates of Docker Compose files for each of your service container files. Each one must end in the `.tmpl` extension, indicating they are not Docker Compose files themselves, but rather templates that will be run through Hyperdrive's [templating system](./templates.md) prior to execution.

- Any other files or folders your module needs. Hyperdrive will ignore these, but it will mount them into your service folder so your modules can use them.


## Directories and Volumes

Your module can access several different directories on the root filesystem, which its containers can mount as volumes:

- The **logging directory** should be used to store logs for your module - from both the adapter and any services. It can be accessed in the templates with [ModuleLogDir](./templates.md#modulelogdir).
- The **config directory** is where your adapter should store your module's configuration files. Your services can look for them here. Use the [ModuleConfigDir](./templates.md#moduleconfigdir) accessor to retrieve it.
- The **data directory** is where your services can store any sensitive data, such as private keys. This is only readable by your module services. The [ModuleDataDir](./templates.md#moduledatadir) accessor will provide it.


## Authentication

All communications for your module, whether it's from the user to the adapter or from the adapter to the services, must be authenticated.

Your adapter will be provided with its own [secret key](./templates.md#modulejwtkeyfile); each command that requires authentication will provide a key as part of the request, which must match this key.

Your services will be provided with a [JWT key](./templates.md#modulejwtkeyfile-1) for authenticating any incoming requests. They can use this in whatever capacity they like, but they must authenticate those requests to prevent an attacker without access to those secrets from using your services.


## Communicating with the Hyperdrive Daemon

The Hyperdrive Daemon provides endpoints for several useful features, such as restarting Docker containers on the system and retrieving the configuration for the entire Hyperdrive system (including the configurations of other modules). It's made accessible to your module via its HTTP API, which can be accessed with the [Hyperdrive Daemon URL](./templates.md#hyperdrivedaemonurl) accessor. All requests must be JWT authenticated with the contents of the [Hyperdrive API Key file](./templates.md#hyperdrivejwtkeyfile).

A full specification for the Hyperdrive API is TBD.


## How Hyperdrive Uses Modules

Below is a high-level overview of typical module usage within Hyperdrive. 


### 1. Installing a Module

Hyperdrive node operators can install a module they want to use by using the `hyperdrive module install` command. This will let them select the module they want from the curated list of officially supported modules, or let them install custom modules from third-party sources if they're trusted.

After this process completes, the module package's contents will be unpacked to the Hyperdrive system file directory:

- `/usr/share/hyperdrive/modules/<FQMN>` on Linux
- `TBD` on MacOS
- `<FQMN>` refers to the [Fully Qualified Module Name](./types.md#fully-qualified-module-name)

When installation finishes, Hyperdrive will start your adapter container so it can call on your adapter whenever it's required. If the adapter container fails to start, the module will be marked as "failed" and not be usable.


### 2. Configuring a Module

When the user runs `hyperdrive service config`, they will open the configurator UI. This UI provides a unified view for the configuration of all Hyperdrive modules. Hyperdrive's [module config specification](./config.md) allows your module to participate in this, so it can dynamically adjust its presentation based on the configuration of other modules or use the value of other module parameters in its own configuration to ensure a unified user experience.

If your adapter container hasn't been started yet, it will be started prior to showing this UI.

Hyperdrive will first call `hd-module get-config` on your adapter to retrieve your configuration and load it into its UI. Once the user wants to save it, it will call `hd-module process-config` and pass the instantiated configuration to your adapter. Your adapter must confirm that the configuration is valid; if not, it should return a list of errors as defined in the adapter specification. If valid, it should provide a list of any externally-bound TCP / UDP ports that it will use, so Hyperdrive can confirm there won't be any bind conflicts.

Next, Hyperdrive will call `hd-module set-config` with the instantiated config. Your adapter should save this, if applicable, in a format that your services expect.


### 3. Running a Module

Once Hyperdrive needs to start its services, it will call `hd-module get-containers` to retrieve the list of Docker containers for your services that are relevant and should be used (even if they're already started). It will retrieve the Docker Compose templates matching these services in your module's system folder, instantiate them with the current configuration, and save them to its list of runtime Docker Compose files. It will then call `docker compose up` against each of these files to start your module's services.

Once your services have started, the user can [interact with your adapter by sending it CLI commands](./adapter.md#hd-module-run). This is where the adapter can interact with your module's services and perform any logic it needs to do.