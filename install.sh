#!/bin/sh
# Download this script from a tagged release, review it, then run: sh install.sh
set -eu
VERSION=0.3.0
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case $(uname -m) in x86_64|amd64) ARCH=amd64;; arm64|aarch64) ARCH=arm64;; *) echo 'Unsupported architecture' >&2; exit 2;; esac
case "$OS" in darwin|linux) ;; *) echo 'Use install.ps1 on Windows' >&2; exit 2;; esac
TEMP=$(mktemp -d)
trap 'rm -rf "$TEMP"' EXIT HUP INT TERM
ASSET="mixrank_${VERSION}_${OS}_${ARCH}.tar.gz"
BASE="https://github.com/nkulavic/mixrank/releases/download/v$VERSION"
curl -fsSL --proto '=https' --tlsv1.2 "$BASE/$ASSET" -o "$TEMP/$ASSET"
curl -fsSL --proto '=https' --tlsv1.2 "$BASE/checksums.txt" -o "$TEMP/checksums.txt"
EXPECTED=$(awk -v file="$ASSET" '$2==file {print $1}' "$TEMP/checksums.txt")
[ ${#EXPECTED} -eq 64 ] || { echo 'Missing release checksum' >&2; exit 2; }
if command -v sha256sum >/dev/null 2>&1; then ACTUAL=$(sha256sum "$TEMP/$ASSET" | awk '{print $1}'); else ACTUAL=$(shasum -a 256 "$TEMP/$ASSET" | awk '{print $1}'); fi
[ "$ACTUAL" = "$EXPECTED" ] || { echo 'Release checksum mismatch' >&2; exit 2; }
tar -xzf "$TEMP/$ASSET" -C "$TEMP" mixrank
chmod 755 "$TEMP/mixrank"
"$TEMP/mixrank" setup "$@"
