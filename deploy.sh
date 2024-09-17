#!/bin/bash

# Step 1: Stop and remove all running Docker containers
docker stack rm kadlab
docker rm -f $(docker ps -aq)

# Step 2: Build the Docker image
docker build . -t kadlab

# Step 3: Initialize Docker Swarm (if not already initialized)
docker swarm init || echo "Swarm already initialized"

# Step 4: Deploy the Docker stack
docker stack deploy -c docker-compose.yml kadlab