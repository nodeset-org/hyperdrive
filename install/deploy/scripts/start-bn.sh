#!/bin/sh
# This script launches beacon clients for Hyperdrive's docker stack; only edit if you know what you're doing ;)

# Performance tuning for ARM systems
UNAME_VAL=$(uname -m)
if [ "$UNAME_VAL" = "arm64" ] || [ "$UNAME_VAL" = "aarch64" ]; then
    # Get the number of available cores
    CORE_COUNT=$(nproc)

    # Don't do performance tweaks on systems with 6+ cores
    if [ "$CORE_COUNT" -gt "5" ]; then
        echo "$CORE_COUNT cores detected, skipping performance tuning"
    else
        echo "$CORE_COUNT cores detected, activating performance tuning"
        PERF_PREFIX="ionice -c 2 -n 0"
        echo "Performance tuning: $PERF_PREFIX"
    fi
fi

# Check for the JWT auth file
if [ ! -f "/secrets/jwtsecret" ]; then
    echo "JWT secret file not found, please try again when the Execution Client has created one."
    exit 1
fi

# Lighthouse startup
if [ "$CLIENT" = "lighthouse" ]; then

    CMD="$PERF_PREFIX /usr/local/bin/lighthouse beacon \
        --network $ETH_NETWORK \
        --datadir /ethclient/lighthouse \
        --port $BN_P2P_PORT \
        --discovery-port $BN_P2P_PORT \
        --execution-endpoint $EC_ENGINE_ENDPOINT \
        --http \
        --http-address 0.0.0.0 \
        --http-port ${BN_API_PORT:-5052} \
        --eth1-blocks-per-log-query 150 \
        --disable-upnp \
        --staking \
        --execution-jwt=/secrets/jwtsecret \
        --quic-port ${BN_P2P_QUIC_PORT:-8001} \
        --historic-state-cache-size 2 \
        $BN_ADDITIONAL_FLAGS"

    # Performance tuning for ARM systems
    UNAME_VAL=$(uname -m)
    if [ "$UNAME_VAL" = "arm64" ] || [ "$UNAME_VAL" = "aarch64" ]; then
        CMD="$CMD --execution-timeout-multiplier 2"
    fi

    if [ ! -z "$MEV_BOOST_URL" ]; then
        CMD="$CMD --builder $MEV_BOOST_URL"
    fi

    if [ ! -z "$BN_MAX_PEERS" ]; then
        CMD="$CMD --target-peers $BN_MAX_PEERS"
    fi

    if [ "$ENABLE_METRICS" = "true" ]; then
        CMD="$CMD --metrics --metrics-address 0.0.0.0 --metrics-port $BN_METRICS_PORT --validator-monitor-auto"
    fi

    if [ ! -z "$CHECKPOINT_SYNC_URL" ]; then
        CMD="$CMD --checkpoint-sync-url $CHECKPOINT_SYNC_URL"
    fi

    if [ "$ENABLE_BITFLY_NODE_METRICS" = "true" ]; then
        CMD="$CMD --monitoring-endpoint $BITFLY_NODE_METRICS_ENDPOINT?apikey=$BITFLY_NODE_METRICS_SECRET&machine=$BITFLY_NODE_METRICS_MACHINE_NAME"
    fi

    exec ${CMD}

fi

# Lodestar startup
if [ "$CLIENT" = "lodestar" ]; then

    CMD="$PERF_PREFIX /usr/local/bin/node --max-http-header-size=65536 /usr/app/packages/cli/bin/lodestar beacon \
        --network $ETH_NETWORK \
        --dataDir /ethclient/lodestar \
        --port $BN_P2P_PORT \
        --execution.urls $EC_ENGINE_ENDPOINT \
        --rest \
        --rest.address 0.0.0.0 \
        --rest.port ${BN_API_PORT:-5052} \
        --jwt-secret /secrets/jwtsecret \
        $BN_ADDITIONAL_FLAGS"

    if [ ! -z "$TTD_OVERRIDE" ]; then
        CMD="$CMD --terminal-total-difficulty-override $TTD_OVERRIDE"
    fi

    if [ ! -z "$MEV_BOOST_URL" ]; then
        CMD="$CMD --builder --builder.urls $MEV_BOOST_URL"
    fi

    if [ ! -z "$BN_MAX_PEERS" ]; then
        CMD="$CMD --targetPeers $BN_MAX_PEERS"
    fi

    if [ "$ENABLE_METRICS" = "true" ]; then
        CMD="$CMD --metrics --metrics.address 0.0.0.0 --metrics.port $BN_METRICS_PORT"
    fi

    if [ ! -z "$EXTERNAL_IP" ]; then
        CMD="$CMD --enr.ip $EXTERNAL_IP --nat"
    fi

    if [ ! -z "$CHECKPOINT_SYNC_URL" ]; then
        CMD="$CMD --checkpointSyncUrl $CHECKPOINT_SYNC_URL"
    fi

    if [ "$ENABLE_BITFLY_NODE_METRICS" = "true" ]; then
        CMD="$CMD --monitoring.endpoint $BITFLY_NODE_METRICS_ENDPOINT?apikey=$BITFLY_NODE_METRICS_SECRET&machine=$BITFLY_NODE_METRICS_MACHINE_NAME"
    fi

    exec ${CMD}

fi

# Nimbus startup
if [ "$CLIENT" = "nimbus" ]; then

    # Handle checkpoint syncing
    if [ ! -z "$CHECKPOINT_SYNC_URL" ]; then
        # Ignore it if a DB already exists
        if [ -f "/ethclient/nimbus/db/nbc.sqlite3" ]; then
            echo "Nimbus database already exists, ignoring checkpoint sync."
        else 
            echo "Starting checkpoint sync for Nimbus..."
            $PERF_PREFIX /home/user/nimbus-eth2/build/nimbus_beacon_node trustedNodeSync --network=$ETH_NETWORK --data-dir=/ethclient/nimbus --trusted-node-url=$CHECKPOINT_SYNC_URL --backfill=false
            echo "Checkpoint sync complete!"
        fi
    fi

    CMD="$PERF_PREFIX /home/user/nimbus-eth2/build/nimbus_beacon_node \
        --non-interactive \
        --enr-auto-update \
        --network=$ETH_NETWORK \
        --data-dir=/ethclient/nimbus \
        --tcp-port=$BN_P2P_PORT \
        --udp-port=$BN_P2P_PORT \
        --web3-url=$EC_ENGINE_ENDPOINT \
        --rest \
        --rest-address=0.0.0.0 \
        --rest-port=${BN_API_PORT:-5052} \
        --jwt-secret=/secrets/jwtsecret \
        $BN_ADDITIONAL_FLAGS"

    if [ ! -z "$MEV_BOOST_URL" ]; then
        CMD="$CMD --payload-builder --payload-builder-url=$MEV_BOOST_URL"
    fi

    if [ ! -z "$BN_MAX_PEERS" ]; then
        CMD="$CMD --max-peers=$BN_MAX_PEERS"
    fi

    if [ "$ENABLE_METRICS" = "true" ]; then
        CMD="$CMD --metrics --metrics-address=0.0.0.0 --metrics-port=$BN_METRICS_PORT"
    fi

    if [ ! -z "$EXTERNAL_IP" ]; then
        CMD="$CMD --nat=extip:$EXTERNAL_IP"
    fi

    if [ ! -z "$NIMBUS_PRUNING_MODE" ]; then
        CMD="$CMD --history=$NIMBUS_PRUNING_MODE"
    fi

    exec ${CMD}

fi

# Prysm startup
if [ "$CLIENT" = "prysm" ]; then

    # Grab the Hoodi genesis state if needed
    if [ "$ETH_NETWORK" = "hoodi" ]; then
        echo "Prysm is configured to use Hoodi, genesis state required."
        if [ ! -f "/ethclient/hoodi-genesis.ssz" ]; then
            echo "Downloading from Github..."
            wget -q https://github.com/eth-clients/hoodi/raw/refs/heads/main/metadata/genesis.ssz -O /ethclient/hoodi-genesis.ssz
            echo "Download complete."
        else
            echo "Genesis state already downloaded, continuing."
        fi
    fi

    CMD="$PERF_PREFIX /app/cmd/beacon-chain/beacon-chain \
        --accept-terms-of-use \
        --$ETH_NETWORK \
        --datadir /ethclient/prysm \
        --p2p-tcp-port $BN_P2P_PORT \
        --p2p-udp-port $BN_P2P_PORT \
        --execution-endpoint $EC_ENGINE_ENDPOINT \
        --rpc-host 0.0.0.0 \
        --rpc-port ${BN_RPC_PORT:-5053} \
        --grpc-gateway-host 0.0.0.0 \
        --grpc-gateway-port ${BN_API_PORT:-5052} \
        --p2p-quic-port ${BN_P2P_QUIC_PORT:-8001} \
        --eth1-header-req-limit 150 \
        --jwt-secret=/secrets/jwtsecret \
        --enable-experimental-backfill \
        $BN_ADDITIONAL_FLAGS"

    if [ ! -z "$MEV_BOOST_URL" ]; then
        CMD="$CMD --http-mev-relay $MEV_BOOST_URL"
    fi

    if [ ! -z "$BN_MAX_PEERS" ]; then
        CMD="$CMD --p2p-max-peers $BN_MAX_PEERS"
    fi

    if [ "$ENABLE_METRICS" = "true" ]; then
        CMD="$CMD --monitoring-host 0.0.0.0 --monitoring-port $BN_METRICS_PORT"
    else
        CMD="$CMD --disable-monitoring"
    fi

    if [ "$ETH_NETWORK" = "hoodi" ]; then
        CMD="$CMD --genesis-state /ethclient/hoodi-genesis.ssz"
    fi

    if [ ! -z "$CHECKPOINT_SYNC_URL" ]; then
        CMD="$CMD --checkpoint-sync-url=$CHECKPOINT_SYNC_URL --genesis-beacon-api-url=$CHECKPOINT_SYNC_URL"
    fi

    exec ${CMD}

fi

# Teku startup
if [ "$CLIENT" = "teku" ]; then

    CMD="$PERF_PREFIX /opt/teku/bin/teku \
        --network=$ETH_NETWORK \
        --data-path=/ethclient/teku \
        --p2p-port=$BN_P2P_PORT \
        --ee-endpoint=$EC_ENGINE_ENDPOINT \
        --rest-api-enabled \
        --rest-api-interface=0.0.0.0 \
        --rest-api-port=${BN_API_PORT:-5052} \
        --rest-api-host-allowlist=* \
        --eth1-deposit-contract-max-request-size=150 \
        --log-destination=CONSOLE \
        --ee-jwt-secret-file=/secrets/jwtsecret \
        --beacon-liveness-tracking-enabled \
        --validators-graffiti-client-append-format=DISABLED \
        $BN_ADDITIONAL_FLAGS"

    if [ "$TEKU_ARCHIVE_MODE" = "true" ]; then
        CMD="$CMD --data-storage-mode=archive"
    fi

    if [ ! -z "$MEV_BOOST_URL" ]; then
        CMD="$CMD --builder-endpoint=$MEV_BOOST_URL"
    fi

    if [ ! -z "$BN_MAX_PEERS" ]; then
        CMD="$CMD --p2p-peer-lower-bound=$BN_MAX_PEERS --p2p-peer-upper-bound=$BN_MAX_PEERS"
    fi

    if [ "$ENABLE_METRICS" = "true" ]; then
        CMD="$CMD --metrics-enabled=true --metrics-interface=0.0.0.0 --metrics-port=$BN_METRICS_PORT --metrics-host-allowlist=*"
    fi

    if [ ! -z "$EXTERNAL_IP" ]; then
        CMD="$CMD --p2p-advertised-ip=$EXTERNAL_IP"
    fi

    if [ ! -z "$CHECKPOINT_SYNC_URL" ]; then
        CMD="$CMD --checkpoint-sync-url=$CHECKPOINT_SYNC_URL"
    fi

    if [ "$ENABLE_BITFLY_NODE_METRICS" = "true" ]; then
        CMD="$CMD --metrics-publish-endpoint=$BITFLY_NODE_METRICS_ENDPOINT?apikey=$BITFLY_NODE_METRICS_SECRET&machine=$BITFLY_NODE_METRICS_MACHINE_NAME"
    fi

    if [ "$TEKU_JVM_HEAP_SIZE" -gt "0" ]; then
        CMD="env JAVA_OPTS=\"-Xmx${TEKU_JVM_HEAP_SIZE}m\" $CMD"
    fi

    exec ${CMD}

fi
