#!/bin/bash

set -e

REPO="Omotolani98/tiny.sh"
VERSION="${1:-latest}"

echo "🔍 Detecting platform..."
ARCH=$(uname -m)
OS=$(uname -s | tr '[:upper:]' '[:lower:]')

case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  aarch64 | arm64) ARCH="arm64" ;;
  *) echo "❌ Unsupported architecture: $ARCH"; exit 1 ;;
esac

if [ "$VERSION" = "latest" ]; then
  VERSION=$(curl -s https://api.github.com/repos/${REPO}/releases/latest | grep tag_name | cut -d '"' -f 4)
fi

VERSION_NO_V="${VERSION#v}"
TARBALL="tiny.sh_${VERSION_NO_V}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/${VERSION}/${TARBALL}"

echo "⬇️ Downloading $TARBALL..."
curl -L "$URL" -o "$TARBALL"

echo "📦 Extracting..."
tar -xzf "$TARBALL"

INSTALL_PATH="/usr/local/bin"
FALLBACK_PATH="$HOME/.local/bin"

if command -v sudo >/dev/null; then
  echo "🚚 Installing to $INSTALL_PATH..."
  chmod +x tiny.sh
  sudo mv tiny.sh $INSTALL_PATH/tiny.sh || {
    echo "⚠️ Failed to move to $INSTALL_PATH. Falling back to $FALLBACK_PATH"
    mkdir -p "$FALLBACK_PATH"
    mv tiny.sh "$FALLBACK_PATH/tiny.sh"
    echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
  }
else
  echo "⚠️ sudo not found. Installing to $FALLBACK_PATH"
  mkdir -p "$FALLBACK_PATH"
  mv tiny.sh "$FALLBACK_PATH/tiny.sh"
  echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
fi

echo "🧹 Cleaning up..."
rm "$TARBALL"

echo "✅ tiny.sh $VERSION installed!"
echo "👉 Restart your shell or run 'source ~/.bashrc' if needed"


