#/bin/sh

## Runs every time the container starts

su vscode
cd ~

# dockerd &
# containerd &

rocketpool --debug service start -y --ignore-slash-timer

echo "{::} Hyperdrive development environment enabled! {::}\n"

exec "$@"