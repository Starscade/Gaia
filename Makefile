.POSIX:


all:

	@./install.sh


check:

	@deno task lint


format:

	@deno task format


www:

	@deno task serve


