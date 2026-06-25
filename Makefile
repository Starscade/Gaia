.POSIX:


all:

	@./scripts/install.sh


check:

	@deno task lint


format:

	@deno task format


www:

	@deno task serve


