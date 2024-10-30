#!/bin/bash

# Get the latest release tag
RELEASE_TAG=$(curl -s https://api.github.com/repos/darwishdev/devkit-cli/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

# Download the release archive
wget https://github.com/darwishdev/devkit-cli/releases/download/$RELEASE_TAG/devkit-cli-refs.tags.$RELEASE_TAG.zip

# Extract the archive
unzip devkit-cli-refs.tags.$RELEASE_TAG.zip 

# rename the file
mkdir -p ~/devkitcli ~/.config/devkitcli && cp -r release/* ~/devkitcli/

touch ~/.config/devkitcli/devkit

# Create a symlink to the binary
sudo ln -s $HOME/devkitcli/devkit /usr/local/bin/devkit

echo "devkit-cli installed successfully! try devkit -h"
