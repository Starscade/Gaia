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

```
export GAIA_API_KEY="your-gemini-key"

# Or...

gaia --env .env <arguments>
```

> [!TIP]
> You can create a template of your current settings with:
> `gaia --environment > .env`. (This includes internal defaults.)


## Usage

###### CHAT

```
gaia "Ahoy!"
```

> [!IMPORTANT]
> Gaia does _not_ preserve context between prompts and responses. To continue a
> conversation, use the `--related` flag:

```
gaia "I like root beer."
gaia --related "Do I like root beer?"
```

> [!TIP]
> If you prefer using Gaia this way, you can make the behaviour permenant by
> setting an alias: `alias ai='gaia --related'`.


###### DIRECT EXECUTION

`gaia "What's my IP address?" | sh`

> [!CAUTION]
> A safer way to execute responses directly is to verify their output before
> piping (e.g. `gaia "What's my IP address?"`), then use `gaia --echo | sh` to
> execute the most recent response verbatim.


###### FILE ANALYSIS

`gaia --read README.md "Summarize this..."`


###### UNIX PIPES

```
cat *.js | gaia "Consolidate these into a single TypeScript module." > mod.ts
```


###### PICTURES

```
gaia --draw "cyberpunk cat" | base64 -d > cyber-cat.jpg
```


###### IN THE BROWSER

```
import Gaia from 'https://esm.sh/gh/Starscade/Gaia/gaia.js'

const ai = new Gaia({
	api_key: '<your_api_key>',
	print_function: (response) => {
		console.log(response.data)
	}
})

const result = await ai.ask({
	user_prompt: 'Ahoy!'
})

if (result.err) console.error(result.err)
```
