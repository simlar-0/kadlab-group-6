#!/bin/bash

go test ./tests -v -coverpkg=./pkg/kademlia_node -coverprofile=cover.out
go tool cover -html cover.out -o cover.html

# MacOS
#open cover.html

# Windows
start cover.html

# Linux
#xdg-open cover.html