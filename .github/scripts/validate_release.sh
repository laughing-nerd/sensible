#!/bin/bash
set -e

# Remove 'refs/tags/' prefix
TAG_NAME="${GITHUB_REF#refs/tags/}"
echo "New tag: $TAG_NAME"

# Get the latest release tag (ignores drafts and prereleases)
LATEST_RELEASE=$(curl -s https://api.github.com/repos/${GITHUB_REPOSITORY}/releases/latest | jq -r .tag_name)

# if no release exists
if [[ "$LATEST_RELEASE" == "null" || -z "$LATEST_RELEASE" ]]; then
  echo "✅ No prior release found, proceeding with release of $TAG_NAME"
  exit 0
fi

echo "Latest release: $LATEST_RELEASE"

# Compare version using sort -V
if [ "$(printf '%s\n' "$TAG_NAME" "$LATEST_RELEASE" | sort -V | head -n1)" != "$LATEST_RELEASE" ]; then
  echo "❌ Error: New tag $TAG_NAME is not greater than the current release $LATEST_RELEASE"
  exit 1
fi

echo "✅ New tag $TAG_NAME is newer than the current release $LATEST_RELEASE"
