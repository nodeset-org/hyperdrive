# Hyperdrive Module Adapter Specification

---

## Disclaimer

The following is a **preliminary** specification and is actively under development. The functions, types, or behaviors described here are subject to change as Hyperdrive matures, and as formal SDKs for different languages are built.

If you are interested in building a module for Hyperdrive, please contact NodeSet (info@nodeset.org) or join our Discord Server (https://discord.gg/dNshadxVkg) to engage with our development team.

---

Every Hyperdrive module is required to come with an **adapter** - a binary executable that serves two purposes:
- It can perform certain functions *without* the main module service running (i.e., it can run in standalone mode).
- It can handle interactive execution performed by the user, which *may* or *may not* be executed without the main module service running.

Many modules rely on binaries and services that have already been built without Hyperdrive support, but want to be integrated into the Hyperdrive ecosystem. In such cases, the CLI can act as a simple "shim" that converts Hyperdrive activities and data into something the existing program can handle. For example, it can convert Hyperdrive's configuration into a form that the existing module's binaries can use, or it can route interactive user commands down to the existing services accordingly.

This specification describes the way Hyperdrive will invoke the CLI binary in various situations, along with the required functions that must be implemented.


## Execution Environment

Your adapter must come in a **Docker container** with all of its prerequisites installed. This way you can write it in whatever language or fashion you like; as long as it contains an executable binary that conforms to the standards below, Hyperdrive will be able to use it.

Since your adapter enters a "wait" mode upon startup, the adapter binary itself will be persistently available for quick and easy access via `docker exec` during operations such as configuration changes. Because of this process, however, your adapter will run before Hyperdrive itself has been configured. Anything (including project name, data directory, API ports, etc.) can change after the fact. Therefore, **your adapter should be stateless** (or at least, not rely on the state of the Hyperdrive configuration).


### Execution Mode

Adapters are executed in one of two different modes:

1. Global mode
1. Project mode

When a module is installed, Hyperdrive will create and run one instance of your adapter's Docker container in **global mode**. This container is not tied to any one Hyperdrive project; all Hyperdrive installations on the node will all interact with it to use your adapter. Global mode commands run independently from individual Hyperdrive projects and are not affected by a project's configuration. They will *not* modify the behavior of your module and thus will not interact with your module's service containers; they are effectively **read-only** commands that will only access things within the adapter itself. Examples include getting the adapter's version, validating a pending configuration, and so on.

Once the user configures a Hyperdrive project, the system will create a *second* instance of your adapter's Docker container. This one will be provided with all of the project-specific details, such as the paths for the configuration and log directories, and the path for the authentication key file (discussed later). Project mode commands should be considered **read and write**, because they can be used to interact with your module's services and modify its configuration.


### Authentication

Commands run on your adapter while in project mode all require authentication. When your adapter container is started in Project mode, it will have the `HD_KEY_FILE` environment variable set. This will hold the full path of a file that contains the string to use for authenticating requests. The request input format for each Project mode command will be in JSON, and include a field named `key` with a string value. The value of `key` should be directly compared to the contents of the `HD_KEY_FILE` file. Only requests with a matching value should be permitted; otherwise your adapter should return an error.


### Docker Compose File

Your module package must include a top-level file named `adapter.tmml`, which is the template for a Docker Compose file that is used to create and run your module's adapter container. This  will be passed through Hyperdrive's [variable substitution](./templates.md) system prior to instantiation.

This file can have any set up required, with the following conditions:

- The `entrypoint` must be the path to your adapter binary, along with any flags required for it to run. This will be used by Hyperdrive as a command prefix when running one of the commands below.
- The `command` arguments must launch your adapter in a mode where it simply **sleeps and idles indefinitely** until Docker tells the container to stop (`SIGTERM` by default), at which point it should exit gracefully to allow the container to stop. It should not take any other behavior or consume any resources beyond what is needed to make the process sleep until it receives the stop signal.
- Logs must be written to the directory provided by the [`ModuleLogDir`](./templates.md#module-adapter-docker-compose-template) template function.
- Many of the commands are authenticated and include a `key` property in the input. For these calls, your adapter must compare them to the contents of the file provided in the [`AdapterKeyFile`](./templates.md#adapterkeyfile) to ensure the caller is permitted to proceed.
- The `container_name` should be set to the [`AdapterContainerName`](./templates.md#adaptercontainername) variable.
- If [`ModuleNetwork`](./templates.md#modulenetwork) is present, your adapter should have it in the `networks` section of its service definition and an entry for it in the top-level `networks` section.


### Environment Variables

Hyperdrive will always start your adapter with the `HD_ADAPTER_MODE` environment variable set. This variable can have two values:

- `global` indicates that it's running in Global mode.
- `project` indicates that it's running in Project mode.

When running in Project mode, these additional environment variables are set:

- `HD_CONFIG_DIR` is the full path (within the container) to the directory that your module should use for storing its configuration settings during a `set-settings` command. Your module's services should load them from this path, which will be passed into their Docker compose instances as well.

- `HD_LOG_DIR` is the full path (within the container) to the directory that your module should use for storing any log files. Your module's services should load them from this path, which will be passed into their Docker compose instances as well.

- `HD_DATA_DIR` is the full path (within the container) to the directory that your module should use for storing sensitive data that should be kept secret. The user will specify a root directory for all sensitive data to preside in, which may be a separate encrypted drive for example; your module must honor this setting to prevent accidentally leaking that data.

- `HD_KEY_FILE` is the full path (within the container) to the file that contains the key to use for authenticating adapter commands while in Project mode.

- `HD_COMPOSE_DIR` is the full path (within the container) to the directory that contains the Docker Compose files for your modules after their templates have been instantiated. The filenames will be the same as the Docker Compose templates contained within your [module package](./module.md#packages), but the extension will be `.yml` instead of `.tmpl`. Your adapter can use these files when starting or stopping your module.

- `HD_COMPOSE_PROJECT` is the name of the project to use as the Docker Compose project name when running Docker Compose commands.


## Adapter API Protocol

Communication between Hyperdrive and the adapter will be done via `STDIN` and `STDOUT`, and delimited by newlines. Your adapter will be invoked with one of the following commands via `docker exec`; each input and output format is defined as part of the function.


# Hyperdrive Adapter Command Specification

The following are commands that the Hyperdrive system itself will call on your module. Your module must be able to handle executing things as a one-off (standalone) mode; in other words, your adapter's Docker container will not start normally in persistent mode when these commands are called, but rather will be run with the specified commands below as the entrypoint; upon finishing the command, the container is then immediately discarded.

**NOTE:** The example input / output JSON below uses newlines for readability, but the JSON input to your program will not; it will all be on one line as newlines are to be treated as "end of input" characters while reading. Your adapter's output may or may not include them as you see fit.


## Global Mode Commands

The following commands are called in **Global mode**.


### `hd-module version`

When executed, this function must return the version of the CLI binary. It does not take any input; it should simply return the CLI version and exit.


#### Input

(None)


#### Output:

```json
{
    "version": "1.2.3"
}
```

where:

- `version` must be a [semantic version](https://semver.org/) string.

The version is expected to have parity with the Hyperdrive module version, as defined in [the descriptor](./descriptor.md). This is simply used for sanity checking prior to startup; mismatches between the two will indicate a module installation / upgrade failure and be reported as errors to the user during startup.


### `hd-module upgrade-instance`

This will be called when Hyperdrive detects that its instance for your module was generated with an older version of the module, and the user has updated the module but hasn't run through the configuration process yet. It should migrate the old [Module Instance](./config.md#instances) to the latest version. If no changes are required, then you can simply return the old configuration with the `version` updated.

This is an opportunity to modify deprecated settings, or invalidate obsolete ones that no longer apply by making them blank or using the default value for example. The user will be informed of what has changed prior to saving the configuration so they have the option to cancel the process and revert the upgrade; your module must not save this configuration until Hyperdrive calls the [set-config](#hd-module-set-config) function.

Furthermore, this is a way for your module to enforce that the instance is set to **disabled** in the event that the user needs to loearn about some kind of breaking changes that would preclude it from working after the upgrade.


#### Input

```json
{
    "instance": {
        ...
    }
}
```

where:

- `instance` is a [Module Instance](./config.md#instances) for your module specifically.


#### Output

A serialized JSON [Module Instance](./config.md#instances) representing the instance after your module has upgraded it. Any parameters that are no longer compatible with the current instance changed to an appropriate value (such as the new default). The `version` of the returned instance must also be updated to the latest version of your module, as reported in its [descriptor](./descriptor.md) file and the [version](#hd-module-version) function.


### `hd-module get-config-metadata`

This should return a [Hyperdrive Configuration Metadata](./config.md#metadata) object for your module as a serialized string in JSON format.


#### Input

(None)


#### Output

A serialized JSON [Hyperdrive Configuration Metadata](./config.md#metadata) object for your module's configuration.



### `hd-module process-settings`

This should process the [Settings](./config.md#settings) of your module's configuration, extracting important information and validating it. This will be called after the user has modified Hyperdrive's configuration (potentially including your module's configuration), but before the configuration needs to be saved.

Your module should use this to return information about the configuration and validate that the provided configuration meets all of your module's requirements and is valid.


#### Input

```json
{
    "oldSettings": {
        ...
    },
    "newSettings": {
        ...
    }
}
```

where:

- `currentSettings` are the current settings for the [complete Hyperdrive installation](TODO), including your module's configuration and the configuration for all other installed modules. They are provided for reference.
- `newSettings` are the newly proposed settings for the [complete Hyperdrive installation](TODO), including your module's configuration and the configuration for all other installed modules. They are provided so your module can compare them against the current settings to determine which services need to be restarted.


#### Output

The following serialized JSON object:

- `errors` (string[], required): A list of error messages to provide to the user when the settings fail validation. Each one should be a reason why the configuration is invalid. They will be displayed directly to the user so they should be human-readable strings. If there are no errors and the settings are valid, this should be an empty array.
- `ports` (object, required): A mapping for externally available TCP/UDP ports that your module's services will bind when running. Each property in the object must have the FQMN of one of your module's properties as its name, and the port value as its value. This is used by Hyperdrive to ensure that your service won't bind ports that are already in use by other services. If your ports are not externally bound, and restricted to Docker's internal network, they don't need to be listed here. This list can be empty if `errors` is not empty for simplicity.
- `servicesToRestart` (string[], required): A list of Docker container service IDs (as written in your module's [service templates](./module.md#packages)) that need to be restarted (if already running) for these settings to take effect. This will be displayed to the user for informational purposes, but Hyperdrive will not actually restart these containers as part of issuing this command. Your adapter doesn't need to check whether or not the service is currently running; it can simply report that a service restart is required even if the service is stopped or doesn't exist.

For example:

```json
{
    "errors": [
        "Remote Logging URL must be a valid URL.",
    ],
    "ports": {
        "publicApiPort": 1234,
        "serverListenerPort": 5678,
    },
    "servicesToRestart": [
        "my-service",
        "my-other-service"
    ],
}
```


## Project Mode Commands

The following commands will be called when your adapter runs in **Project mode**, and belongs to a specific Hyperdrive project.


### `hd-module set-settings`

This will be called prior to starting / restarting your module's services. It will provide the settings for the entire [complete Hyperdrive installation](TODO) in serialized JSON format to `STDIN` (terminated with an empty `\n` character); your adapter must read this properly. Your CLI can then save whatever configuration it needs in a format your module services can interpret (if they do not already pull the Hyperdrive configuration from its daemon on startup).

This configuration is guaranteed to be valid according to the `process-settings` command above, as that will be called prior to this.

No response is expected from this command during a successful run. If any errors occur while saving the config, they should be printed to `STDERR`. Hyperdrive will detect them and indicate a configuration failure to the user, then abort the startup procedure.


#### Input

```json
{
    "key": "...",
    "settings": {
        ...
    }
}
```

where:

- `key` must match the contents of the file in `HD_KEY_FILE`.
- `settings` are the settings for the [complete Hyperdrive installation](TODO), including your module's configuration and the configuration for all other installed modules.


#### Output

(None)


### `hd-module start`

*NOTE: this command is currently a WIP and designed for MVP purposes only. The final form is expected to change.*

Hyperdrive calls this command after your module's configuration settings have been saved with `set-settings`. Your adapter should start (or restart) whatever Docker containers it needs to based on the current settings (which will be provided with the call). For MVP purposes your adapter should execute a `docker compose up` command with the relevant files passed in based on the provided settings. Your module's service file templates will be instantiated and placed into the `HD_COMPOSE_DIR` directory prior to this call. It should use Docker Compose to start the services; for the Compose project name, use the `HD_COMPOSE_PROJECT` property.


#### Input

```json
{
    "key": "...",
    "settings": {
        ...
    },
}
```

where:

- `key` must match the contents of the file in `HD_KEY_FILE`.
- `settings` are the settings for the [complete Hyperdrive installation](TODO), including your module's configuration and the configuration for all other installed modules.


#### Output

If the operation worked successfully, your adapter should return with Exit Code 0. If it failed, you should print an error to STDERR and return a non-zero Exit Code. Hyperdrive itself doesn't consume any printed output of this command; it will simply be displayed to the user.


### `hd-module run`

This is called when the user wants to run a command on your adapter via Hyperdrive's CLI. Since your adapter effectively serves as your module's CLI, commands here will be forwarded to your adapter. These commands can be anything, from a simple "print the module adapter's help text" to complete functions that affect multiple systems.

The user will invoke these commands with one of the following syntax options:

- `hyperdrive <module shortcut> <command>`
- `hyperdrive <module name> <command>`
- `hyperdrive <FQMN> <command>`

The first format uses the [module shortcut](./descriptor.md#structure) for convenience. If two or more modules have conflicting shortcuts, it will not be available.

The second format uses the [module name](./descriptor.md#structure) itself. If two or more modules have conflicting module names, it will not be available.

The third format uses the [Fully Qualified Module Name](./types.md#fully-qualified-module-name) as a way to ultimately disambiguate any conflicts. While verbose, this syntax will always invoke the correct module.


#### Input

```json
{
    "key": "...",
    "command": "..."
}
```

where:

- `key` must match the contents of the file in `HD_KEY_FILE`.
- `command` is the command to run (without your adapter binary name).

**NOTE:** This function will be run *interactively*, meaning it can prompt the user for input if necessary.


#### Output

This should output whatever output your command has; it will be viewed directly by the user so it doesn't need to be in JSON format. Any errors that occur should be printed to `STDERR`.














### `hd-module call-config-function`

This is called when the user wants to run a command on your adapter to dynamically compute values for configuration/template substitutions.


#### Input

```json
{
    "funcName": "...",
}
```

where:

 - `funcName` is the name of the adapter-defined function to invoke.



#### Output

```json
{
    "result": "...",
}
```

where:

 - `result` is the string containing the returned computed value.
