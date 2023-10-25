# Hyperdrive testing image
FROM mcr.microsoft.com/devcontainers/go

# Declare username variable
ARG USERNAME=testuser

# Create a new user using the username variable
RUN useradd -m $USERNAME

# Set password for the created user (e.g., "password123")
RUN echo "$USERNAME:password" | chpasswd

# Change to the created user
USER $USERNAME

# Set up directories and PATH for user's binaries
RUN mkdir -p /home/$USERNAME/bin && \
    echo "export PATH=\$PATH:/home/$USERNAME/bin" >> /home/$USERNAME/.bashrc

# Add rocketpool
ADD --chown=$USERNAME:$USERNAME https://github.com/rocket-pool/smartnode-install/releases/latest/download/rocketpool-cli-linux-amd64 /home/$USERNAME/bin/rocketpool

# Add compiled Go files
# ADD --chown=$USERNAME:$USERNAME ../bin/* /home/$USERNAME/hyperdrive/

# Only copy over necessary go related files
COPY ./src ./src

# Set working directory
WORKDIR /go/src

# Build the Go application
RUN go build -o ../bin/hyperdrive hyperdrive.go

# Return to the root user to set permissions
USER root
RUN chmod +x /home/$USERNAME/bin/rocketpool

# Add and execute the setup script
COPY build/setup.sh /home/$USERNAME/setup.sh
RUN chmod +x /home/$USERNAME/setup.sh \
    && sh -c /home/$USERNAME/setup.sh >> /home/$USERNAME/setup.log

#TODO: Make sure the user-settings.yml file is copied over to the correct location
COPY build/user-settings.yml ~/user-settings.yml

# Switch back to the defined user for any further operations
USER $USERNAME
