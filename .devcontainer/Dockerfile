# Hyperdrive testing image
FROM mcr.microsoft.com/devcontainers/go

ADD --chown=vscode:vscode https://github.com/rocket-pool/smartnode-install/releases/latest/download/rocketpool-cli-linux-amd64 bin/rocketpool

USER root
RUN groupadd docker \
    && usermod -aG docker vscode \
    && chmod +x bin/rocketpool