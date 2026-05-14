# Gaia AI Agent

Gaia is a lightweight CLI interface for Google Gemini with an emphasis on code generation and execution directly in the terminal.

## Installation

1. `git clone https://github.com/Starscade/Gaia.git`.
2. `cd ./Gaia`.
3. `make`.
> [!INFO]
> Gaia installs to `~/.local/bin` by default. You may change this in the `Makefile`.

## Configuration

`export GAIA_AGENT_KEY="your-gemini-key"`

## Usage

### Direct Execution

`gaia -x "What's my IP address?"`
> [!CAUTION]
> Use with extreme discretion!

### UNIX Pipe Support

`cat foo.js | gaia "Rewrite this in Go." > foo.go`
