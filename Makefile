.PHONY: all clean fmt test

export CGO_ENABLED := 0
export PATH := $(HOME)/.local/bin:$(PATH)


all:

	@ \
		go mod tidy \
		&& make fmt \
		&& go build -ldflags="-s -w" -v -x \
		-o ~/.local/bin/gaia ./cmd/gaia \
		&& ( \
			./scripts/auto-version.sh \
			&& printf "\n \033[1;32mOK\033[0m\n\n" \
		) \
		|| printf "\n \033[1;31mERR\033[0m\n\n"


clean:

	go mod tidy \
	&& go clean -cache -x


fmt:

	@go fmt ./...


test:

	@go test -v -x ./...
