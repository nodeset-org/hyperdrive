#!/bin/sh
# This script launches validator clients for Hyperdrive's docker stack; only edit if you know what you're doing ;)

# Lighthouse startup
if [ "$CLIENT" = "lighthouse" ]; then

    # Set up the CC + fallback string
    BN_URL_STRING=$BN_API_ENDPOINT
    if [ ! -z "$FALLBACK_BN_API_ENDPOINT" ]; then
        BN_URL_STRING="$BN_API_ENDPOINT,$FALLBACK_BN_API_ENDPOINT"
    fi

    CMD="/usr/local/bin/lighthouse validator \
        --network $ETH_NETWORK \
        --datadir /validators/lighthouse \
        --init-slashing-protection \
        --logfile-max-number 0 \
        --beacon-nodes $BN_URL_STRING \
        --suggested-fee-recipient $FEE_RECIPIENT \
        $VC_ADDITIONAL_FLAGS"

    if [ "$DOPPELGANGER_DETECTION" = "true" ]; then
        CMD="$CMD --enable-doppelganger-protection"
    fi

    if [ "$ENABLE_MEV_BOOST" = "true" ]; then
        CMD="$CMD --builder-proposals --prefer-builder-proposals"
    fi

    if [ "$ENABLE_METRICS" = "true" ]; then
        CMD="$CMD --metrics --metrics-address 0.0.0.0 --metrics-port $VC_METRICS_PORT"
    fi

    if [ "$ENABLE_BITFLY_NODE_METRICS" = "true" ]; then
        CMD="$CMD --monitoring-endpoint $BITFLY_NODE_METRICS_ENDPOINT?apikey=$BITFLY_NODE_METRICS_SECRET&machine=$BITFLY_NODE_METRICS_MACHINE_NAME"
    fi

    exec ${CMD} --graffiti "$GRAFFITI"

fi

# Lodestar startup
if [ "$CLIENT" = "lodestar" ]; then

    # Remove any lock files that were left over accidentally after an unclean shutdown
    find /validators/lodestar/validators -name voting-keystore.json.lock -delete

    # Set up the CC + fallback string
    BN_URL_STRING=$BN_API_ENDPOINT
    if [ ! -z "$FALLBACK_BN_API_ENDPOINT" ]; then
        BN_URL_STRING="$BN_API_ENDPOINT,$FALLBACK_BN_API_ENDPOINT"
    fi

    CMD="/usr/app/node_modules/.bin/lodestar validator \
        --network $ETH_NETWORK \
        --dataDir /validators/lodestar \
        --beacon-nodes $BN_URL_STRING \
        $FALLBACK_BN_STRING \
        --keystoresDir /validators/lodestar/validators \
        --secretsDir /validators/lodestar/secrets \
        --suggestedFeeRecipient $FEE_RECIPIENT \
        $VC_ADDITIONAL_FLAGS"

    if [ "$DOPPELGANGER_DETECTION" = "true" ]; then
        CMD="$CMD --doppelgangerProtection"
    fi

    if [ "$ENABLE_MEV_BOOST" = "true" ]; then
        CMD="$CMD --builder"
    fi

    if [ "$ENABLE_METRICS" = "true" ]; then
        CMD="$CMD --metrics --metrics.address 0.0.0.0 --metrics.port $VC_METRICS_PORT"
    fi

    if [ "$ENABLE_BITFLY_NODE_METRICS" = "true" ]; then
        CMD="$CMD --monitoring.endpoint $BITFLY_NODE_METRICS_ENDPOINT?apikey=$BITFLY_NODE_METRICS_SECRET&machine=$BITFLY_NODE_METRICS_MACHINE_NAME"
    fi

    exec ${CMD} --graffiti "$GRAFFITI"

fi


# Nimbus startup
if [ "$CLIENT" = "nimbus" ]; then

    # Nimbus won't start unless the validator directories already exist
    mkdir -p /validators/nimbus/validators
    mkdir -p /validators/nimbus/secrets

    # Set up the fallback arg
    if [ ! -z "$FALLBACK_BN_API_ENDPOINT" ]; then
        FALLBACK_BN_ARG="--beacon-node=$FALLBACK_BN_API_ENDPOINT"
    fi

    CMD="/home/user/nimbus_validator_client \
        --non-interactive \
        --beacon-node=$BN_API_ENDPOINT $FALLBACK_BN_ARG \
        --data-dir=/ethclient/nimbus_vc \
        --validators-dir=/validators/nimbus/validators \
        --secrets-dir=/validators/nimbus/secrets \
        --doppelganger-detection=$DOPPELGANGER_DETECTION \
        --suggested-fee-recipient=$FEE_RECIPIENT \
        --block-monitor-type=event \
        $VC_ADDITIONAL_FLAGS"

    if [ "$ENABLE_MEV_BOOST" = "true" ]; then
        CMD="$CMD --payload-builder"
    fi

    if [ "$ENABLE_METRICS" = "true" ]; then
        CMD="$CMD --metrics --metrics-address=0.0.0.0 --metrics-port=$VC_METRICS_PORT"
    fi

    # Graffiti breaks if it's in the CMD string instead of here because of spaces
    exec ${CMD} --graffiti="$GRAFFITI"

fi


# Prysm startup
if [ "$CLIENT" = "prysm" ]; then

    # Make the Prysm dir
    mkdir -p /validators/prysm-non-hd/

    # Get rid of the protocol prefix
    BN_RPC_ENDPOINT=$(echo $BN_RPC_ENDPOINT | sed -E 's/.*\:\/\/(.*)/\1/')
    if [ ! -z "$FALLBACK_BN_RPC_ENDPOINT" ]; then
        FALLBACK_BN_RPC_ENDPOINT=$(echo $FALLBACK_BN_RPC_ENDPOINT | sed -E 's/.*\:\/\/(.*)/\1/')
    fi

    # Set up the CC + fallback string
    BN_URL_STRING=$BN_RPC_ENDPOINT
    if [ ! -z "$FALLBACK_BN_RPC_ENDPOINT" ]; then
        BN_URL_STRING="$BN_RPC_ENDPOINT,$FALLBACK_BN_RPC_ENDPOINT"
    fi

    CMD="/app/cmd/validator/validator \
        --accept-terms-of-use \
        --$ETH_NETWORK \
        --wallet-dir /validators/prysm-non-hd \
        --wallet-password-file /validators/prysm-non-hd/direct/accounts/secret \
        --beacon-rpc-provider $BN_URL_STRING \
        --suggested-fee-recipient $FEE_RECIPIENT \
        $VC_ADDITIONAL_FLAGS"

    if [ "$ENABLE_MEV_BOOST" = "true" ]; then
        CMD="$CMD --enable-builder"
    fi

    if [ "$DOPPELGANGER_DETECTION" = "true" ]; then
        CMD="$CMD --enable-doppelganger"
    fi

    if [ "$ENABLE_METRICS" = "true" ]; then
        CMD="$CMD --monitoring-host 0.0.0.0 --monitoring-port $VC_METRICS_PORT"
    else
        CMD="$CMD --disable-account-metrics"
    fi

    exec ${CMD} --graffiti "$GRAFFITI"

fi


# Teku startup
if [ "$CLIENT" = "teku" ]; then

    # Teku won't start unless the validator directories already exist
    mkdir -p /validators/teku/keys
    mkdir -p /validators/teku/passwords

    # Remove any lock files that were left over accidentally after an unclean shutdown
    rm -f /validators/teku/keys/*.lock

    # Set up the CC + fallback string
    BN_URL_STRING=$BN_API_ENDPOINT
    if [ ! -z "$FALLBACK_BN_API_ENDPOINT" ]; then
        BN_URL_STRING="$BN_API_ENDPOINT,$FALLBACK_BN_API_ENDPOINT"
    fi

    CMD="/opt/teku/bin/teku validator-client \
        --network=$ETH_NETWORK \
        --data-path=/validators/teku \
        --validator-keys=/validators/teku/keys:/validators/teku/passwords \
        --beacon-node-api-endpoints=$BN_URL_STRING \
        --validators-keystore-locking-enabled=false \
        --log-destination=CONSOLE \
        --validators-proposer-default-fee-recipient=$FEE_RECIPIENT \
        $VC_ADDITIONAL_FLAGS"

    if [ "$DOPPELGANGER_DETECTION" = "true" ]; then
        CMD="$CMD --doppelganger-detection-enabled"
    fi

    if [ "$ENABLE_MEV_BOOST" = "true" ]; then
        CMD="$CMD --validators-builder-registration-default-enabled=true"
    fi

    if [ "$ENABLE_METRICS" = "true" ]; then
        CMD="$CMD --metrics-enabled=true --metrics-interface=0.0.0.0 --metrics-port=$VC_METRICS_PORT --metrics-host-allowlist=*"
    fi

    if [ "$ENABLE_BITFLY_NODE_METRICS" = "true" ]; then
        CMD="$CMD --metrics-publish-endpoint=$BITFLY_NODE_METRICS_ENDPOINT?apikey=$BITFLY_NODE_METRICS_SECRET&machine=$BITFLY_NODE_METRICS_MACHINE_NAME"
    fi

    exec ${CMD} --validators-graffiti="$GRAFFITI"

fi

