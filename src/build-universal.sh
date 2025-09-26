#!/bin/bash
set -e

APP_NAME="Taz GAG Macro"
BUILD_DIR="build"
INTEL_DIR="$BUILD_DIR/bin/macOS/intel"
ARM_DIR="$BUILD_DIR/bin/macOS/arm"
UNIVERSAL_DIR="$BUILD_DIR/bin/macOS/universal"
IMAGES_DIR="./images"

echo "Cleaning old builds"
rm -rf "$INTEL_DIR" "$ARM_DIR" "$UNIVERSAL_DIR"
mkdir -p "$INTEL_DIR" "$ARM_DIR" "$UNIVERSAL_DIR"

echo "Building Intel (amd64)"
GOOS=darwin GOARCH=amd64 wails build -o "$APP_NAME"
mv "$BUILD_DIR/bin/$APP_NAME.app" "$INTEL_DIR/$APP_NAME-intel.app"

echo "Building Apple Silicon (arm64)"
GOOS=darwin GOARCH=arm64 wails build -o "$APP_NAME"
mv "$BUILD_DIR/bin/$APP_NAME.app" "$ARM_DIR/$APP_NAME-arm64.app"

echo "Creating universal app"
cp -R "$ARM_DIR/$APP_NAME-arm64.app" "$UNIVERSAL_DIR/$APP_NAME.app"
lipo -create \
  "$INTEL_DIR/$APP_NAME-intel.app/Contents/MacOS/$APP_NAME" \
  "$ARM_DIR/$APP_NAME-arm64.app/Contents/MacOS/$APP_NAME" \
  -output "$UNIVERSAL_DIR/$APP_NAME.app/Contents/MacOS/$APP_NAME"

echo "Universal build created at: $UNIVERSAL_DIR/$APP_NAME.app"

echo "Copying Images folder..."
cp -R "$IMAGES_DIR" "$UNIVERSAL_DIR/$APP_NAME.app/Contents/MacOS/"

echo "Signing app..."
codesign --force --deep --sign - "$UNIVERSAL_DIR/$APP_NAME.app"

lipo -info "$UNIVERSAL_DIR/$APP_NAME.app/Contents/MacOS/$APP_NAME"
codesign --verify --deep --verbose=2 "$UNIVERSAL_DIR/$APP_NAME.app"
spctl -a -vv "$UNIVERSAL_DIR/$APP_NAME.app"


