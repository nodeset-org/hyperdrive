# Hyperdrive testing image
FROM mcr.microsoft.com/devcontainers/go

ADD --chown=vscode:vscode https://github.com/rocket-pool/smartnode-install/releases/latest/download/rocketpool-cli-linux-amd64 /home/vscode/bin/rocketpool

RUN groupadd docker \
    && usermod -aG docker vscode \
    && chmod +x /home/vscode/bin/rocketpool