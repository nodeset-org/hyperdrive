#!/bin/sh

INSTALL_DIR=$HOME/.hyperdrive

# Make the base directory
mkdir -p $INSTALL_DIR

# Copy the install artifacts
cp -r ./deploy/templates $INSTALL_DIR
cp -r ./deploy/override $INSTALL_DIR

# Make the folders that will be used at runtime
mkdir -p $INSTALL_DIR/runtime
mkdir -p $INSTALL_DIR/data/sockets
