.POSIX:


all:

	@./scripts/install.sh


clean:

	@go mod tidy
	@go clean -cache -x


fmt:

	@go fmt ./...


test:

	@go test -v -x ./...


