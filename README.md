# Gaia AI Agent

Gaia is a portable client for Google Gemini that works on both the command line
as well as the web.

## Installation

1. `git clone https://github.com/Starscade/Gaia.git`
2. `cd ./Gaia`
3. `make`

> [!NOTE]
> Gaia installs to `~/.local/bin` by default. You may change this by setting
> `GAIA_INSTALL_DIR` before running `make`.

## Configuration

Set `export GAIA_API_KEY="your-gemini-key"` or use a `.env`.

> [!TIP]
> You can create a template of your current settings with:
> `gaia --print-env > .env`. (This includes internal defaults.)

## Usage

###### DIRECT EXECUTION

`gaia --ask "What's my IP address?" | sh`

###### FILE ANALYSIS

`gaia --read README.md --ask "Summarize this..."`

> [!CAUTION]
> A safer way to execute responses directly is to verify their output before
> piping (e.g. `gaia --ask "What's my IP address?"`), then use
> `gaia --echo | sh` to execute the most recent response verbatim.

> [!IMPORTANT]
> Gaia does _not_ preserve context between query/response pairs. To continue a
> conversation, use the `--related` flag:

```
gaia --ask "I like root beer."
gaia --related --ask "Do I like root beer?"
```

> [!TIP]
> If you prefer using Gaia this way, you can make the behaviour permenant by
> setting an alias: `alias ai='gaia --related'`.
