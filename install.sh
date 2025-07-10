#!/bin/bash

set -e

echo "🚀 Starting devkit-cli installation..."

# Detect OS
OS=$(uname -s)
ARCH="amd64"
PLATFORM=""

case "$OS" in
  Linux*)
    PLATFORM="linux"
    ;;
  Darwin*)
    PLATFORM="darwin"
    ;;
  MINGW*|MSYS*|CYGWIN*|Windows_NT)
    echo "❌ Native Windows is not supported by this script. Please use WSL or install manually."
    exit 1
    ;;
  *)
    echo "❌ Unsupported OS: $OS"
    exit 1
    ;;
esac

echo "🧠 Detected OS: $PLATFORM-$ARCH"

# Get latest release tag
echo "📦 Fetching the latest release version..."
RELEASE_TAG=$(curl -s https://api.github.com/repos/darwishdev/devkit-cli/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

# Prepare download URL

TMP_DIR="/tmp/devkit-install"
mkdir -p "$TMP_DIR"
ASSET_FILENAME="devkit-${PLATFORM}-${ARCH}.zip"
FILENAME="$TMP_DIR/$ASSET_FILENAME"
DOWNLOAD_URL="https://github.com/darwishdev/devkit-cli/releases/download/${RELEASE_TAG}/${ASSET_FILENAME}"

echo "⬇️ Downloading $FILENAME (version $RELEASE_TAG)..."

wget -q --show-progress -O "$FILENAME" "$DOWNLOAD_URL"

echo "📂 Extracting $FILENAME..."
unzip -o "$FILENAME" -d "$TMP_DIR"

echo "📁 Setting up the installation directory..."
mkdir -p ~/devkitcli ~/.config/devkitcli /tmp/devkit-install

mv "$TMP_DIR/devkit-${PLATFORM}-${ARCH}" "$TMP_DIR/devkit"
cp "$TMP_DIR/devkit" ~/devkitcli/
rm -rf "$TMP_DIR"
CONFIG_FILE=~/.config/devkitcli/devkit

echo "🛠️ Adding placeholder config if not exists..."
if [ ! -f "$CONFIG_FILE" ]; then
cat <<EOL > "$CONFIG_FILE"
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
fi

echo "🔗 Creating symlink to devkit.env..."
ln -sf "$CONFIG_FILE" ~/devkitcli/devkit.env

echo "✅ devkit-cli installed successfully! Run \`devkit -h\` to get started."

echo -e "\n📌 Next Steps:"
echo "Edit your config at: ~/.config/devkitcli/devkit"
echo "Fill in the required values like API keys and service names."
