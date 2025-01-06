# Hyperdrive

Hyperdrive is an all-in-one node management system for NodeSet node operators.

Check [the NodeSet docs pages](https://nodeset-org.gitbook.io/nodeset/node-operators/hyperdrive) for full documentation and setup guides.


## Installation

Installing Hyperdrive can be done in two ways: via the `apt` package manager for Debian-based systems, or manually via the CLI (for any Linux or macOS system).


### Via the Package Manager (for Debian-based systems with `apt`)

If your system uses the `apt` package manager, you can install Hyperdrive by enabling our repository.


#### Install Docker

Start by installing Docker for your system following the [Docker installation instructions](https://docs.docker.com/engine/install/).

Next, add your user to the group of Docker administrators:
```
sudo usermod -aG docker $USER
```

Finally, exit the terminal session and start a new one (log out and back in or close and re-open SSH) for the new permissions to take effect.


#### Install Hyperdrive

1. Update the system packages and install some prerequisites:
    ```
    sudo apt update && sudo apt install curl gnupg apt-transport-https ca-certificates
    ```

2. Save the Hyperdrive repository signing key:
    ```
    sudo install -m 0755 -d /etc/apt/keyrings && sudo curl -fsSL https://packagecloud.io/nodeset/hyperdrive/gpgkey -o /etc/apt/keyrings/hyperdrive.asc
    ```

3. Add the Hyperdrive repository to your `apt` list:
    ```
    sudo tee -a /etc/apt/sources.list.d/hyperdrive.list << EOF
    deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/hyperdrive.asc] https://packagecloud.io/nodeset/hyperdrive/any/ any main
    deb-src [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/hyperdrive.asc] https://packagecloud.io/nodeset/hyperdrive/any/ any main
    EOF
    ```

4. Install Hyperdrive via `apt`:
    ```
    sudo apt update && sudo apt install hyperdrive
    ```


### Manual Install (for all systems)

If you can't or don't want to use the `apt` process, you can install Hyperdrive manually instead.

1. Download the CLI from [the latest GitHub release](https://github.com/nodeset-org/hyperdrive/releases/latest). Note there are four options: two for Linux and two for Darwin (macOS); both are available for `amd64` and `arm64`. To have parity with the package installer, we recommend saving it to `/usr/bin/hyperdrive`. For example, on an `x64` Linux system, you could do:
   ```
   sudo wget https://github.com/nodeset-org/hyperdrive/releases/latest/download/hyperdrive-cli-linux-amd64 -O /usr/bin/hyperdrive && sudo chmod +x /usr/bin/hyperdrive
   ```
    Make sure you run `chmod +x` on it before trying to use it.

2. Install Hyperdrive via the CLI:
   ```
   hyperdrive service install
   ```

This will also handle installing all of the dependencies and permissions for you.


## Updating Hyperdrive

### Via the Package Manager (for Debian-based systems with `apt`)

If you installed Hyperdrive via the package manager, you simply need to run the following to update it when a new release is out (along with any other system packages that are out of date):
```
sudo apt update && sudo apt dist-upgrade && sudo apt auto-remove
```


### Manual Update (for all systems)

If you installed Hyperdrive manually, start by downloading the new CLI using the same process you followed in step 1 of the [manual installation](#manual-install-for-all-systems) section.

Once it's downloaded, run the following command:

```
hyperdrive service install -d
```

Note the `-d` which skips Operating System dependencies, since you already have them.


## Testing

Running the integration tests requires an externally-managed instance of [Hardhat](https://hardhat.org/hardhat-runner/docs/getting-started#overview) installed and running. Among other things, Hardhat is a special standalone Execution client that simulates an Execution layer blockchain. It provides testing functions such as the ability to take snapshots of the EVM, mine new blocks, fast forward time, and so on.

We recommend using the [Hardhat start script](https://github.com/nodeset-org/osha/blob/main/hardhat/start.sh) provided in the [OSHA repository](https://github.com/nodeset-org/osha) if you don't already have experience running Hardhat.


### Configuration

When running the tests, make sure you follow two important rules:
1. Set the `HARDHAT_URL` environment variable to the URL of your instance prior to running tests
2. Set package parallelization to 1 (so multiple tests don't run concurrently, which will break Hardhat's snapshotting system)

For example:

`HARDHAT_URL="http://localhost:8545" go test -p 1 ./...`

During testing you'll notice Hardhat's node will print many event messages such as `evm_snapshot`. This is expected and part of the test suite, as each test that modifies the EL state will snapshot Hardhat prior to execution and revert to the snapshot at the end.


## Attribution

Adapted from the [Rocket Pool Smart Node](https://github.com/rocket-pool/smartnode) with love.