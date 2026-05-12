package main

const (
	DEFAULT_AGENT_INTELLECT              = "MINIMAL"
	DEFAULT_AGENT_MODEL                  = "gemini-flash-lite-latest"
	DEFAULT_AGENT_NAME                   = "Gaia"
	DEFAULT_AGENT_PERSONA                = "If no language is specified, write code in POSIX-complient sh (or PostgreSQL if dealing with SQL). Otherwise, write the code in the language that the user mentions. Always use the most portable syntax."
	DEFAULT_AGENT_PERSONA_CODE           = "You are a programmer. You respond exclusively in plaintext code snippets that can be executed (or compiled) as is. Never format your responses using markdown. " + DEFAULT_AGENT_PERSONA
	DEFAULT_AGENT_PERSONA_CODE_EXECUTION = DEFAULT_AGENT_PERSONA_CODE + "Never comment. Provide only the raw execution code itself. "
	DEFAULT_AGENT_PERSONA_VERBOSE        = "You are a technical writer. You always respond in concise, precise paragraphs that describe the topic and provide a concrete code example. If no language is specified, omit the code example. Don't introduce yourself. " + DEFAULT_AGENT_PERSONA
	DEFAULT_DB_PATH                      = "/tmp/gaia.sqlite3"

	ENV_AGENT_INTELLECT = "GAIA_AGENT_INTELLECT"
	ENV_AGENT_MODEL     = "GAIA_AGENT_MODEL"
	ENV_AGENT_NAME      = "GAIA_AGENT_NAME"
	ENV_AGENT_PERSONA   = "GAIA_AGENT_PERSONA"
	ENV_APIKEY          = "GEMINI_API_KEY" // Ensures compatibility with other Gemini API tools.
	ENV_DB_PATH         = "GAIA_DB_PATH"

	ERR_NO_APIKEY = "No " + ENV_APIKEY + "!"
	ERR_NO_PROMPT = "No prompt provided!"

	FLAG_EXECUTE_HELP         = "Execute the prompt as code."
	FLAG_EXECUTE_OPTION_SHORT = "x"
	FLAG_HELP_HELP            = "Show this message."
	FLAG_HELP_OPTION_LONG     = "help"
	FLAG_HELP_OPTION_SHORT    = "h"
	FLAG_NSFW_HELP            = "Free speech mode."
	FLAG_NSFW_OPTION_LONG     = "nsfw"
	FLAG_TOPIC_HELP           = "Preserve topic context from previous query."
	FLAG_TOPIC_OPTION_LONG    = "topic"
	FLAG_TOPIC_OPTION_SHORT   = "t"
	FLAG_VERBOSE_HELP         = "Speak casually."
	FLAG_VERBOSE_OPTION_SHORT = "v"

	PRINT_GAIA       = "\n \033[1mGAIA - Gaia AI Agent\033[0m\n"
	PRINT_HELP_USAGE = "\n \033[1mUsage\033[0m: gaia -x \"What's my IP address?\"\n\n"

	SQL_CREATE_TABLE     = "CREATE TABLE IF NOT EXISTS conversations (id TEXT, ts TEXT, is_agent INT, body TEXT)"
	SQL_INSERT_ROW       = "INSERT INTO conversations (id, ts, is_agent, body) VALUES (?, ?, ?, ?)"
	SQL_GET_CONVERSATION = "SELECT is_agent, body FROM conversations WHERE id = ?"
	SQL_GET_CURRENT_ID   = "SELECT id FROM conversations ORDER BY ts DESC LIMIT 1"
)
