#!/bin/bash

# Define the output directory
OUTPUT_DIR=$1

# Define the application name
APP_NAME=titan

echo "*************************************"
echo "* Releasing $APP_NAME"
echo "*************************************"

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

  # echo "Compress with UPX for $GOOS platform"
  # case "$GOOS" in
  #   linux*)     upx --best --lzma $DESTINATION || true;;
  #   darwin*)    upx --best --lzma --force-macos $DESTINATION || true;;
  #   windows*)   upx --best --lzma $DESTINATION || true;;
  #   *)          echo "$0 - UNKNOWN:${GOOS}" && exit 1
  # esac

  case "$GOOS" in
    linux*)     echo "Setting executable permissions for $GOOS platform" && chmod 755 $DESTINATION;;
    darwin*)    echo "Setting executable permissions for $GOOS platform" && chmod 755 $DESTINATION;;
    windows*)   echo "Cannot set executable permissions for windows";;
    *)          echo "$0 - UNKNOWN:${GOOS}" && exit 1
  esac

  echo "Binary for $GOOS/$GOARCH can be found in $DESTINATION"
done

echo "Build completed. Binaries are located in the $OUTPUT_DIR directory."
