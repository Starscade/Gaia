#!/bin/sh

LATEST_MINOR_VERSION="$(git tag | tail -n 1 | cut -d . -f 1,2)"
ROOT_VERSION="$(git tag | grep "$LATEST_MINOR_VERSION" | head -n 1)"
CURRENT_PATCH="$(git rev-list --count "$ROOT_VERSION"..HEAD)"
NEXT_PATCH="$((1 + $CURRENT_PATCH))"
GIT_BRANCH="$(git branch | head -n 1 | cut -d' ' -f 2)"
NEW_VERSION="${LATEST_MINOR_VERSION}.${NEXT_PATCH} (${GIT_BRANCH})"
CURRENT_VERSION="$(gaia --version)"

# If no change, don't auto-increment.

if test "$CURRENT_VERSION" != "$NEW_VERSION"; then
	sed -i \
	"s/v[[:digit:]]\+\.[[:digit:]]\+\.[[:digit:]]\+ ([^\']*)/$NEW_VERSION/" \
	./main.ts
fi
