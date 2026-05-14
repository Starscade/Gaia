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
> You can create a template of your current settings with: `gaia --print-env > .env`. (This includes internal defaults.)

## Usage

###### DIRECT EXECUTION

`gaia "What's my IP address?" | sh`
> [!CAUTION]
> Execute with extreme discretion! A safer way to run responses as code is to verify the output of `gaia "What's my IP address?"` without piping, then use `gaia --recall | sh` to execute it verbatim.

###### WITHIN A PIPE

`cat foo.js | gaia "Rewrite this in Go." > main.go && go run .`
