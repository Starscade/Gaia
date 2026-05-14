# Gaia AI Agent

Gaia is a lightweight CLI tool for Google Gemini with an emphasis on code generation and execution.

## Installation

1. `git clone https://github.com/Starscade/Gaia.git`
2. `cd ./Gaia`
3. `make`
> [!NOTE]
> Gaia installs to `~/.local/bin` by default. You may change this in the `Makefile`.

## Configuration

Set your environment one variable at a time (e.g. `export GAIA_AGENT_KEY="your-gemini-key"`) or use a `.env`.
> [!TIP]
> You can create a template of your current settings with: `gaia --print-env > .env`. (This will only print variables that beginn with `GAIA_`.)

## Usage

### Direct Execution

`gaia -x "What's my IP address?"`
> [!CAUTION]
> Use with extreme discretion!

### UNIX Pipe Support

`cat foo.js | gaia "Rewrite this in Go." > foo.go`
