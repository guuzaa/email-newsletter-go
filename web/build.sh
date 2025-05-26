#!/bin/bash

# Build script for the newsletter frontend
echo "🏗️  Building Newsletter Frontend..."

# Check if bun is installed
if ! command -v bun &> /dev/null; then
    echo "❌ Bun is not installed. Please install it from https://bun.sh"
    exit 1
fi

# Change to web directory
cd "$(dirname "$0")"

# Install dependencies if node_modules doesn't exist
if [ ! -d "node_modules" ]; then
    echo "📦 Installing dependencies..."
    bun install
fi

# Build the project
echo "🔨 Building for production..."
bun run build

echo "🎉 Build complete! Built files are in dist/ directory and ready for Go embedding."
