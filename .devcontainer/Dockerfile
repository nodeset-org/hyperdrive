# Hyperdrive testing image

FROM mcr.microsoft.com/devcontainers/go

RUN useradd -ms /bin/bash nodeop
USER nodeop
WORKDIR /home/nodeop

# add RP
ADD https://github.com/rocket-pool/smartnode-install/releases/latest/download/rocketpool-cli-linux-amd64 /home/nodeop/bin/rocketpool
# RUN chmod +x ~/bin/rocketpool
