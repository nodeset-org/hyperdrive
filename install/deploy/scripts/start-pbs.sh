#!/bin/sh
# This script launches PBS clients for Hyperdrive's docker stack; only edit if you know what you're doing ;)


# Commit-Boost
if [ "$CLIENT" = "commitBoost" ]; then
    # Common-Boost PBS doesn't accept command line argumnets, everything must be in the config file which is
    # specified by an environment variable.
    export CB_CONFIG="/cb_config.toml"
    if [ ! -f "$CB_CONFIG" ]; then
        echo "Commit-Boost config file not found at $CB_CONFIG. Please ensure the config file is correctly mounted."
        exit 1
    fi
    CMD="/usr/local/bin/commit-boost-pbs"
    exec ${CMD}
fi


# MEV-Boost
if [ "$CLIENT" = "mevBoost" ]; then
    CMD="/app/mev-boost \
        -$ETH_NETWORK \
        -addr 0.0.0.0:$PBS_PORT \
        -relay-check \
        -relays $PBS_RELAYS"
    
    exec ${CMD}
fi