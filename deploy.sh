#!/bin/bash

# Default value
DEFAULT_VALUE=50
# Use the first argument if provided, otherwise use the default value
VALUE=${1:-$DEFAULT_VALUE}

# Step 1: Stop and remove all running Docker containers
docker-compose down
docker rm -f $(docker ps -aq)

# Step 2: Build the Docker image
docker build . -t kadlab

# Step 3: Deploy the Docker Compose stack with scaling
docker-compose up --scale kademliaNodes=$VALUE -d