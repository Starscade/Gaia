#!/bin/sh

INSTALL_DIR="${GAIA_INSTALL_DIR:-${HOME}/.local/bin}"
INSTALL_NAME=gaia

check_command() {
	CMD_NAME="$1"
	command -v "$CMD_NAME" > /dev/null 2>&1 \
		|| {
			printf " \033[1m${CMD_NAME}\033[0m not found. Is it installed?\n" \
			&& exit
		}
}

check_command deno
check_command make

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

( deno task make \
	&& mkdir -p "$INSTALL_DIR" \
	&& install -m 755 _gaia "${INSTALL_DIR}/${INSTALL_NAME}" \
	&& rm -v _gaia \
	&& print_status
) || print_status 1
