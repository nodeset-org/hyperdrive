# The daemon image
FROM debian:bookworm-slim
ARG TARGETOS TARGETARCH
COPY ./build/hyperdrive-const-daemon-${TARGETOS}-${TARGETARCH} /usr/bin/hyperdrive-const-daemon
RUN apt update && \
    apt install ca-certificates -y && \
	# Cleanup
	apt clean && \
        rm -rf /var/lib/apt/lists/*

# Container entry point
ENTRYPOINT ["/usr/bin/hyperdrive-const-daemon"]