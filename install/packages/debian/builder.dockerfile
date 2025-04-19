# Image for building Hyperdrive debian packages
FROM debian:bookworm-slim

# Add backports and the unstable repo
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
