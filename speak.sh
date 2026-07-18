#!/bin/sh

gaia --speak "$1" \
| base64 -d \
| ffmpeg -f s16le \
         -ar 24000 \
         -ac 1 \
         -i pipe:0 \
         -c:a libopus \
         "${2}.opus"
