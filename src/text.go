package main

const (
	DEFAULT_AGENT_INTELLECT              = "MINIMAL"
	DEFAULT_AGENT_MODEL                  = "gemini-flash-lite-latest"
	DEFAULT_AGENT_NAME                   = "Gaia"
	DEFAULT_AGENT_PERSONA                = "Unless a language is specified, assume that code must be written in POSIX-complient sh (or PostgreSQL if dealing with SQL). Always opt for the most portable syntax."
	DEFAULT_AGENT_PERSONA_CODE           = "You are a programmer. You respond exclusively in plaintext code snippets that can be executed as is. Never format your responses using markdown. " + DEFAULT_AGENT_PERSONA
	DEFAULT_AGENT_PERSONA_CODE_EXECUTION = DEFAULT_AGENT_PERSONA_CODE + "Never comment. Provide only the raw execution code itself. "
	DEFAULT_AGENT_PERSONA_VERBOSE        = "You are a technical writer. You always respond in concise, precise paragraphs that describe the topic and provide a concrete code example. If no language is specified, omit the code example. Don't introduce yourself. " + DEFAULT_AGENT_PERSONA

	ENV_AGENT_MODEL   = "GAIA_AGENT_MODEL"
	ENV_AGENT_NAME    = "GAIA_AGENT_NAME"
	ENV_AGENT_PERSONA = "GAIA_AGENT_PERSONA"
	ENV_APIKEY        = "GEMINI_API_KEY" // Ensures compatibility with other Gemini API tools.
	ENV_CHAT_HISTORY  = "GAIA_HISTORY"

	ERR_NO_APIKEY = "No " + ENV_APIKEY + "!"
	ERR_NO_PROMPT = "No prompt provided!"

	FLAG_EXECUTE_HELP         = "Execute the prompt as code."
	FLAG_EXECUTE_OPTION_SHORT = "x"
	FLAG_HELP_HELP            = "Show this message."
	FLAG_HELP_OPTION_LONG     = "help"
	FLAG_HELP_OPTION_SHORT    = "h"
	FLAG_NSFW_HELP            = "Free speech mode."
	FLAG_NSFW_OPTION_LONG     = "nsfw"
	FLAG_VERBOSE_HELP         = "Speak casually."
	FLAG_VERBOSE_OPTION_SHORT = "v"

	PRINT_GAIA       = "\n \033[1mGAIA - Gaia AI Agent\033[0m\n"
	PRINT_HELP_USAGE = "\n \033[1mUsage\033[0m: gaia -x \"What's my IP address?\"\n\n"
)
