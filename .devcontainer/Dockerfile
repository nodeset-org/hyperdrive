# Hyperdrive testing image

# this doesn't work?
#ARG RP_REMOTE=https://github.com/rocket-pool/smartnode-install/releases/latest/download/rocketpool-cli-linux-amd64

FROM mcr.microsoft.com/devcontainers/go

RUN useradd -ms /bin/bash nodeop
WORKDIR /home/nodeop

# add RP
ADD https://github.com/rocket-pool/smartnode-install/releases/latest/download/rocketpool-cli-linux-amd64 /home/nodeop/bin/rocketpool
RUN chmod +x /home/nodeop/bin/rocketpool
RUN mkdir /home/nodeop/.rocketpool
RUN /home/nodeop/bin/rocketpool --allow-root service install -y

USER nodeop