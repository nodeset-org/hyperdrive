#!/bin/bash

# This script will build all of the artifacts involved in a new Hyperdrive release.

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
    echo -n "Building CLI binaries... "
    docker buildx build --rm -f docker/cli.dockerfile --output build/$VERSION --target cli . || fail "Error building CLI binaries."
    echo "done!"
}


# Builds the hyperdrive distro packages
build_distro_packages() {
    echo -n "Building deb packages..."
    docker buildx build --rm -f install/packages/debian/package.dockerfile --output build/$VERSION --target package . || fail "Error building deb packages."
    rm build/$VERSION/*.build build/$VERSION/*.buildinfo build/$VERSION/*.changes
    echo "done!"
}


# Builds the .tar.xz file packages with the HD configuration files
build_install_packages() {
    echo -n "Building Hyperdrive installer packages... "
    tar cfJ build/$VERSION/hyperdrive-install.tar.xz install/deploy || fail "Error building installer package."
    cp install/install.sh build/$VERSION
    echo "done!"
}


# Print usage
usage() {
    echo "Usage: build.sh [options] -v <version number>"
    echo "This script assumes it is in the hyperdrive repository directory."
    echo "Options:"
    echo $'\t-a\tBuild all of the artifacts'
    echo $'\t-c\tBuild the CLI binaries for all platforms'
    echo $'\t-t\tBuild the distro packages (.deb)'
    echo $'\t-p\tBuild the Hyperdrive installer packages'
    exit 0
}


# =================
# === Main Body ===
# =================

# Parse arguments
while getopts "actpv:" FLAG; do
    case "$FLAG" in
        a) CLI=true DISTRO=true PACKAGES=true ;;
        c) CLI=true ;;
        t) DISTRO=true ;;
        p) PACKAGES=true ;;
        v) VERSION="$OPTARG" ;;
        *) usage ;;
    esac
done
if [ -z "$VERSION" ]; then
    usage
fi

# Cleanup old artifacts
rm -rf build/$VERSION/*
mkdir -p build/$VERSION

# Make a multiarch builder, ignore if it's already there
docker buildx create --name multiarch-builder --driver docker-container --use > /dev/null 2>&1

# Build the artifacts
if [ "$CLI" = true ]; then
    build_cli
fi
if [ "$DISTRO" = true ]; then
    build_distro_packages
fi
if [ "$PACKAGES" = true ]; then
    build_install_packages
fi


# =======================
# === Manual Routines ===
# =======================

# Builds the deb package builder image
build_deb_builder() {
    cd hyperdrive || fail "Directory ${PWD}/hyperdrive does not exist or you don't have permissions to access it."

    echo -n "Building deb builder..."
    docker buildx build --rm --platform=linux/amd64,linux/arm64 -t nodeset/hyperdrive-deb-builder:$VERSION -f install/packages/debian/builder.dockerfile --push . || fail "Error building deb builder."
    echo "done!"

    cd ..
}
