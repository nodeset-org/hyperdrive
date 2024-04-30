# The builder for building the CLIs
FROM golang:1.21-bookworm AS builder
COPY . /hyperdrive
ENV CGO_ENABLED=0
WORKDIR /hyperdrive/hyperdrive-cli

# Build x64 version
RUN GOOS=linux GOARCH=amd64 go build -o /build/hyperdrive-cli-linux-amd64
RUN GOOS=darwin GOARCH=amd64 go build -o /build/hyperdrive-cli-darwin-amd64

# Build the arm64 version
RUN GOOS=linux GOARCH=arm64 go build -o /build/hyperdrive-cli-linux-arm64
RUN GOOS=darwin GOARCH=arm64 go build -o /build/hyperdrive-cli-darwin-arm64

# Copy the output
FROM scratch AS cli
COPY --from=builder /build /