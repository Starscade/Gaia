#!/bin/sh

gaia --draw "$1" \
| base64 -d \
> "${2}.jpg"
