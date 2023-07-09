# Hyperdrive testing image

# this doesn't work??? if this is used during the ADD step later it will create a folder with docker files inside /home/nodeop/bin/rocketpool instead of expanding the variable as expected
# ARG RP_REMOTE=https://github.com/rocket-pool/smartnode-install/releases/latest/download/rocketpool-cli-linux-amd64

FROM mcr.microsoft.com/devcontainers/go

RUN useradd -ms /bin/bash nodeop
WORKDIR /home/nodeop

# install RP
USER nodeop
ADD https://github.com/rocket-pool/smartnode-install/releases/latest/download/rocketpool-cli-linux-amd64 bin/rocketpool

USER root
RUN groupadd docker \
    && usermod -aG docker root \
    && usermod -aG docker nodeop \
    && chmod +x bin/rocketpool

# RP requires this var to be set, but the noninteractive shell
# that docker uses doesn't set it by default
ENV USER root
RUN bin/rocketpool --allow-root service install -y -p .rocketpool \
    # I'm a cowboy pew pew
    && chmod -R 777 .rocketpool

# TODO: enable and start docker daemons

USER nodeop
COPY user-settings.yml .rocketpool/user-settings.yml