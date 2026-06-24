#!/bin/sh

find "$(pwd)/cmd/gaia" \
     "$(pwd)/internal" \
     -type f \
     -name "*" \
     -exec cat {} + \
| ~/.local/bin/gaia "$@"
