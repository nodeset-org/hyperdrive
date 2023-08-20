#/bin/sh

## Runs as a user via postStartCommand from devcontainer.json

cd ~

rocketpool --debug service install -y -p .rocketpool
cp bin/hyperdrive/.devcontainer/user-settings.yml .rocketpool/user-settings.yml

#sudo dockerd && containerd

rocketpool --debug service start -y --ignore-slash-timer

rocketpool --debug wallet init -p thisisabigtest --confirm-mnemonic

echo "{::} Hyperdrive development environment enabled! {::}\n"