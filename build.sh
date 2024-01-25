#!/bin/bash

# This script will build all of the artifacts involved in a new Hyperdrive release.
# NOTE: You MUST put this in a directory that has the `hyperdrive` repository cloned as a subdirectory.

PROJECT_NAME="hyperdrive"

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
    cd $PROJECT_NAME || fail "Directory ${PWD}/$PROJECT_NAME does not exist or you don't have permissions to access it."

    echo -n "Building CLI binaries... "
    docker buildx build --rm -f docker/cli.dockerfile --output ../$VERSION --target cli . || fail "Error building CLI binaries."
    rm -rf hyperdrive-cli/build
    echo "done!"

    cd ..
}


# Builds the .tar.xz file packages with the RP configuration files
build_install_packages() {
    cd $PROJECT_NAME || fail "Directory ${PWD}/$PROJECT_NAME does not exist or you don't have permissions to access it."
    rm -f hyperdrive-install.tar.xz

    echo -n "Building Hyperdrive installer packages... "
    tar cfJ hyperdrive-install.tar.xz install || fail "Error building installer package."
    #tar cfJ hyperdrive-install.tar.xz install/deploy || fail "Error building installer package."
    mv hyperdrive-install.tar.xz ../$VERSION
    #cp install/install.sh ../$VERSION
    echo "done!"

    cd ..
}


# Builds the daemon binaries and Docker Hyperdrive images, and pushes them to Docker Hub
# NOTE: You must install qemu first; e.g. sudo apt-get install -y qemu qemu-user-static
build_daemon() {
    cd $PROJECT_NAME || fail "Directory ${PWD}/$PROJECT_NAME does not exist or you don't have permissions to access it."

    echo "Building Docker Hyperdrive image..."
    docker buildx build --platform=linux/amd64 -t nodeset/$PROJECT_NAME:$VERSION-amd64 -f docker/daemon.dockerfile --load . || fail "Error building amd64 Docker Hyperdrive image."
    docker buildx build --platform=linux/arm64 -t nodeset/$PROJECT_NAME:$VERSION-arm64 -f docker/daemon.dockerfile --load . || fail "Error building arm64 Docker Hyperdrive image."
    echo "done!"

    echo -n "Pushing to Docker Hub... "
    docker push nodeset/$PROJECT_NAME:$VERSION-amd64 || fail "Error pushing amd64 Docker Hyperdrive image to Docker Hub."
    docker push nodeset/$PROJECT_NAME:$VERSION-arm64 || fail "Error pushing arm Docker Hyperdrive image to Docker Hub."
    echo "done!"

    cd ..
}


# Builds the Docker Manifests and pushes them to Docker Hub
build_docker_manifest() {
    echo -n "Building Docker manifest... "
    rm -f ~/.docker/manifests/docker.io_nodeset_$PROJECT_NAME-$VERSION
    docker manifest create nodeset/$PROJECT_NAME:$VERSION --amend nodeset/$PROJECT_NAME:$VERSION-amd64 --amend nodeset/$PROJECT_NAME:$VERSION-arm64
    echo "done!"

    echo -n "Pushing to Docker Hub... "
    docker manifest push --purge nodeset/$PROJECT_NAME:$VERSION
    echo "done!"
}


# Builds the 'latest' Docker Manifests and pushes them to Docker Hub
build_latest_docker_manifest() {
    echo -n "Building 'latest' Docker manifest... "
    rm -f ~/.docker/manifests/docker.io_nodeset_$PROJECT_NAME-latest
    docker manifest create nodeset/$PROJECT_NAME:latest --amend nodeset/$PROJECT_NAME:$VERSION-amd64 --amend nodeset/$PROJECT_NAME:$VERSION-arm64
    echo "done!"

    echo -n "Pushing to Docker Hub... "
    docker manifest push --purge nodeset/$PROJECT_NAME:latest
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
    echo $'\t-n\tBuild the Daemon manifests, and push them to Docker Hub'
    exit 0
}


# =================
# === Main Body ===
# =================

# Parse arguments
while getopts "acpdnlrfv:" FLAG; do
    case "$FLAG" in
        a) CLI=true PACKAGES=true DAEMON=true MANIFEST=true LATEST_MANIFEST=true ;;
        c) CLI=true ;;
        p) PACKAGES=true ;;
        d) DAEMON=true ;;
        n) MANIFEST=true ;;
        l) LATEST_MANIFEST=true ;;
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
if [ "$MANIFEST" = true ]; then
    build_docker_manifest
fi
if [ "$LATEST_MANIFEST" = true ]; then
    build_latest_docker_manifest
fi