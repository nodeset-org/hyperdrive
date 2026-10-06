#!/bin/sh
# This script launches PBS clients for Hyperdrive's docker stack; only edit if you know what you're doing ;)


# Commit-Boost
if [ "$CLIENT" = "commitBoost" ]; then
    # Commit-Boost PBS takes its settings from the config file named by CB_CONFIG,
    # not from command-line flags.
    export CB_CONFIG="/cb_config.toml"
    if [ ! -f "$CB_CONFIG" ]; then
        echo "Commit-Boost config file not found at $CB_CONFIG. Please ensure the config file is correctly mounted."
        exit 1
    fi
    # v0.10.0 unified the CLI, PBS, and signer into `commit-boost` with subcommands
    # (Commit-Boost/commit-boost-client#425). v0.9.x images still ship the
    # standalone commit-boost-pbs binary, which is what the default tag uses.
    if [ -x /usr/local/bin/commit-boost ]; then
        exec /usr/local/bin/commit-boost pbs
    elif [ -x /usr/local/bin/commit-boost-pbs ]; then
        exec /usr/local/bin/commit-boost-pbs
    else
        echo "Commit-Boost PBS binary not found (tried /usr/local/bin/commit-boost and /usr/local/bin/commit-boost-pbs)."
        exit 127
    fi
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
