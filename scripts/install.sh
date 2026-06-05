#!/bin/sh

INSTALL_DIR="${GAIA_INSTALL_DIR:-$HOME/.local/bin}"
INSTALL_NAME=gaia

test -n "$1" && INSTALL_NAME="$1"

print_status() {
	COLOR=2
	STATUS=OK
	if test -n "$1"; then
		COLOR=1
		STATUS=ERR
	fi
	printf "\n \033[1;3${COLOR}m${STATUS}\033[0m\n\n"
}

( go mod tidy \
	&& go fmt ./... \
	&& CGO_ENABLED=0 \
		go build \
		-ldflags="-s -w" \
		-v -x -o \
		"${INSTALL_DIR}/${INSTALL_NAME}" "./cmd/${INSTALL_NAME}" \
	&& ./scripts/auto-version.sh \
	&& print_status
) || print_status 1
