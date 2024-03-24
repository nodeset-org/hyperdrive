# The daemon image
FROM debian:bookworm-slim
ARG TARGETOS TARGETARCH
COPY ./build/hyperdrive-stakewise-daemon-${TARGETOS}-${TARGETARCH} /usr/bin/hyperdrive-stakewise-daemon
RUN apt update && \
    apt install ca-certificates -y && \
	# Cleanup
	apt clean && \
        rm -rf /var/lib/apt/lists/*

# Container entry point
ENTRYPOINT ["/usr/bin/hyperdrive-stakewise-daemon"]