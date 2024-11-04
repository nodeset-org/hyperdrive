#!/bin/sh

##
# Hyperdrive service installation script
# Prints progress messages to stdout
# All command output is redirected to stderr
##

COLOR_RED='\033[0;31m'
COLOR_YELLOW='\033[33m'
COLOR_RESET='\033[0m'

# Require root access for installation
#if [ "$(id -u)" -ne "0" ]; then
#     echo "This script requires root."
#     exit 1
#fi

# Print a failure message to stderr and exit
fail() {
    MESSAGE=$1
    >&2 echo -e "\n${COLOR_RED}**ERROR**\n$MESSAGE${COLOR_RESET}"
    exit 1
}

warn() {
    MESSAGE=$1
    >&2 echo -e "\n${COLOR_YELLOW}**WARNING**\n$MESSAGE${COLOR_RESET}"
}

# Get CPU architecture
UNAME_VAL=$(uname -m)
ARCH=""
case $UNAME_VAL in
    x86_64)  ARCH="amd64" ;;
    aarch64) ARCH="arm64" ;;
    arm64)   ARCH="arm64" ;;
    *)       fail "CPU architecture not supported: $UNAME_VAL" ;;
esac

# Get the platform type
PLATFORM=$(uname -s)
if [ "$PLATFORM" = "Linux" ]; then
    if command -v lsb_release &>/dev/null ; then
        PLATFORM=$(lsb_release -si)
    elif [ -f "/etc/centos-release" ]; then
        PLATFORM="CentOS"
    elif [ -f "/etc/fedora-release" ]; then
        PLATFORM="Fedora"
    fi
fi

##
# Utils
##

# Print progress
TOTAL_STEPS="7"
progress() {
    STEP_NUMBER=$1
    MESSAGE=$2
    echo "Step $STEP_NUMBER of $TOTAL_STEPS: $MESSAGE"
}

# Docker installation steps
add_user_docker() {
    usermod -aG docker $USER || fail "Could not add user to docker group."
}

# Install
install() {
    ##
    # Initialization
    ##

    # Parse arguments
    PACKAGE_VERSION="latest"
    while getopts "dl:i:r:v:b:h" FLAG; do
        case "$FLAG" in
            d) NO_DEPS=true ;;
            l) LOCAL_PACKAGE_PATH="$OPTARG" ;;
            i) HD_INSTALL_PATH="$OPTARG" ;;
            r) HD_RUNTIME_PATH="$OPTARG" ;;
            v) PACKAGE_VERSION="$OPTARG" ;;
            b) BASH_COMPLETION_PATH="$OPTARG" ;;
            *) fail "Incorrect usage." ;;
        esac
    done

    # Get package files URL
    PACKAGE_NAME="hyperdrive-install.tar.xz"
    if [ "$PACKAGE_VERSION" = "latest" ]; then
        PACKAGE_URL="https://github.com/nodeset-org/hyperdrive/releases/latest/download/$PACKAGE_NAME"
    else
        PACKAGE_URL="https://github.com/nodeset-org/hyperdrive/releases/download/$PACKAGE_VERSION/$PACKAGE_NAME"
    fi

    # Create temporary data folder; clean up on exit
    TEMPDIR=$(mktemp -d 2>/dev/null) || fail "Could not create temporary data directory."
    trap 'rm -rf "$TEMPDIR"' EXIT

    # Get temporary data paths
    PACKAGE_FILES_PATH="$TEMPDIR/install/deploy"
    AUTOCOMPLETE_FILES_PATH="$TEMPDIR/install/autocomplete"


    ##
    # Installation
    ##

    # OS dependencies
    if [ -z "$NO_DEPS" ]; then

    case "$PLATFORM" in

        # Ubuntu / Debian / Raspbian
        Ubuntu|Debian|Raspbian)

            # Get platform name
            PLATFORM_NAME=$(echo "$PLATFORM" | tr '[:upper:]' '[:lower:]')

            # Install OS dependencies
            progress 1 "Installing OS dependencies..."
            { apt-get -y update || fail "Could not update OS package definitions."; } >&2
            { apt-get -y install apt-transport-https ca-certificates curl gnupg gnupg-agent lsb-release software-properties-common chrony || fail "Could not install OS packages."; } >&2

            # Check for existing Docker installation
            progress 2 "Checking if Docker is installed..."
            dpkg-query -W -f='${Status}' docker-ce 2>&1 | grep -q -P '^install ok installed$' > /dev/null
            if [ $? != "0" ]; then
                echo "Installing Docker..."
                if [ ! -f /etc/apt/sources.list.d/docker.list ]; then
                    # Install the Docker repo
                    { mkdir -p /etc/apt/keyrings || fail "Could not create APT keyrings directory."; } >&2
                    { curl -fsSL "https://download.docker.com/linux/$PLATFORM_NAME/gpg" | gpg --dearmor -o /etc/apt/keyrings/docker.gpg || fail "Could not add docker repository key."; } >&2
                    { echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/$PLATFORM_NAME $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null || fail "Could not add docker repository."; } >&2
                fi
                { apt-get -y update || fail "Could not update OS package definitions."; } >&2
                { apt-get -y install docker-ce docker-ce-cli docker-compose-plugin containerd.io || fail "Could not install Docker packages."; } >&2
            fi

            # Check for existing docker-compose-plugin installation
            progress 2 "Checking if docker-compose-plugin is installed..."
            dpkg-query -W -f='${Status}' docker-compose-plugin 2>&1 | grep -q -P '^install ok installed$' > /dev/null
            if [ $? != "0" ]; then
                echo "Installing docker-compose-plugin..."
                if [ ! -f /etc/apt/sources.list.d/docker.list ]; then
                    # Install the Docker repo, removing the legacy one if it exists
                    { add-apt-repository --remove "deb [arch=$(dpkg --print-architecture)] https://download.docker.com/linux/$PLATFORM_NAME $(lsb_release -cs) stable"; } 2>/dev/null
                    { mkdir -p /etc/apt/keyrings || fail "Could not create APT keyrings directory."; } >&2
                    { curl -fsSL "https://download.docker.com/linux/$PLATFORM_NAME/gpg" | gpg --dearmor -o /etc/apt/keyrings/docker.gpg || fail "Could not add docker repository key."; } >&2
                    { echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/$PLATFORM_NAME $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null || fail "Could not add Docker repository."; } >&2
                fi
                { apt-get -y update || fail "Could not update OS package definitions."; } >&2
                { apt-get -y install docker-compose-plugin || fail "Could not install docker-compose-plugin."; } >&2
                { systemctl restart docker || fail "Could not restart docker daemon."; } >&2
            else
                echo "Already installed."
            fi

            # Add user to docker group
            progress 3 "Adding user to docker group..."
            >&2 add_user_docker

        ;;

        # Centos
        CentOS)

            # Install OS dependencies
            progress 1 "Installing OS dependencies..."
            { yum install -y yum-utils chrony || fail "Could not install OS packages."; } >&2
            { systemctl start chronyd || fail "Could not start chrony daemon."; } >&2

            # Install docker
            progress 2 "Installing docker..."
            { yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo || fail "Could not add docker repository."; } >&2
            { yum install -y docker-ce docker-ce-cli docker-compose-plugin containerd.io || fail "Could not install docker packages."; } >&2
            { systemctl start docker || fail "Could not start docker daemon."; } >&2
            { systemctl enable docker || fail "Could not set docker daemon to auto-start on boot."; } >&2

            # Check for existing docker-compose-plugin installation
            progress 2 "Checking if docker-compose-plugin is installed..."
            yum -q list installed docker-compose-plugin 2>/dev/null 1>/dev/null
            if [ $? != "0" ]; then
                echo "Installing docker-compose-plugin..."
                { yum install -y docker-compose-plugin || fail "Could not install docker-compose-plugin."; } >&2
                { systemctl restart docker || fail "Could not restart docker daemon."; } >&2
            else
                echo "Already installed."
            fi

            # Add user to docker group
            progress 3 "Adding user to docker group..."
            >&2 add_user_docker

        ;;

        # Fedora
        Fedora)

            # Install OS dependencies
            progress 1 "Installing OS dependencies..."
            { dnf -y install dnf-plugins-core chrony || fail "Could not install OS packages."; } >&2
            { systemctl start chronyd || fail "Could not start chrony daemon."; } >&2

            # Install docker
            progress 2 "Installing docker..."
            { dnf config-manager --add-repo https://download.docker.com/linux/fedora/docker-ce.repo || fail "Could not add docker repository."; } >&2
            { dnf -y install docker-ce docker-ce-cli docker-compose-plugin containerd.io || fail "Could not install docker packages."; } >&2
            { systemctl start docker || fail "Could not start docker daemon."; } >&2
            { systemctl enable docker || fail "Could not set docker daemon to auto-start on boot."; } >&2

            # Check for existing docker-compose-plugin installation
            progress 2 "Checking if docker-compose-plugin is installed..."
            dnf -q list installed docker-compose-plugin 2>/dev/null 1>/dev/null
            if [ $? != "0" ]; then
                echo "Installing docker-compose-plugin..."
                { dnf install -y docker-compose-plugin || fail "Could not install docker-compose-plugin."; } >&2
                { systemctl restart docker || fail "Could not restart docker daemon."; } >&2
            else
                echo "Already installed."
            fi

            # Add user to docker group
            progress 3 "Adding user to docker group..."
            >&2 add_user_docker

        ;;

        # Unsupported OS
        *)
            RED='\033[0;31m'
            echo ""
            echo -e "${RED}**ERROR**"
            echo "Automatic dependency installation for the $PLATFORM operating system is not supported."
            echo "Please install docker and docker-compose-plugin manually, then try again with the '-d' flag to skip OS dependency installation."
            echo "Be sure to add yourself to the docker group with 'usermod -aG docker $USER' after installing docker."
            echo "Log out and back in, or restart your system after you run this command."
            echo -e "${RESET}"
            exit 1
        ;;

    esac
    else
        echo "Skipping steps 1 - 3 (OS dependencies & docker)"
    fi

    # Create hyperdrive dir & files - default to Linux paths for now if not set
    if [ -z "$HD_INSTALL_PATH" ]; then
        HD_INSTALL_PATH="/usr/share/hyperdrive"
    fi
    if [ -z "$HD_RUNTIME_PATH" ]; then
        HD_RUNTIME_PATH="/var/lib/hyperdrive"
    fi
    if [ -z "$BASH_COMPLETION_PATH" ]; then
        BASH_COMPLETION_PATH="/usr/share/bash-completion/completions"
    fi

    progress 4 "Creating Hyperdrive directory structure..."
    { mkdir -p "$HD_INSTALL_PATH" || fail "Could not create the Hyperdrive resources directory."; } >&2
    { mkdir -p "$HD_RUNTIME_PATH/data" || fail "Could not create the Hyperdrive system data directory."; } >&2
    { mkdir -p "$HD_RUNTIME_PATH/global" || fail "Could not create the Hyperdrive system global directory."; } >&2
    { chmod 0700 "$HD_RUNTIME_PATH/data" || fail "Could not set the Hyperdrive data directory permissions."; } >&2
    { chmod 0700 "$HD_RUNTIME_PATH/global" || fail "Could not set the Hyperdrive global directory permissions."; } >&2

    # Download and extract package files
    progress 5 "Downloading Hyperdrive package files..."
    if [ -z "$LOCAL_PACKAGE_PATH" ]; then
        { curl -L "$PACKAGE_URL" | tar -xJ -C "$TEMPDIR" || fail "Could not download and extract the Hyperdrive package files."; } >&2
    else
        if [ ! -f $LOCAL_PACKAGE_PATH ]; then
            fail "Installer package [$LOCAL_PACKAGE_PATH] does not exist." >&2
        fi
        { tar -f "$LOCAL_PACKAGE_PATH" -xJ -C "$TEMPDIR" || fail "Could not extract the local Hyperdrive package files."; } >&2
    fi
    { test -d "$PACKAGE_FILES_PATH" || fail "Could not extract the Hyperdrive package files."; } >&2

    # Copy package files
    progress 6 "Copying package files to Hyperdrive system directory..."
    { find "$PACKAGE_FILES_PATH" -exec cp -r {} "$HD_INSTALL_PATH" \; || fail "Could not copy deployment artifacts ($PACKAGE_FILES_PATH) to the Hyperdrive system directory ($HD_INSTALL_PATH)."; } >&2
    { find "$HD_INSTALL_PATH/scripts" -name "*.sh" -exec chmod +x {} \; 2>/dev/null || fail "Could not set executable permissions on package files."; } >&2

    # Copy bash completion helper
    if [ -d "$BASH_COMPLETION_PATH" ]; then
        { cp "$AUTOCOMPLETE_FILES_PATH/bash_autocomplete" "$BASH_COMPLETION_PATH/hyperdrive" 2>/dev/null || fail "Could not install bash completion file."; } >&2
    else
        warn "No directory at expected path for bash_completion '$BASH_COMPLETION_PATH' - skipping."
    fi

    # Clean up unnecessary files from old installations
    progress 7 "Cleaning up obsolete files from previous installs..."
}

install "$@"

