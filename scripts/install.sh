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

( ./scripts/auto-version.sh \
	&& deno task make \
	&& install -m 755 _gaia "${INSTALL_DIR}/${INSTALL_NAME}" \
	&& rm -v _gaia \
	&& print_status
) || print_status 1
