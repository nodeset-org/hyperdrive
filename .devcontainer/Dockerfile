# Hyperdrive testing image

# this doesn't work??? if this is used during the ADD step later it will create a folder with docker files inside /home/nodeop/bin/rocketpool instead of expanding the variable as expected
# ARG RP_REMOTE=https://github.com/rocket-pool/smartnode-install/releases/latest/download/rocketpool-cli-linux-amd64

FROM mcr.microsoft.com/devcontainers/go

RUN useradd -ms /bin/bash nodeop
WORKDIR /home/nodeop

# install RP
USER nodeop
ADD https://github.com/rocket-pool/smartnode-install/releases/latest/download/rocketpool-cli-linux-amd64 bin/rocketpool
# TODO: rocketpool installation is supposed to create these directories but does not...
# RUN mkdir .rocketpool .rocketpool/templates .rocketpool/override
USER root
RUN chmod +x bin/rocketpool
RUN bin/rocketpool --allow-root service install -y -p .rocketpool

USER nodeop
COPY user-settings.yml .rocketpool/user-settings.yml