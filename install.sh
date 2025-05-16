#!/usr/bin/env bash

set -e

REPO="valcinei/jiboia-tunnel"
BINARY="jiboia"
LATEST=$(curl -s https://api.github.com/repos/$REPO/releases/latest | grep tag_name | cut -d '"' -f4)
OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

if [[ "$ARCH" == "x86_64" ]]; then
  ARCH="amd64"
elif [[ "$ARCH" == "aarch64" || "$ARCH" == "arm64" ]]; then
  ARCH="arm64"
else
  echo "Unsupported architecture: $ARCH"
  exit 1
fi

URL="https://github.com/$REPO/releases/download/$LATEST/${BINARY}-${OS}-${ARCH}.zip"
echo "Downloading $URL..."

curl -L "$URL" -o "$BINARY.zip"
unzip "$BINARY.zip"
chmod +x "$BINARY-${OS}-${ARCH}"
sudo mv "$BINARY-${OS}-${ARCH}" /usr/local/bin/$BINARY

echo "✅ $BINARY installed to /usr/local/bin"
$BINARY --help
