# The builder for building the Debian package
FROM nodeset/hyperdrive-deb-builder:v1.0.1 AS builder
ARG BUILDPLATFORM

# Debian packages need a very particular folder structure, so we're basically converting the repo structure into what it wants here 
COPY ./install/packages/debian/debian /hyperdrive/debian/debian
COPY ./install/deploy /hyperdrive/debian/deploy
COPY ./src/ /hyperdrive/debian/
WORKDIR /hyperdrive/debian

# Build the native arch package and source package
RUN DEB_BUILD_OPTIONS=noautodbgsym debuild -us -uc

# Build the other package (binary only since we already made the source)
RUN if [ "$BUILDPLATFORM" = "linux/arm64" ]; then \
        REMOTE_ARCH=amd64; \
    else \
        REMOTE_ARCH=arm64; \
    fi && \
    DEB_BUILD_OPTIONS=noautodbgsym debuild -us -uc -b -a${REMOTE_ARCH} --no-check-builddeps

# Copy the output
FROM scratch AS package
COPY --from=builder /hyperdrive/hyperdrive* /