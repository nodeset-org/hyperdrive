#!/bin/bash

# builds the dockerfile locally
docker build -t hyperdrive-image -f build/Dockerfile .

# runs the dockerfile that was just built
docker run -it hyperdrive-image
