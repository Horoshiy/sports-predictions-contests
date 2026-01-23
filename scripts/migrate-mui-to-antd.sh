#!/bin/bash

# Bulk migration script for MUI to Ant Design
# This script performs common replacements across all component files

FRONTEND_DIR="frontend/src"

echo "Starting MUI to Ant Design migration..."

# Find all TypeScript/TSX files with MUI imports
FILES=$(grep -rl "@mui" "$FRONTEND_DIR" --include="*.tsx" --include="*.ts" 2>/dev/null)

for file in $FILES; do
  echo "Processing: $file"
  
  # Backup original file
  cp "$file" "$file.bak"
  
  # Common import replacements
  sed -i.tmp 's/@mui\/material/@antd/g' "$file"
  sed -i.tmp 's/@mui\/icons-material/@ant-design\/icons/g' "$file"
  sed -i.tmp 's/@mui\/x-date-pickers/dayjs/g' "$file"
  
  # Common component replacements
  sed -i.tmp 's/Box/Space/g' "$file"
  sed -i.tmp 's/Typography/Text/g' "$file"
  sed -i.tmp 's/Chip/Tag/g' "$file"
  sed -i.tmp 's/CircularProgress/Spin/g' "$file"
  sed -i.tmp 's/IconButton/Button/g' "$file"
  sed -i.tmp 's/TextField/Input/g' "$file"
  sed -i.tmp 's/Dialog/Modal/g' "$file"
  
  # Clean up temp files
  rm -f "$file.tmp"
done

echo "Migration complete. Backup files created with .bak extension"
echo "Please review changes and test thoroughly"
