#!/bin/sh

export CGO_ENABLED=0
INSTALL_DIR="${GAIA_INSTALL_DIR:-$HOME/.local/bin}"

go mod tidy \
	&& go fmt ./... \
	&& go build \
		-ldflags="-s -w" \
		-v -x -o \
		"${INSTALL_DIR}/gaia" ./cmd/gaia \
		&& ./scripts/auto-version.sh \
		&& printf "\n \033[1;32mOK\033[0m\n\n" \
		|| printf "\n \033[1;31mERR\033[0m\n\n"
