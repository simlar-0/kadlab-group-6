#!/bin/bash

# Default value
DEFAULT_VALUE=50
# Use the first argument if provided, otherwise use the default value
VALUE=${1:-$DEFAULT_VALUE}

# Define the project name
PROJECT_NAME="kadlab_group_6"

# Stop and remove all running Docker containers
docker-compose -p $PROJECT_NAME down
docker rm -f $(docker ps -aq --filter "label=com.docker.compose.project=$PROJECT_NAME")

# Build the Docker image
docker build . -t kadlab

# Deploy the Docker Compose stack with scaling
docker-compose -p $PROJECT_NAME up --scale kademlia-node=$VALUE -d