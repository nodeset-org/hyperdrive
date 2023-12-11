# The builder for building the daemon
FROM golang:1.21-bookworm AS builder
COPY . /hyperdrive
ENV CGO_ENABLED=1
ENV CGO_CFLAGS="-O -D__BLST_PORTABLE__"
RUN cd /hyperdrive/hyperdrive-daemon && go build

# The actual image
FROM debian:bookworm-slim
RUN apt update && \
	# Install dependencies
	apt install -y curl gpg lsb-release && \
	# Add the Docker repo GPG key
	curl -fsSL "https://download.docker.com/linux/debian/gpg" | gpg --dearmor -o /etc/apt/keyrings/docker.gpg && \
	# Add the Docker repo to apt
	echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/debian $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null && \
	mkdir -p /etc/apt/keyrings && \
	apt update && \
	# Install the CLI
	apt install -y docker-ce-cli && \
	# Cleanup
	apt purge -y curl gpg lsb-release && \
	apt clean && \
        rm -rf /var/lib/apt/lists/*

# Grab the daemon from the builder
COPY --from=builder /hyperdrive/hyperdrive-daemon/hyperdrive-daemon /usr/local/bin/hyperdrive-daemon
ENTRYPOINT ["hyperdrive-daemon"]
