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

# --- v1.0.0: Feature A ---
git checkout -b feature/a
touch feature-a.txt
git add feature-a.txt
git commit -m "Add Feature A"
git checkout main
git merge --no-ff feature/a -m "Merge pull request #1 from user/feature/a

## 📝 Description
Added Feature A.

## 📣 Customer-Facing Release Notes
- Released Feature A.

## ⚙️ Configuration Changes
- Added config for Feature A.

## 🔌 Required Hardware Changes
- Requires Hardware A.
"
git tag v1.0.0

# --- v1.1.0: Feature B ---
git checkout -b feature/b
touch feature-b.txt
git add feature-b.txt
git commit -m "Add Feature B"
git checkout main
git merge --no-ff feature/b -m "Merge pull request #2 from user/feature/b

## 📝 Description
Added Feature B.

## 📣 Customer-Facing Release Notes
- Released Feature B.

## ⚙️ Configuration Changes
- Added config for Feature B.

## 🔌 Required Hardware Changes
- Requires Hardware B.
"
git tag v1.1.0

# --- v1.2.0: Feature C ---
git checkout -b feature/c
touch feature-c.txt
git add feature-c.txt
git commit -m "Add Feature C"
git checkout main
git merge --no-ff feature/c -m "Merge pull request #3 from user/feature/c

## 📝 Description
Added Feature C.

## 📣 Customer-Facing Release Notes
- Released Feature C.

## ⚙️ Configuration Changes
- Added config for Feature C.

## 🔌 Required Hardware Changes
- Requires Hardware C.
"
git tag v1.2.0

# Run the tool: v1.0.0 to v1.2.0
echo "---------------------------------------------------"
echo "Running semantic-release-notes from v1.0.0 to v1.2.0..."
/Users/jpollak/semantic-release-notes/semantic-release-notes -repo "$TEST_REPO" -from v1.0.0 -to v1.2.0

# Run the tool: v1.1.0 to v1.2.0
echo "---------------------------------------------------"
echo "Running semantic-release-notes from v1.1.0 to v1.2.0..."
/Users/jpollak/semantic-release-notes/semantic-release-notes -repo "$TEST_REPO" -from v1.1.0 -to v1.2.0

# Run the tool: v1.0.0 (everything up to v1.0.0)
echo "---------------------------------------------------"
echo "Running semantic-release-notes for v1.0.0..."
/Users/jpollak/semantic-release-notes/semantic-release-notes -repo "$TEST_REPO" -to v1.0.0

# Run the tool: Missing tag (should fail)
echo "---------------------------------------------------"
echo "Running semantic-release-notes with missing tag (expecting failure)..."
if /Users/jpollak/semantic-release-notes/semantic-release-notes -repo "$TEST_REPO" -from v9.9.9 -to v1.0.0; then
    echo "Error: Tool should have failed for missing tag v9.9.9"
    exit 1
else
    echo "Success: Tool failed as expected for missing tag."
fi

# --- v1.3.0: Cleanup and XML comments ---
git checkout -b chore/cleanup-xml
touch cleanup.txt
git add cleanup.txt
git commit -m "Cleanup with XML comments"
git checkout main
git merge --no-ff chore/cleanup-xml -m "Merge pull request #4 from user/chore/cleanup-xml

## 📝 Description
Cleanup with XML comments.

## 📣 Customer-Facing Release Notes
<!-- This is a comment -->
- Cleaned up some stuff.
<!-- Another comment -->

## ⚙️ Configuration Changes
None

## 🔌 Required Hardware Changes
N/A
"
git tag v1.3.0

# Run the tool: v1.2.0 to v1.3.0 (Testing cleanup)
echo "---------------------------------------------------"
echo "Running semantic-release-notes from v1.2.0 to v1.3.0 (Testing cleanup)..."
/Users/jpollak/semantic-release-notes/semantic-release-notes -repo "$TEST_REPO" -from v1.2.0 -to v1.3.0

# --- v1.4.0: Multiline Description ---
git checkout -b feature/multiline
touch multiline.txt
git add multiline.txt
git commit -m "Add multiline feature"
git checkout main
git merge --no-ff feature/multiline -m "Merge pull request #5 from user/feature/multiline

## 📝 Description
Multiline Feature.
This is a detailed description.
It has multiple lines.

## 📣 Customer-Facing Release Notes
- Added multiline feature.
"
git tag v1.4.0

# Run the tool: v1.3.0 to v1.4.0 (Testing multiline)
echo "---------------------------------------------------"
echo "Running semantic-release-notes from v1.3.0 to v1.4.0 (Testing multiline)..."
/Users/jpollak/semantic-release-notes/semantic-release-notes -repo "$TEST_REPO" -from v1.3.0 -to v1.4.0

# Clean up
rm -rf "$TEST_REPO"
