# Image for building Hyperdrive debian packages
FROM debian:bookworm-slim

# Add the testing (trixie) repo because as of right now bookworm doesn't have Go 1.24 yet
RUN cat <<'EOF' > /etc/apt/sources.list
deb http://deb.debian.org/debian testing main contrib non-free-firmware
deb-src http://deb.debian.org/debian testing main contrib non-free-firmware
EOF

# Install dependencies
RUN apt update && \
    apt install -y -t bookworm devscripts lintian binutils-x86-64-linux-gnu binutils-aarch64-linux-gnu nano && \
    apt install -y -t testing golang-any && \
    apt install -y -t testing dh-golang && \
    # Cleanup
    apt clean && \
    rm -rf /var/lib/apt/lists/*
