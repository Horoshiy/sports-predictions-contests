#!/bin/bash
set -e

echo "Installing Playwright browsers..."
cd "$(dirname "$0")/../frontend"

# Install dependencies if needed
if [ ! -d "node_modules" ]; then
  echo "Installing npm dependencies..."
  npm install
fi

# Install Playwright browsers
echo "Installing Playwright browsers (chromium, firefox, webkit)..."
npx playwright install --with-deps

echo "✅ Playwright browsers installed successfully!"
