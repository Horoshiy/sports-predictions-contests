#!/bin/bash

# Script to generate screenshots for documentation
# This script uses Playwright to capture screenshots of the application

set -e

echo "🖼️  Generating screenshots for documentation..."

# Check if frontend directory exists
if [ ! -d "frontend" ]; then
    echo "❌ Error: frontend directory not found"
    exit 1
fi

# Create screenshots directory if it doesn't exist
mkdir -p docs/screenshots

cd frontend

# Check if node_modules exists
if [ ! -d "node_modules" ]; then
    echo "📦 Installing dependencies..."
    npm install
fi

# Check if Playwright browsers are installed
if [ ! -d "node_modules/@playwright/test" ]; then
    echo "🎭 Installing Playwright..."
    npm install --save-dev @playwright/test
fi

echo "🎭 Installing Playwright browsers..."
npx playwright install --with-deps

# Create a temporary screenshot script
cat > /tmp/generate-screenshots.ts << 'EOF'
import { chromium } from '@playwright/test';
import * as path from 'path';

async function generateScreenshots() {
  const browser = await chromium.launch({ headless: true });
  const context = await browser.newContext({
    viewport: { width: 1920, height: 1080 }
  });
  const page = await context.newPage();

  const baseURL = 'http://localhost:3000';
  const screenshotsDir = path.join(__dirname, '..', 'docs', 'screenshots');

  console.log('📸 Capturing screenshots...');

  try {
    // Login page
    console.log('  - Login page');
    await page.goto(`${baseURL}/login`);
    await page.waitForLoadState('networkidle');
    await page.screenshot({ path: `${screenshotsDir}/login-page.png`, fullPage: true });

    // Register page
    console.log('  - Register page');
    await page.goto(`${baseURL}/register`);
    await page.waitForLoadState('networkidle');
    await page.screenshot({ path: `${screenshotsDir}/register-page.png`, fullPage: true });

    // Login with test credentials
    await page.goto(`${baseURL}/login`);
    await page.fill('input[type="email"]', 'user1@example.com');
    await page.fill('input[type="password"]', 'password123');
    await page.click('button[type="submit"]');
    await page.waitForURL(`${baseURL}/contests`);

    // Contests list
    console.log('  - Contests list');
    await page.goto(`${baseURL}/contests`);
    await page.waitForLoadState('networkidle');
    await page.screenshot({ path: `${screenshotsDir}/contests-list.png`, fullPage: true });

    // Contest details (if available)
    console.log('  - Contest details');
    const contestLink = await page.locator('a[href*="/contests/"]').first();
    if (await contestLink.count() > 0) {
      await contestLink.click();
      await page.waitForLoadState('networkidle');
      await page.screenshot({ path: `${screenshotsDir}/contest-details.png`, fullPage: true });
    }

    // Predictions page
    console.log('  - Predictions page');
    await page.goto(`${baseURL}/predictions`);
    await page.waitForLoadState('networkidle');
    await page.screenshot({ path: `${screenshotsDir}/predictions-page.png`, fullPage: true });

    // Leaderboard
    console.log('  - Leaderboard');
    await page.goto(`${baseURL}/contests`);
    await page.waitForLoadState('networkidle');
    const leaderboardTab = await page.locator('text=Leaderboard').first();
    if (await leaderboardTab.count() > 0) {
      await leaderboardTab.click();
      await page.waitForTimeout(1000);
      await page.screenshot({ path: `${screenshotsDir}/leaderboard.png`, fullPage: true });
    }

    // Profile page
    console.log('  - Profile page');
    await page.goto(`${baseURL}/profile`);
    await page.waitForLoadState('networkidle');
    await page.screenshot({ path: `${screenshotsDir}/profile-page.png`, fullPage: true });

    // Analytics dashboard
    console.log('  - Analytics dashboard');
    await page.goto(`${baseURL}/analytics`);
    await page.waitForLoadState('networkidle');
    await page.screenshot({ path: `${screenshotsDir}/analytics-dashboard.png`, fullPage: true });

    // Sports management
    console.log('  - Sports management');
    await page.goto(`${baseURL}/sports`);
    await page.waitForLoadState('networkidle');
    await page.screenshot({ path: `${screenshotsDir}/sports-management.png`, fullPage: true });

    // Teams page
    console.log('  - Teams page');
    await page.goto(`${baseURL}/teams`);
    await page.waitForLoadState('networkidle');
    await page.screenshot({ path: `${screenshotsDir}/teams-page.png`, fullPage: true });

    console.log('✅ All screenshots captured successfully!');
  } catch (error) {
    console.error('❌ Error capturing screenshots:', error);
    throw error;
  } finally {
    await browser.close();
  }
}

generateScreenshots().catch(console.error);
EOF

# Check if services are running
echo "🔍 Checking if services are running..."
if ! curl -s http://localhost:3000 > /dev/null; then
    echo "⚠️  Warning: Frontend is not running on http://localhost:3000"
    echo "   Please start the services first:"
    echo "   make dev"
    echo "   make docker-services"
    exit 1
fi

# Run the screenshot generation script
echo "📸 Generating screenshots..."
npx tsx /tmp/generate-screenshots.ts

# Clean up
rm /tmp/generate-screenshots.ts

cd ..

echo ""
echo "✅ Screenshots generated successfully!"
echo "📁 Screenshots saved to: docs/screenshots/"
echo ""
echo "Generated screenshots:"
ls -lh docs/screenshots/*.png 2>/dev/null || echo "No screenshots found"
