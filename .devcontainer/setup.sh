#/bin/sh

## Installs dependencies only when the container is built

su vscode
cd ~
sudo chown vscode:vscode .rocketpool 
rocketpool --debug service install -y -p .rocketpool
cp ~/bin/hyperdrive/.devcontainer/user-settings.yml ~/.rocketpool/user-settings.yml