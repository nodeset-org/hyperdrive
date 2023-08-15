#/bin/sh
cd ~
sudo chown vscode:vscode .rocketpool 
rocketpool --debug service install -y -p .rocketpool

# dockerd &
# containerd &
cp ~/bin/hyperdrive/.devcontainer/user-settings.yml ~/.rocketpool/user-settings.yml
rocketpool --debug service start -y --ignore-slash-timer

echo "{::} Hyperdrive development environment enabled! {::}\n"