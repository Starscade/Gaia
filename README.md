# Gaia AI Agent

Gaia is a lightweight Go CLI tool that interfaces with the Gemini API to provide code generation, execution, and assistance directly in your terminal.

## Installation

1. Clone the repository.
2. Ensure your `GOPATH` and `PATH` are configured to include `~/.local/bin`.
3. Build the project using the provided Makefile: `make`.

## Configuration

Gaia requires a Gemini API key: `export GEMINI_API_KEY="your-api-key-here"`.

## Usage

### Direct Prompting
`gaia "How do I list files in the current directory?"`

### Execution Mode
Gaia can generate and execute code directly: `gaia -x "echo Hello World"`.
> [!CAUTION]
> Be very careful when writing prompts for this!

### Pipeline Support
Pipe data into Gaia for analysis or transformation: `cat file.txt | gaia -i "Summarize this file"`.

### Flags
- `-x`: Execute the AI response as a shell script.
- `-i`: Accept input via stdin.
- `-v`: Verbose mode (switches persona to technical writer).

## Environment Variables
- `GEMINI_API_KEY`: Required. Your Google Gemini API key.
- `GAIA_AGENT_MODEL`: Override the default model (`gemini-flash-lite-latest`).
- `GAIA_AGENT_NAME`: Customize the AI's persona name.
- `GAIA_AGENT_PERSONA`: Override the default system prompt.

---

> P.S. This README was written for Gaia by Gaia! (Edited by Angus.)
