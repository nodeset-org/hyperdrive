#/bin/sh

## Runs as a user

cd ~

echo Installing RP service...
rocketpool --debug service install -y -p .rocketpool

cp bin/hyperdrive/.devcontainer/user-settings.yml ~/.rocketpool/user-settings.yml

echo Starting RP service...
rocketpool --debug service start -y --ignore-slash-timer

echo Initializing RP node wallet...
rocketpool -s --debug wallet init -p thisisabigtest --confirm-mnemonic

echo "{::} Hyperdrive development environment enabled! {::}\n"