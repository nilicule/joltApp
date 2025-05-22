#!/bin/bash

# Build script for Jolt

# Ensure the script exits on any error
set -e

# Create build directory if it doesn't exist
mkdir -p build

# Build the application
echo "Building Jolt..."
go build -o build/Jolt

# Create the application bundle structure
echo "Creating application bundle..."
mkdir -p build/Jolt.app/Contents/{MacOS,Resources}

# Copy the executable to the bundle
cp build/Jolt build/Jolt.app/Contents/MacOS/

# Get version from version.txt file
VERSION=$(cat version.txt)

# Create Info.plist to hide dock icon
cat > build/Jolt.app/Contents/Info.plist << EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleName</key>
    <string>Jolt</string>
    <key>CFBundleDisplayName</key>
    <string>Jolt</string>
    <key>CFBundleIdentifier</key>
    <string>org.rc6.jolt</string>
    <key>CFBundleVersion</key>
    <string>${VERSION}</string>
    <key>CFBundleShortVersionString</key>
    <string>${VERSION}</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>CFBundleSignature</key>
    <string>????</string>
    <key>CFBundleExecutable</key>
    <string>Jolt</string>
    <key>LSMinimumSystemVersion</key>
    <string>10.13</string>
    <key>LSUIElement</key>
    <true/>
</dict>
</plist>
EOF

# Create releases directory if it doesn't exist
mkdir -p releases

# Get version from version.txt file
VERSION=$(cat version.txt)
ZIP_FILENAME="Jolt-${VERSION}.zip"

# Create zip file of the application bundle
echo "Creating zip file..."
cd build
zip -r "${ZIP_FILENAME}" Jolt.app
cd ..

# Move zip file to releases directory
echo "Moving zip file to releases directory..."
mv "build/${ZIP_FILENAME}" releases/

echo "Application bundle created: build/Jolt.app"
echo "Zip file created: releases/${ZIP_FILENAME}"
echo "You can now run the application by double-clicking on build/Jolt.app"
