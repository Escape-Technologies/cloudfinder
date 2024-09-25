#!/bin/bash

# GitHub repository details
REPO="Escape-Technologies/cloudfinder"

# Determine the latest version tag from GitHub API
LATEST_TAG=$(curl -s https://api.github.com/repos/$REPO/releases/latest | grep tag_name | cut -d '"' -f 4)

# Check if the tag was found
if [ -z "$LATEST_TAG" ]; then
    echo "Unable to find the latest release tag. Exiting."
    exit 1
fi

# Remove the 'v' prefix to get the release version
LATEST_RELEASE=${LATEST_TAG#v}

# Determine OS and ARCH
OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Translate uname output to Go architecture naming convention
if [ "$ARCH" == "x86_64" ]; then
    ARCH="amd64"
elif [ "$ARCH" == "arm64" ]; then
    ARCH="arm64"
else
    echo "Unsupported architecture: $ARCH"
    exit 1
fi

# Define the filename based on the OS and architecture
FILENAME="cloudfinder_${LATEST_RELEASE}_${OS}_${ARCH}.tar.gz"
DOWNLOAD_URL="https://github.com/$REPO/releases/download/${LATEST_TAG}/${FILENAME}"

# Download the tar.gz file
echo "Downloading $FILENAME from $DOWNLOAD_URL..."
curl -L -o "$FILENAME" "$DOWNLOAD_URL"

# Check if the file was downloaded successfully
if [ ! -f "$FILENAME" ]; then
    echo "Failed to download $FILENAME. Exiting."
    exit 1
fi

# Extract cloudfinder bin from the tar.gz file
echo "Extracting $FILENAME..."
tar -xzvf "$FILENAME" cloudfinder

# Check if the binary exists after extraction
if [ ! -f "cloudfinder" ]; then
    echo "Extraction failed or binary not found. Exiting."
    exit 1
fi

# Move the binary to /usr/local/bin for global access
echo "Installing cloudfinder to /usr/local/bin..."
sudo mv cloudfinder /usr/local/bin/

# Clean up
rm "$FILENAME"

echo "cloudfinder ${LATEST_TAG} installed successfully!"
