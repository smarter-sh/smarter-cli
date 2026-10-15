#!/usr/bin/env bash
# Copy the Smarter platform's published OpenAPI document into the cli's test data.
#
# The platform commits its api's OpenAPI document, generated from its Pydantic
# models. The cli's contract test (contract_test.go) runs every command against
# a mock platform, and checks that every request is an operation of this copy,
# that every /cli/ operation has a command, and that the kinds the cli sends are
# the document's SAMKinds. See doc/adr/022-pinned-openapi.md.
#
# Run this after the platform's api changes (make contracts), then run the tests. It records
# the platform version it copied in testdata/contracts/SOURCE. Pin a release: the cli on main
# should be checked against the platform's latest stable release.
#
# usage, from the repository root:
#   scripts/sync-contracts.sh [path of a smarter checkout, default ../smarter] [git ref]
#
# With a git ref (e.g. origin/main, or a release tag), the file is read at that ref with git show,
# without changing the checkout. Without one, it is copied from the checkout's working tree.
set -euo pipefail

SMARTER="${1:-../smarter}"
REF="${2:-}"
TARGET="testdata/contracts"
SOURCE="smarter/smarter/apps/api/v1/contract/openapi.json"

if [ ! -d "$SMARTER/smarter/smarter/apps/api/v1/contract" ]; then
  echo "No smarter checkout at $SMARTER" >&2
  exit 1
fi

mkdir -p "$TARGET"
if [ -n "$REF" ]; then
  git -C "$SMARTER" show "$REF:$SOURCE" > "$TARGET/openapi.json"
else
  cp "$SMARTER/$SOURCE" "$TARGET/openapi.json"
fi

VERSION=$(git -C "$SMARTER" describe --tags --always "${REF:-HEAD}" 2>/dev/null || echo "$SMARTER")
if [ -z "$REF" ] && [ -n "$(git -C "$SMARTER" status --porcelain -- "$SOURCE" 2>/dev/null)" ]; then
  VERSION="$VERSION, with uncommitted changes"
fi
# record where the copy came from; commit it with the copy.
echo "openapi.json: smarter $VERSION" > "$TARGET/SOURCE"
echo "Copied the OpenAPI document of $VERSION into $TARGET"
