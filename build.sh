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
    tar cfJ build/$VERSION/hyperdrive-install.tar.xz install || fail "Error building installer package."
    cp install/install.sh build/$VERSION
    echo "done!"
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
    else
        docker buildx build --rm --load --build-arg BINARIES_PATH=build/$VERSION -t nodeset/hyperdrive:$VERSION -f docker/daemon.dockerfile . || fail "Error building Hyperdrive Docker image."
    fi
    echo "done!"
}


# Builds the Stakewise daemon image and pushes it to Docker Hub
# NOTE: You must install qemu first; e.g. sudo apt-get install -y qemu qemu-user-static
build_sw_daemon() {
    echo "Building Stakewise daemon binaries..."
    docker buildx build --rm --platform=linux/amd64,linux/arm64 -f docker/modules/stakewise/sw_daemon-build.dockerfile --output build/$VERSION --target daemon . || fail "Error building Stakewise daemon binaries."
    echo "done!"

    # Flatted the folders to make it easier to upload artifacts to github
    mv build/$VERSION/linux_amd64/hyperdrive-stakewise-daemon build/$VERSION/hyperdrive-stakewise-daemon-linux-amd64
    mv build/$VERSION/linux_arm64/hyperdrive-stakewise-daemon build/$VERSION/hyperdrive-stakewise-daemon-linux-arm64

    # Clean up the empty directories
    rmdir build/$VERSION/linux_amd64 build/$VERSION/linux_arm64

    echo "Building Stakewise Docker image..."
    # If uploading, make and push a manifest
    if [ "$UPLOAD" = true ]; then
        docker buildx build --rm --platform=linux/amd64,linux/arm64 --build-arg BINARIES_PATH=build/$VERSION -t nodeset/hyperdrive-stakewise:$VERSION -f docker/modules/stakewise/sw_daemon.dockerfile --push . || fail "Error building Stakewise Docker image."
    else
        docker buildx build --rm --load --build-arg BINARIES_PATH=build/$VERSION -t nodeset/hyperdrive-stakewise:$VERSION -f docker/modules/stakewise/sw_daemon.dockerfile . || fail "Error building Stakewise Docker image."
    fi
    echo "done!"
}

# Builds the Constellation daemon image and pushes it to Docker Hub
# NOTE: You must install qemu first; e.g. sudo apt-get install -y qemu qemu-user-static
build_constellation_daemon() {
    echo "Building Constellation daemon binaries..."
    docker buildx build --rm --platform=linux/amd64,linux/arm64 -f docker/modules/constellation/const_daemon-build.dockerfile --output ../$VERSION --target daemon . || fail "Error building Constellation daemon binaries."
    echo "done!"


    # Flatted the folders to make it easier to upload artifacts to github
    mv build/$VERSION/linux_amd64/hyperdrive-constellation-daemon build/$VERSION/hyperdrive-constellation-daemon-linux-amd64
    mv build/$VERSION/linux_arm64/hyperdrive-constellation-daemon build/$VERSION/hyperdrive-constellation-daemon-linux-arm64

    # Clean up the empty directories
    rmdir build/$VERSION/linux_amd64 build/$VERSION/linux_arm64

    echo "Building Constellation Docker image..."
    # If uploading, make and push a manifest
    if [ "$UPLOAD" = true ]; then
        docker buildx build --rm --platform=linux/amd64,linux/arm64 --build-arg BINARIES_PATH=build/$VERSION -t nodeset/hyperdrive-constellation:$VERSION -f docker/modules/constellation/const_daemon.dockerfile --push . || fail "Error building Constellation Docker image."
    else
        docker buildx build --rm --load --build-arg BINARIES_PATH=build/$VERSION -t nodeset/hyperdrive-constellation:$VERSION -f docker/modules/constellation/const_daemon.dockerfile . || fail "Error building Constellation Docker image."
    fi
    echo "done!"

    # # Copy the daemon binaries to a build folder so the image can access them
    # mkdir -p ./build
    # cp ../$VERSION/linux_amd64/* ./build
    # cp ../$VERSION/linux_arm64/* ./build
    # echo "done!"

    # echo "Building Constellation Docker image..."
    # docker buildx build --rm --platform=linux/amd64,linux/arm64 -t nodeset/hyperdrive-constellation:$VERSION -f docker/modules/constellation/const_daemon.dockerfile --push . || fail "Error building Constellation Docker image."
    # echo "done!"

    # # Cleanup
    # mv ../$VERSION/linux_amd64/* ../$VERSION
    # mv ../$VERSION/linux_arm64/* ../$VERSION
    # rm -rf ../$VERSION/linux_amd64/
    # rm -rf ../$VERSION/linux_arm64/
    # rm -rf ./build

    # cd ..
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
    echo "This script assumes it is in a directory that contains subdirectories for all of the Hyperdrive repositories."
    echo "Options:"
    echo $'\t-a\tBuild all of the artifacts'
    echo $'\t-c\tBuild the CLI binaries for all platforms'
    echo $'\t-t\tBuild the distro packages (.deb)'
    echo $'\t-p\tBuild the Hyperdrive installer packages'
    echo $'\t-d\tBuild the Hyperdrive daemon image, and push it to Docker Hub'
    echo $'\t-s\tBuild the Hyperdrive Stakewise daemon image, and push it to Docker Hub'
    echo $'\t-x\tBuild the Hyperdrive Constellation daemon image, and push it to Docker Hub'
    echo $'\t-l\tTag the given version as "latest" on Docker Hub'
    echo $'\t-u\tWhen passed with a build, upload the resulting image tags to Docker Hub'
    exit 0
}


# =================
# === Main Body ===
# =================

# Parse arguments
while getopts "actpdsxluv:" FLAG; do
    case "$FLAG" in
        a) CLI=true DISTRO=true PACKAGES=true DAEMON=true SW_DAEMON=true CONST_DAEMON=true;;
        c) CLI=true ;;
        d) DAEMON=true ;;
        l) LATEST=true ;;
        x) CONST_DAEMON=true ;;
        p) PACKAGES=true ;;
        s) SW_DAEMON=true ;;
        t) DISTRO=true ;;
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
if [ "$DAEMON" = true ]; then
    build_daemon
fi
if [ "$SW_DAEMON" = true ]; then
    build_sw_daemon
fi
if [ "$CONST_DAEMON" = true ]; then
    build_constellation_daemon
fi
if [ "$LATEST" = true ]; then
    tag_latest
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
