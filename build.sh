#!/bin/bash

# This script will build all of the artifacts involved in a new Hyperdrive release.
# NOTE: You MUST put this in a directory that has the `hyperdrive` repository cloned as a subdirectory.

# =================
# === Functions ===
# =================

# Print a failure message to stderr and exit
fail() {
    MESSAGE=$1
    RED='\033[0;31m'
    RESET='\033[;0m'
    >&2 echo -e "\n${RED}**ERROR**\n$MESSAGE${RESET}\n"
    exit 1
}


# Builds all of the CLI binaries
build_cli() {
    cd hyperdrive || fail "Directory ${PWD}/hyperdrive does not exist or you don't have permissions to access it."

    echo -n "Building CLI binaries... "
    docker buildx build --rm -f docker/cli.dockerfile --output ../$VERSION --target cli . || fail "Error building CLI binaries."
    echo "done!"

    cd ..
}


# Builds the .tar.xz file packages with the HD configuration files
build_install_packages() {
    cd hyperdrive || fail "Directory ${PWD}/hyperdrive does not exist or you don't have permissions to access it."
    rm -f hyperdrive-install.tar.xz

    echo -n "Building Hyperdrive installer packages... "
    tar cfJ hyperdrive-install.tar.xz install || fail "Error building installer package."
    mv hyperdrive-install.tar.xz ../$VERSION
    echo "done!"

    cd ..
}


# Builds the Hyperdrive image and pushes it to Docker Hub
# NOTE: You must install qemu first; e.g. sudo apt-get install -y qemu qemu-user-static
build_daemon() {
    cd hyperdrive || fail "Directory ${PWD}/hyperdrive does not exist or you don't have permissions to access it."

    # Make a multiarch builder, ignore if it's already there
    docker buildx create --name multiarch-builder --driver docker-container --use > /dev/null 2>&1

    echo "Building and pushing Docker Hyperdrive image..."
    docker buildx build --rm --platform=linux/amd64,linux/arm64 -t nodeset/hyperdrive:$VERSION -f docker/daemon.dockerfile --push . || fail "Error building Hyperdrive daemon image."
    echo "done!"

    cd ..
}


# Tags the 'latest' Docker Hub image
tag_latest() {
    echo -n "Tagging 'latest' Docker image... "
    docker tag nodeset/hyperdrive:$VERSION nodeset/hyperdrive:latest
    echo "done!"

    echo -n "Pushing to Docker Hub... "
    docker push nodeset/hyperdrive:latest
    echo "done!"
}


# Print usage
usage() {
    echo "Usage: build-release.sh [options] -v <version number>"
    echo "This script assumes it is in a directory that contains subdirectories for all of the Hyperdrive repositories."
    echo "Options:"
    echo $'\t-a\tBuild all of the artifacts'
    echo $'\t-c\tBuild the CLI binaries for all platforms'
    echo $'\t-p\tBuild the Hyperdrive installer packages'
    echo $'\t-d\tBuild the Daemon Hyperdrive images, and push them to Docker Hub'
    echo $'\t-l\tTag the given version as "latest" on Docker Hub'
    exit 0
}


# =================
# === Main Body ===
# =================

# Parse arguments
while getopts "acpdlv:" FLAG; do
    case "$FLAG" in
        a) CLI=true PACKAGES=true DAEMON=true MANIFEST=true LATEST=true ;;
        c) CLI=true ;;
        p) PACKAGES=true ;;
        d) DAEMON=true ;;
        l) LATEST=true ;;
        v) VERSION="$OPTARG" ;;
        *) usage ;;
    esac
done
if [ -z "$VERSION" ]; then
    usage
fi

# Cleanup old artifacts
rm -f ./$VERSION/*
mkdir -p ./$VERSION

# Build the artifacts
if [ "$CLI" = true ]; then
    build_cli
fi
if [ "$PACKAGES" = true ]; then
    build_install_packages
fi
if [ "$DAEMON" = true ]; then
    build_daemon
fi
if [ "$LATEST" = true ]; then
    tag_latest
fi