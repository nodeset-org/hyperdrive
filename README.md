# Hyperdrive

Hyperdrive is an all-in-one node management system for NodeSet node operators.

Check [the NodeSet docs pages](https://nodeset-org.gitbook.io/nodeset/node-operators/hyperdrive) for full documentation and setup guides.

## Installation

We provide packaged versions of each release so you can manage your installation via your system package manager. If you've never installed Hyperdrive before, first you should add our package repository to your list.

E.g. for Debian:
`sudo add-apt-repository 'deb https://packagecloud.io/nodeset/hyperdrive'`

`sudo apt-get update`

`sudo apt-get install hyperdrive`

To finalize the installation, define your configuration by running `hyperdrive service install`, then `hyperdrive service config`.

## Updating

Update via your package manger. E.g. for Debian:

`sudo apt update && sudo apt dist-upgrade && sudo apt autoremove`

## Attribution

Adapted from the [Rocket Pool Smart Node](https://github.com/rocket-pool/smartnode) with love.
