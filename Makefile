export PATH := $(HOME)/.local/bin:$(PATH)


all:

	@CGO_ENABLED=0 go build -ldflags="-s -w" -v -x -o ~/.local/bin/gaia main.go


