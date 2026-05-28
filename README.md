# Gaia AI Agent

Gaia is a headless client for Google Gemini.

## Installation

1. `git clone https://github.com/Starscade/Gaia.git`
2. `cd ./Gaia`
3. `make`

> [!NOTE]
> Gaia installs to `~/.local/bin` by default. You may change this by setting `GAIA_INSTALL_DIR`.

## Configuration

Set `export GAIA_AGENT_KEY="your-gemini-key"` or use a `.env`.

> [!TIP]
> You can create a template of your current settings with: `gaia --print-env > .env`. (This includes internal defaults.)

## Usage

###### DIRECT EXECUTION
`gaia "What's my IP address?" | sh`

###### FILE ANALYSIS
`gaia --read README.md "Summarize this..."`

###### MULTI-FILE / FOLDER ANALYSIS
`gaia --read ./{foo,bar}/*.md "Give me a project overview."`

###### UNIX PIPELINES
`cat foo.js bar.js baz.js | gaia "Rewrite these in Go." > main.go && go run .`

> [!CAUTION]
> A safer way to execute responses directly is to verify their output before piping (e.g. `gaia "What's my IP address?"`), then use `gaia --echo | sh` to execute the most recent response verbatim.

> [!IMPORTANT]
> Gaia does *not* preserve context between query/response pairs. To continue a conversation, use the `--related` flag:

```
gaia "I like root beer."
gaia --related "Do I like root beer?"
```

> [!TIP]
> If you prefer using Gaia this way, you can make the behaviour permenant by setting an alias: `alias ai='gaia --related'`.
