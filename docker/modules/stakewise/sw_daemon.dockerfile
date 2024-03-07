# The builder for building the daemon
FROM --platform=${BUILDPLATFORM} golang:1.21-bookworm AS builder
ARG TARGETOS TARGETARCH BUILDPLATFORM
COPY . /hyperdrive
ENV CGO_ENABLED=1
ENV CGO_CFLAGS="-O -D__BLST_PORTABLE__"
RUN if [ "$BUILDPLATFORM" = "linux/amd64" -a "$TARGETARCH" = "arm64" ]; then \
        # Install the GCC cross compiler
        apt update && apt install -y gcc-aarch64-linux-gnu g++-aarch64-linux-gnu && \
        export CC=aarch64-linux-gnu-gcc && export CC_FOR_TARGET=gcc-aarch64-linux-gnu; \
    elif [ "$BUILDPLATFORM" = "linux/arm64" -a "$TARGETARCH" = "amd64" ]; then \
        apt update && apt install -y gcc-x86-64-linux-gnu g++-x86-64-linux-gnu && \
        export CC=x86_64-linux-gnu-gcc && export CC_FOR_TARGET=gcc-x86-64-linux-gnu; \
    fi && \
    cd /hyperdrive/src/modules/stakewise/stakewise-daemon && \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /build/hyperdrive-stakewise-daemon-${TARGETOS}-${TARGETARCH}

# The daemon image
FROM debian:bookworm-slim
ARG TARGETOS TARGETARCH
COPY --from=builder /build/hyperdrive-stakewise-daemon-${TARGETOS}-${TARGETARCH} /usr/bin/hyperdrive-stakewise-daemon
RUN apt update && \
    apt install ca-certificates -y && \
	# Cleanup
	apt clean && \
        rm -rf /var/lib/apt/lists/*

# Container entry point
ENTRYPOINT ["/usr/bin/hyperdrive-stakewise-daemon"]