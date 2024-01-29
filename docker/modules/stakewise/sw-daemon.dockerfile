# The builder for building the daemon
FROM golang:1.21-bookworm AS builder
ARG TARGETARCH
COPY . /hyperdrive
ENV CGO_ENABLED=1
ENV CGO_CFLAGS="-O -D__BLST_PORTABLE__"
RUN cd /hyperdrive/modules/stakewise && go build -o /build/hyperdrive-stakewise-daemon-linux-${TARGETARCH}

# The daemon image
FROM debian:bookworm-slim
ARG TARGETARCH
COPY --from=builder /build/hyperdrive-stakewise-daemon-linux-${TARGETARCH} /usr/bin/hyperdrive-stakewise-daemon
RUN apt update && \
    apt install ca-certificates -y && \
	# Cleanup
	apt clean && \
        rm -rf /var/lib/apt/lists/*

# Container entry point
ENTRYPOINT ["/usr/bin/hyperdrive-stakewise-daemon"]