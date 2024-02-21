# The builder for building the daemon
FROM --platform=${BUILDPLATFORM} golang:1.21-bookworm AS builder
ARG TARGETOS TARGETARCH
COPY . /hyperdrive
ENV CGO_ENABLED=1
ENV CGO_CFLAGS="-O -D__BLST_PORTABLE__"
RUN if [ "$TARGETARCH" = "arm64" ]; then \
        # Install the GCC cross compiler
        apt update && apt install -y gcc-aarch64-linux-gnu g++-aarch64-linux-gnu && \
        export CC=aarch64-linux-gnu-gcc && export CC_FOR_TARGET=gcc-aarch64-linux-gnu; \
    fi && \
    cd /hyperdrive/src/hyperdrive-daemon && \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /build/hyperdrive-daemon-${TARGETOS}-${TARGETARCH}

# The daemon image
FROM debian:bookworm-slim
ARG TARGETOS TARGETARCH
COPY --from=builder /build/hyperdrive-daemon-${TARGETOS}-${TARGETARCH} /usr/bin/hyperdrive-daemon
RUN apt update && \
    apt install ca-certificates -y && \
	# Cleanup
	apt clean && \
        rm -rf /var/lib/apt/lists/*

# Container entry point
ENTRYPOINT ["/usr/bin/hyperdrive-daemon"]