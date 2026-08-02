# Gaia AI Agent

Gaia is a standalone ESM module for Google Gemini.
The included PWA is for reference only.

## Installation

1. `git clone https://github.com/Starscade/Gaia.git`
2. `cd ./Gaia`
3. `deno task serve`

## Usage

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
