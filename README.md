# Gaia AI Agent

Gaia is a lightweight tool for bringing Google Gemini to the command line.

## Installation

1. `git clone https://github.com/Starscade/Gaia.git`
2. `cd ./Gaia`
3. `make`
> [!NOTE]
> Gaia installs to `~/.local/bin` by default. You may change this in the `Makefile`.

## Configuration

Set `export GAIA_AGENT_KEY="your-gemini-key"` or use a `.env`.
> [!TIP]
> You can create a template of your current settings with: `gaia --print-env > .env`. (This includes internal defaults.)

## Usage

- `cat foo.js bar.js baz.js | gaia "Rewrite these in Go." > main.go && go run .`
- `gaia "What's my IP address?" | sh`

> [!CAUTION]
> A safer way to execute responses is to verify their output before piping (e.g. `gaia "What's my IP address?"`), then use `gaia --echo | sh` to execute the most recent response verbatim.

> [!IMPORTANT]
> Gaia does *not* preserve context between query/response pairs. To continue a conversation, use the `--related` flag:

```
gaia "My favorite flavor is vanilla."
gaia --related "What's my favorite flavor?"
```

> [!TIP]
> If you prefer using Gaia this way, create an alias like: `alias ai='gaia --related'`.
