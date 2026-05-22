#!/bin/sh


LATEST_MINOR_VERSION="$(git tag | tail -n 1 | cut -d . -f 1,2)"
ROOT_VERSION="$(git tag | grep "$LATEST_MINOR_VERSION" | head -n 1)"
CURRENT_PATCH="$(git rev-list --count "$ROOT_VERSION"..HEAD)"
GIT_BRANCH="$(git branch | head -n 1 | cut -d' ' -f 2)"
NEW_VERSION="${LATEST_MINOR_VERSION}.${CURRENT_PATCH} (${GIT_BRANCH})"


# If no change, don't auto-increment.

test -z "$(git status --porcelain)" \
	|| sed -i \
		"s/v[[:digit:]]\.[[:digit:]]\.[[:digit:]] \([^\"]*\)/$NEW_VERSION/" \
		internal/text/text.go
