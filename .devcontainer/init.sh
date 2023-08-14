#/bin/sh
cd ~/bin
rocketpool --debug service install -y -p .rocketpool \

# dockerd &
# containerd & 
cd ~
rocketpool --debug service start
cp user-settings.yml .rocketpool/user-settings.yml

echo "Hyperdrive development environment enabled"