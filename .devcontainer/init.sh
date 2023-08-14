rocketpool service install -y -p .rocketpool \
    # I'm a cowboy pew pew
    && chmod -R 777 .rocketpool

# dockerd &
# containerd & 
rocketpool service start
cp user-settings.yml .rocketpool/user-settings.yml

echo "Hyperdrive development environment enabled"