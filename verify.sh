#!/bin/bash
set -e

# Create a temporary directory for the test repo
TEST_REPO=$(mktemp -d)
echo "Creating test repo at $TEST_REPO"

cd "$TEST_REPO"
git init
git config user.email "you@example.com"
git config user.name "Your Name"
git config merge.log false

# Initial commit
touch README.md
git add README.md
git commit -m "Initial commit"

# --- v1.0.0: Foundation ---
git checkout -b feature/init
touch init.txt
git add init.txt
git commit -m "Initial feature"
git checkout main
git merge --no-ff feature/init -m "feat(core): Initial release [ISSUE-1]

## 📝 Description
Setting up the foundation.

## 📣 Customer-Facing Release Notes
- Initial release of the platform.
"
git tag v1.0.0

# --- v1.1.0: Features and Fixes ---
git checkout -b feature/ui-update
touch ui.txt
git add ui.txt
git commit -m "UI Update"
git checkout main
git merge --no-ff feature/ui-update -m "feat(ui): Update dashboard layout [ISSUE-2]

## 📝 Description
New dashboard with widgets.

## 📣 Customer-Facing Release Notes
- New dashboard available.
"

git checkout -b fix/api-bug
touch api.txt
git add api.txt
git commit -m "API Fix"
git checkout main
git merge --no-ff fix/api-bug -m "fix(api): Fix timeout issue [ISSUE-3]

## ⚙️ Configuration Changes
- Increased default timeout to 30s.
"
git tag v1.1.0

# --- v1.2.0: Hardware and Cleanup ---
git checkout -b feature/hardware-support
touch hw.txt
git add hw.txt
git commit -m "Hardware Support"
git checkout main
git merge --no-ff feature/hardware-support -m "feat(hw): Add support for Model X [ISSUE-4]

## 🔌 Required Hardware Changes
- Model X requires firmware v2.0.
"

git checkout -b chore/cleanup
touch cleanup.txt
git add cleanup.txt
git commit -m "Cleanup"
git checkout main
git merge --no-ff chore/cleanup -m "chore: Remove legacy code [ISSUE-5]"
git tag v1.2.0

# --- Verification ---

TOOL_PATH="/Users/jpollak/semantic-release-notes/semantic-release-notes"

echo "---------------------------------------------------"
echo "Test 1: v1.0.0 to v1.1.0 (Features, Fixes, Config)"
OUTPUT=$($TOOL_PATH -repo "$TEST_REPO" -from v1.0.0 -to v1.1.0)
echo "$OUTPUT"

if [[ "$OUTPUT" != *"# Release Notes (v1.0.0...v1.1.0)"* ]]; then
    echo "Error: Header format incorrect."
    exit 1
fi
if [[ "$OUTPUT" != *"Update dashboard layout"* ]]; then
    echo "Error: Missing UI feature."
    exit 1
fi
if [[ "$OUTPUT" != *"Fix timeout issue"* ]]; then
    echo "Error: Missing API fix."
    exit 1
fi
if [[ "$OUTPUT" != *"Configuration Changes"* ]]; then
    echo "Error: Missing Configuration Changes section."
    exit 1
fi

echo "---------------------------------------------------"
echo "Test 2: v1.1.0 to v1.2.0 (Hardware, Chore)"
OUTPUT=$($TOOL_PATH -repo "$TEST_REPO" -from v1.1.0 -to v1.2.0)
echo "$OUTPUT"

if [[ "$OUTPUT" != *"# Release Notes (v1.1.0...v1.2.0)"* ]]; then
    echo "Error: Header format incorrect."
    exit 1
fi
if [[ "$OUTPUT" != *"Add support for Model X"* ]]; then
    echo "Error: Missing hardware feature."
    exit 1
fi
if [[ "$OUTPUT" != *"Required Hardware Changes"* ]]; then
    echo "Error: Missing Required Hardware Changes section."
    exit 1
fi
# Chores usually go to "Other Changes" or similar if not mapped, or just listed if mapped to Other.
# Our parser maps unknown types to Other. 'chore' maps to Other.
if [[ "$OUTPUT" != *"Remove legacy code"* ]]; then
    echo "Error: Missing chore commit."
    exit 1
fi

echo "---------------------------------------------------"
echo "Test 3: v1.0.0 to v1.2.0 (Full Range)"
OUTPUT=$($TOOL_PATH -repo "$TEST_REPO" -from v1.0.0 -to v1.2.0)
# Should contain everything from 1.1.0 and 1.2.0
if [[ "$OUTPUT" != *"Update dashboard layout"* ]] || [[ "$OUTPUT" != *"Add support for Model X"* ]]; then
    echo "Error: Missing commits in full range test."
    exit 1
fi

# Clean up
rm -rf "$TEST_REPO"
echo "---------------------------------------------------"
echo "All tests passed!"


