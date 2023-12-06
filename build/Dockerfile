# Hyperdrive testing image
FROM mcr.microsoft.com/devcontainers/go

# Declare username variable
ARG USERNAME=testuser

# Create a new user using the username variable
RUN useradd -m $USERNAME

# TODO: Env var
# Set password for the created user (e.g., "password123")
RUN echo "$USERNAME:password" | chpasswd

RUN echo "$USERNAME ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers.d/$USERNAME && \
    chmod 0440 /etc/sudoers.d/$USERNAME

# Change to the created user
USER $USERNAME

# Set up directories and PATH for user's binaries
RUN mkdir -p /home/$USERNAME/bin && \
    echo "export PATH=\$PATH:/home/$USERNAME/bin" >> /home/$USERNAME/.bashrc

# Add rocketpool
ADD --chown=$USERNAME:$USERNAME https://github.com/rocket-pool/smartnode-install/releases/latest/download/rocketpool-cli-linux-amd64 /home/$USERNAME/bin/rocketpool

# Only copy over necessary go related files
COPY ./hyperdrive-daemon ./hyperdrive-daemon
COPY ./hyperdrive ./hyperdrive

COPY ./go.mod ./hyperdrive/go.mod
COPY ./go.sum ./hyperdrive/go.sum

COPY ./go.mod ./hyperdrive-daemon/go.mod
COPY ./go.sum ./hyperdrive-daemon/go.sum

# Set working directory and build the Go application
WORKDIR /go/hyperdrive
RUN go build -o ../bin/hyperdrive main.go

# Set working directory and build the Go application
WORKDIR /go/hyperdrive-daemon
RUN go build -o ../bin/hyperdrive-daemon main.go

# # Return to the root user to set permissions
# USER root

# # Add and execute the setup script
# COPY build/setup.sh /home/$USERNAME/setup.sh

# RUN mkdir -p /home/$USERNAME/.rocketpool && \
#     chmod +x /home/$USERNAME/bin/rocketpool && \
#     chown -R $USERNAME:$USERNAME /etc/apt && \
#     chown -R $USERNAME:$USERNAME /var && \
#     chown -R $USERNAME:$USERNAME /home/$USERNAME/.rocketpool && \
#     chmod +x /home/$USERNAME/setup.sh  && \
#     sh -c /home/$USERNAME/setup.sh >> /home/$USERNAME/setup.log

# #TODO: Make sure the user-settings.yml file is copied over to the correct location
# COPY build/user-settings.yml /home/$USERNAME/.rocketpool/user-settings.yml

# # Switch back to the defined user for any further operations
# USER $USERNAME
