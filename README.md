# Gaia AI Agent

Gaia is a lightweight CLI tool for Google Gemini tailored for raw code generation and execution.

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

Gaia within a pipe: `cat foo.js bar.js baz.js | gaia "Rewrite these in Go." > main.go && go run .`.

Gaia direct: `gaia "What's my IP address?" | sh`.
> [!CAUTION]
> A safer way to execute responses is to verify their output before piping (e.g. `gaia "What's my IP address?"`), then use `gaia --replay | sh` to execute the most recent response verbatim.

> [!IMPORTANT]
> Gaia does *not* preserve context between query/response pairs. To continue a conversation, use the `--related` flag:

```
gaia "My favorite flavor is vanilla."
gaia --related "What's my favorite flavor?"
```

> [!TIP]
> If you plan on using Gaia this way, aliases are your friend: `alias ai='gaia --related'`.
