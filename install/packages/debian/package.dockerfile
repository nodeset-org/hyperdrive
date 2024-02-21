# The builder for building the Debian package
FROM nodeset/hyperdrive-deb-builder:v0.1.0 AS builder
COPY . /hyperdrive
WORKDIR /hyperdrive/install/packages/debian

# Build the amd64 package and source package
RUN DEB_BUILD_OPTIONS=noautodbgsym debuild -us -uc

# Build the arm64 package (binary only since we already made the source)
RUN DEB_BUILD_OPTIONS=noautodbgsym debuild -us -uc -b -aarm64 --no-check-builddeps

# Copy the output
FROM scratch AS package
COPY --from=builder /hyperdrive/install/packages/hyperdrive* /