#!/bin/bash

# This script will build the Hyperdrive daemon.

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


# Builds the Hyperdrive image and pushes it to Docker Hub
# NOTE: You must install qemu first; e.g. sudo apt-get install -y qemu qemu-user-static
build_daemon() {
    echo "Building Hyperdrive binaries..."
    docker buildx build --rm --platform=linux/amd64,linux/arm64 -f docker/daemon-build.dockerfile --output build/$VERSION --target daemon . || fail "Error building Hyperdrive daemon binaries."
    echo "done!"

    # Flatted the folders to make it easier to upload artifacts to github
    mv build/$VERSION/linux_amd64/hyperdrive-daemon build/$VERSION/hyperdrive-daemon-linux-amd64
    mv build/$VERSION/linux_arm64/hyperdrive-daemon build/$VERSION/hyperdrive-daemon-linux-arm64

    # Clean up the empty directories
    rmdir build/$VERSION/linux_amd64 build/$VERSION/linux_arm64

    echo "Building Hyperdrive Docker image..."
    # If uploading, make and push a manifest
    if [ "$UPLOAD" = true ]; then
        docker buildx build --rm --platform=linux/amd64,linux/arm64 --build-arg BINARIES_PATH=build/$VERSION -t nodeset/hyperdrive:$VERSION -f docker/daemon.dockerfile --push . || fail "Error building Hyperdrive Docker image."
    elif [ "$LOCAL_UPLOAD" = true ]; then
        if [ -z "$LOCAL_DOCKER_REGISTRY" ]; then
            fail "LOCAL_DOCKER_REGISTRY must be set to upload to a local registry."
        fi
        docker buildx build --rm --platform=linux/amd64,linux/arm64 --build-arg BINARIES_PATH=build/$VERSION -t $LOCAL_DOCKER_REGISTRY/nodeset/hyperdrive:$VERSION -f docker/daemon.dockerfile --push . || fail "Error building Hyperdrive Docker image."
    else
        docker buildx build --rm --load --build-arg BINARIES_PATH=build/$VERSION -t nodeset/hyperdrive:$VERSION -f docker/daemon.dockerfile . || fail "Error building Hyperdrive Docker image."
    fi
    echo "done!"
}


# Tags the 'latest' Docker Hub image
tag_latest() {
    echo -n "Tagging 'latest' Docker image... "
    docker tag nodeset/hyperdrive:$VERSION nodeset/hyperdrive:latest
    echo "done!"

    if [ "$UPLOAD" = true ]; then
        echo -n "Pushing to Docker Hub... "
        docker push nodeset/hyperdrive:latest
        echo "done!"
    else
        echo "The image tag only exists locally. Rerun with -u to upload to Docker Hub."
    fi
}


# Print usage
usage() {
    echo "Usage: build.sh [options] -v <version number>"
    echo "This script assumes it is in the hyperdrive-daemon repository directory."
    echo "Options:"
    echo $'\t-a\tBuild all of the artifacts'
    echo $'\t-c\tWhen passed with a build, upload the resulting image tags to a local Docker registry specified in $LOCAL_DOCKER_REGISTRY'
    echo $'\t-d\tBuild the Hyperdrive daemon image and Docker container'
    echo $'\t-l\tTag the given version as "latest" on Docker Hub'
    echo $'\t-u\tWhen passed with a build, upload the resulting image tags to Docker Hub'
    exit 0
}


# =================
# === Main Body ===
# =================

# Parse arguments
while getopts "acdluv:" FLAG; do
    case "$FLAG" in
        a) DAEMON=true ;;
        c) LOCAL_UPLOAD=true ;;
        d) DAEMON=true ;;
        l) LATEST=true ;;
        u) UPLOAD=true ;;
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
# NOTE: if using a local repo with a private CA, you will have to follow these steps to add the CA to the builder:
# https://stackoverflow.com/a/73585243

# Build the artifacts
if [ "$DAEMON" = true ]; then
    build_daemon
fi
if [ "$LATEST" = true ]; then
    tag_latest
fi
