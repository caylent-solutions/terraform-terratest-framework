#!/bin/bash
set -e

# This script bumps the version based on conventional commit messages
# Usage: ./scripts/bump-version.sh [major|minor|patch]

# Get the current version from VERSION file
CURRENT_VERSION=$(cat VERSION)
echo "Current version: $CURRENT_VERSION"

# Parse the version
MAJOR=$(echo $CURRENT_VERSION | cut -d. -f1)
MINOR=$(echo $CURRENT_VERSION | cut -d. -f2)
PATCH=$(echo $CURRENT_VERSION | cut -d. -f3)

# Determine version bump type from commit messages if not specified
if [ -z "$1" ]; then
  # Get commit messages since last tag
  LAST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "")
  
  if [ -z "$LAST_TAG" ]; then
    # If no tags exist, get all commits
    COMMITS=$(git log --pretty=format:"%s")
  else
    # Get commits since last tag
    COMMITS=$(git log ${LAST_TAG}..HEAD --pretty=format:"%s")
  fi
  
  # Check for breaking changes
  if echo "$COMMITS" | grep -q "^BREAKING CHANGE:" || \
     echo "$COMMITS" | grep -q "^breaking!:" || \
     echo "$COMMITS" | grep -q "!:" || \
     echo "$COMMITS" | grep -q "feat!:" || \
     echo "$COMMITS" | grep -q "fix!:" || \
     echo "$COMMITS" | grep -q "refactor!:" || \
     echo "$COMMITS" | grep -q "docs!:" || \
     echo "$COMMITS" | grep -q "style!:" || \
     echo "$COMMITS" | grep -q "test!:" || \
     echo "$COMMITS" | grep -q "chore!:" || \
     echo "$COMMITS" | grep -q "ci!:" || \
     echo "$COMMITS" | grep -q "build!:" || \
     echo "$COMMITS" | grep -q "perf!:"; then
    BUMP_TYPE="major"
  # Check for feature commits
  elif echo "$COMMITS" | grep -q "^feat:" || \
       echo "$COMMITS" | grep -q "^feature:"; then
    BUMP_TYPE="minor"
  # All other commit types result in a patch bump
  else
    BUMP_TYPE="patch"
  fi
else
  BUMP_TYPE=$1
fi

echo "Bump type: $BUMP_TYPE"

# Bump version
case $BUMP_TYPE in
  major)
    MAJOR=$((MAJOR + 1))
    MINOR=0
    PATCH=0
    ;;
  minor)
    MINOR=$((MINOR + 1))
    PATCH=0
    ;;
  patch)
    PATCH=$((PATCH + 1))
    ;;
  *)
    echo "Invalid bump type: $BUMP_TYPE"
    echo "Usage: $0 [major|minor|patch]"
    exit 1
    ;;
esac

# Create new version
NEW_VERSION="${MAJOR}.${MINOR}.${PATCH}"
echo "New version: $NEW_VERSION"

# Update VERSION file
echo $NEW_VERSION > VERSION

# Update version in other files
sed -i "s/Version=v[0-9]*\.[0-9]*\.[0-9]*/Version=v${NEW_VERSION}/g" Makefile

echo "Version bumped to $NEW_VERSION"

# Stage the changes
git add VERSION Makefile

# Commit the changes
git commit -m "chore: bump version to $NEW_VERSION"

# Create a tag
git tag -a "v$NEW_VERSION" -m "Version $NEW_VERSION"

echo "Changes committed and tagged as v$NEW_VERSION"
echo "Run 'git push && git push --tags' to push changes to remote"