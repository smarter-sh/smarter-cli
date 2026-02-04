#!/usr/bin/env bash
#------------------------------------------------------------------------------
# Script to bump version numbers in various files.
#------------------------------------------------------------------------------

VERSION=$1
if [ -z "$VERSION" ]; then
  echo "Usage: $0 <new-version>"
  exit 1
fi

# Update VERSION file
echo "$VERSION" > VERSION
echo "Updated VERSION file to $VERSION"

# Update VERIFICATION.txt
sed -i '' -E "s/v[0-9]+\.[0-9]+\.[0-9]+/v$VERSION/g; s/[0-9]+\.[0-9]+\.[0-9]+/$VERSION/g" VERIFICATION.txt
echo "Updated VERIFICATION.txt to $VERSION"


# Update Makefile choco-push target
sed -i '' -E "s/(choco push smarter\\.)[0-9]+\\.[0-9]+\\.[0-9]+(\\.nupkg --source https:\/\/push.chocolatey.org\/)\\//\\1$VERSION\\2\\//g" Makefile
echo "Updated Makefile choco-push target to $VERSION"



# Update debian/rules download URL and filename
sed -i '' -E "s|(releases/download/)v[0-9]+\.[0-9]+\.[0-9]+|\1v$VERSION|g; s/smarter-ubuntu-latest-[0-9]+\.[0-9]+\.[0-9]+/smarter-ubuntu-latest-$VERSION/g" debian/rules
echo "Updated debian/rules to $VERSION"
