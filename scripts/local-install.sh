#!/usr/bin/env bash
set -euo pipefail

. $(dirname $0)/commons.sh

PLUGINS_DIR="$HOME/.terraform.d/plugins/"
mkdir -p "$PLUGINS_DIR"

info "Installing provider to $PLUGINS_DIR"
cp -v "$BUILD_DIR"* "$PLUGINS_DIR/"

## Terraform >= v0.13 requires different layout
PROVIDER_NAME=$(cat "$PROVIDER_NAME_FILE")
PROVIDER_VERSION=$(cat "$PROVIDER_VERSION_FILE")
# patched versus microsoft-origin below
PROVIDER_REGISTRY='babylonhealth.com'
PROVIDER_ORGANIZATION='babylonhealth'
# patched versus microsoft-origin above
PROVIDER_SOURCE_ADDRESS="${PROVIDER_ORGANIZATION}/${PROVIDER_NAME}"

PLUGINS_DIR="${PLUGINS_DIR}${PROVIDER_REGISTRY}/${PROVIDER_SOURCE_ADDRESS}/${PROVIDER_VERSION}/${OS}_${PROC}"
info "Installing provider to $PLUGINS_DIR"
mkdir -p "$PLUGINS_DIR"
cp -v "$BUILD_DIR"* "$PLUGINS_DIR/"
