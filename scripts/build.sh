#!/bin/bash

# Define the output directory
OUTPUT_DIR=./bin

# Create the output directory if it doesn't exist
mkdir -p $OUTPUT_DIR

# Define the application name
APP_NAME=titan

# Platforms and architectures to build for
platforms=("linux/amd64" "darwin/arm64" "windows/amd64")

for platform in "${platforms[@]}"; do
  IFS="/" read -r GOOS GOARCH <<< "$platform"
  OUTPUT_NAME="$APP_NAME"
  if [ "$GOOS" = "windows" ]; then
    OUTPUT_NAME+=".exe"
  fi
  echo "Building for $GOOS/$GOARCH..."
  DESTINATION=$OUTPUT_DIR/$GOOS-$GOARCH/$OUTPUT_NAME
  env GOOS=$GOOS GOARCH=$GOARCH CGO_ENABLED=0 \
    go build -ldflags="-s -w" -trimpath -o $DESTINATION ./cmd/titan.go
  if [ ! "$GOOS" = "windows" ]; then
      echo "Setting executable permissions for $GOOS platform"
      chmod 755 $DESTINATION
  else
      echo "Cannot set executable permissions for windows"
  fi

done

echo "Build completed. Binaries are located in the $OUTPUT_DIR directory."
