#!/bin/bash

go test ./tests -v -coverpkg=./pkg/kademlia_node -coverprofile=cover.out
go tool cover -html cover.out -o cover.html

# Function to open cover.html based on the operating system
open_cover_html() {
    case "$(uname -s)" in
        (Darwin)
            # macOS
            open cover.html
            ;;
        (Linux)
            # Linux
            xdg-open cover.html
            ;;
        (CYGWIN*|MINGW32*|MSYS*|MINGW*)
            # Windows
            start cover.html
            ;;
        (*)
            echo "Unsupported OS"
            ;;
    esac
}

# Call the function to open cover.html
open_cover_html