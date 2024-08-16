#!/bin/bash

# Build the Go application
go build -o ttg

# Check if the build was successful
if [ $? -ne 0 ]; then
    echo "Build failed. Exiting."
    exit 1
fi

# Move the binary to /usr/local/bin
sudo mv ttg /usr/local/bin/

# Check if the move was successful
if [ $? -ne 0 ]; then
    echo "Failed to move the binary. Exiting."
    exit 1
fi

echo "Build and move completed successfully."