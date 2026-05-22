.PHONY: all fmt test

export PATH := $(HOME)/.local/bin:$(PATH)


all:

	@\
		go mod tidy \
		&& make fmt \
		&& ./scripts/auto-version.sh \
		&& CGO_ENABLED=0 \
		go build -ldflags="-s -w" -v -x \
		-o ~/.local/bin/gaia ./cmd/gaia \
		&& printf "\n \033[1;32mOK\033[0m\n\n" \
		|| printf "\n \033[1;31mERR\033[0m\n\n"


fmt:

	@go fmt ./...


test:

	@go test -v -x ./...
