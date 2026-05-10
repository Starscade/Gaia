.PHONY: all fmt

export PATH := $(HOME)/.local/bin:$(PATH)


all:

	@\
		make fmt \
		&& CGO_ENABLED=0 \
		go build -ldflags="-s -w" -v -x \
		-o ~/.local/bin/gaia ./src \
		&& printf "\n \033[1;32mOK\033[0m\n" \
		|| printf "\n \033[1;31mERR\033[0m\n"


fmt:

	@go fmt src/*.go


