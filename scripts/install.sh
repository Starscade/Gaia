#!/bin/sh

export CGO_ENABLED=0
INSTALL_DIR="${GAIA_INSTALL_DIR:-$HOME/.local/bin}"

go mod tidy \
	&& go fmt ./... \
	&& go build \
		-ldflags="-s -w" \
		-v -x -o \
		"${INSTALL_DIR}/gaia" ./cmd/gaia \
		&& ./scripts/auto-version.sh
