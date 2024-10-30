#!/bin/bash
echo "Starting devkit-cli installation..."


echo "Fetching the latest release version..."
RELEASE_TAG=$(curl -s https://api.github.com/repos/darwishdev/devkit-cli/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')


echo "Downloading devkit-cli version $RELEASE_TAG..."
wget https://github.com/darwishdev/devkit-cli/releases/download/$RELEASE_TAG/devkit-cli-refs.tags.$RELEASE_TAG.zip

echo "Extracting files..."
unzip devkit-cli-refs.tags.$RELEASE_TAG.zip 

echo "Setting up the installation directory..."
mkdir -p ~/devkitcli ~/.config/devkitcli && cp -r release/* ~/devkitcli/ && rm -rf release
echo "adding the placeholders on config files"
cat <<EOL > ~/.config/devkitcli/devkit
GIT_USER=yourgituser
DOCKER_HUB_USER=youdockeruser
BUF_USER=yourbufuser
GITHUB_TOKEN=yourgittoken
API_SERVICE_NAME=devkit
API_VERSION=v1
GOOGLE_CLIENT_ID=google_client_id
GOOGLE_CLIENT_SECRET=google_client_secret
RESEND_API_KEY=resend_key
EOL
echo "Creating configuration file..."
touch ~/.config/devkitcli/devkit

sudo ln -s $HOME/devkitcli/devkit /usr/local/bin/devkit.env

echo "devkit-cli installed successfully! try devkit -h"

# Step 7: Prompt user to edit the configuration file
echo -e "\nNext Steps:"
echo "Please open the configuration file located at ~/.config/devkitcli/devkit to add your settings."
echo "You will need to configure necessary fields, such as API keys and other required settings, to start using devkit-cli effectively."

echo "Installation complete. Happy coding with devkit-cli!"

