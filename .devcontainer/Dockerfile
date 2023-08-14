# Hyperdrive testing image

# this doesn't work??? if this is used during the ADD step later it will 
# create a folder with docker files inside /home/nodeop/bin/rocketpool instead of expanding the variable as expected
# ARG RP_REMOTE=https://github.com/rocket-pool/smartnode-install/releases/latest/download/rocketpool-cli-linux-amd64

FROM mcr.microsoft.com/devcontainers/go

WORKDIR /home/vscode

# RUN passwd -d vscode

ADD --chown=vscode:vscode https://github.com/rocket-pool/smartnode-install/releases/latest/download/rocketpool-cli-linux-amd64 bin/rocketpool

USER root
RUN groupadd docker \
    && usermod -aG docker vscode \
    && chmod +x bin/rocketpool

USER vscode
#echo devpwd | sudo -S sh 
# CMD [ "/bin/sh", "-c", "init.sh >> init.log" ]

WORKDIR /tmp
COPY init.sh init.sh
# RUN /bin/sh init.sh >> init.log