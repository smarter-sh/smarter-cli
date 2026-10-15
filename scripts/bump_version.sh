#!/usr/bin/env bash
#------------------------------------------------------------------------------
# Script to bump version numbers in various files.
#
# semantic-release runs it (release.config.js) on the CI's Ubuntu runner, and
# maintainers may run it on macOS, so it edits files with perl -pi, which
# behaves the same on both; sed -i doesn't (BSD sed requires an argument that
# GNU sed reads as the script).
#------------------------------------------------------------------------------
set -euo pipefail

VERSION=${1:-}
if [ -z "$VERSION" ]; then
  echo "Usage: $0 <new-version>"
  exit 1
fi
export VERSION

# Update VERSION file
echo "$VERSION" > VERSION
echo "Updated VERSION file to $VERSION"

# Update VERIFICATION.txt
perl -pi -e 's/v[0-9]+\.[0-9]+\.[0-9]+/v$ENV{VERSION}/g; s/(?<![v0-9.])[0-9]+\.[0-9]+\.[0-9]+/$ENV{VERSION}/g' VERIFICATION.txt
echo "Updated VERIFICATION.txt to $VERSION"

# Update Makefile choco-push target
perl -pi -e 's/(choco push smarter\.)[0-9]+\.[0-9]+\.[0-9]+(\.nupkg)/${1}$ENV{VERSION}${2}/g' Makefile
echo "Updated Makefile choco-push target to $VERSION"

# Update debian/rules download URL and filename
perl -pi -e 's|(releases/download/)v[0-9]+\.[0-9]+\.[0-9]+|${1}v$ENV{VERSION}|g; s/smarter-ubuntu-latest-[0-9]+\.[0-9]+\.[0-9]+/smarter-ubuntu-latest-$ENV{VERSION}/g' debian/rules
echo "Updated debian/rules to $VERSION"
