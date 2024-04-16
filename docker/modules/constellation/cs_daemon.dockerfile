# The daemon image
FROM debian:bookworm-slim
ARG BINARIES_PATH
ARG TARGETOS TARGETARCH
COPY ${BINARIES_PATH}/hyperdrive-constellation-daemon-${TARGETOS}-${TARGETARCH} /usr/bin/hyperdrive-constellation-daemon
RUN apt update && \
    apt install ca-certificates -y && \
    # Cleanup
    apt clean && \
    rm -rf /var/lib/apt/lists/*

# Container entry point
ENTRYPOINT ["/usr/bin/hyperdrive-constellation-daemon"]