FROM golang:1.23.1-alpine3.20

# Add the commands needed to put your compiled go binary in the container and
# run it when the container starts.
#
# See https://docs.docker.com/engine/reference/builder/ for a reference of all
# the commands you can use in this file.
#
# In order to use this file together with the docker-compose.yml file in the
# same directory, you need to ensure the image you build gets the name
# "kadlab", which you do by using the following command:
#
# $ docker build . -t kadlab

# Set working directory
WORKDIR /app

# Copy the source code into the container
COPY . .

# Install dependencies
RUN go mod download

# Build the binary
RUN go build -o kadlab ./cmd/kademlia_main

# Run the binary
ENTRYPOINT ["./kadlab"]