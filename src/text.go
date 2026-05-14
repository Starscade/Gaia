package main

const (
	DEFAULT_AGENT_INTELLECT = "MINIMAL"
	DEFAULT_AGENT_MODEL     = "gemini-flash-lite-latest"
	DEFAULT_AGENT_PERSONA   = "Your name is Gaia. You are a programmer. You respond exclusively in plaintext code snippets that can be executed (or compiled) as is. Never format your responses using markdown. If no language is specified, write code in POSIX-complient sh (or PostgreSQL if dealing with SQL). Otherwise, write the code in the language that the user mentions. Always use the most portable shell syntax (e.g. the oldest, most widely supported). Prefer tab indentation to spaces. Never introduce yourself."
	DEFAULT_CENSOR_RATING   = "BLOCK_MEDIUM_AND_ABOVE"
	DEFAULT_DB_PATH         = "/tmp/gaia.sqlite3"
	DEFAULT_ENV_PATH        = ".env"

	ENV_AGENT_INTELLECT = "GAIA_AGENT_INTELLECT"
	ENV_AGENT_MODEL     = "GAIA_AGENT_MODEL"
	ENV_AGENT_PERSONA   = "GAIA_AGENT_PERSONA"
	ENV_API_KEY         = "GAIA_AGENT_KEY"
	ENV_CENSOR_RATING   = "GAIA_AGENT_CENSOR_RATING"
	ENV_DB_PATH         = "GAIA_DB_PATH"
	ENV_DOTENV_PATH     = "GAIA_ENV_FILE"

	ERR_NO_API_KEY = "No " + ENV_API_KEY + "!"
	ERR_NO_PROMPT  = "No prompt provided!"

	FLAG_ENV_HELP                        = "Override the environment with a file."
	FLAG_ENV_OPTION_LONG                 = "env"
	FLAG_HELP_HELP                       = "Show this message."
	FLAG_HELP_OPTION_LONG                = "help"
	FLAG_HELP_OPTION_SHORT               = "h"
	FLAG_NSFW_HELP                       = "Disable safeguards."
	FLAG_NSFW_OPTION_LONG                = "nsfw"
	FLAG_PRESERVE_CONTEXT_HELP           = "Preserve conversation context."
	FLAG_PRESERVE_CONTEXT_OPTION_LONG    = "feed"
	FLAG_PRESERVE_CONTEXT_OPTION_SHORT   = "f"
	FLAG_PRINT_ENV_HELP                  = "Print Gaia's current settings."
	FLAG_PRINT_ENV_OPTION_LONG           = "print-env"
	FLAG_PRINT_LAST_RESPONSE_HELP        = "Print the last response verbatim."
	FLAG_PRINT_LAST_RESPONSE_OPTION_LONG = "replay"

	PRINT_GAIA       = "\n \033[1mGAIA - Gaia AI Agent\033[0m\n"
	PRINT_HELP_USAGE = "\n \033[1mUsage\033[0m: `gaia \"What's my IP address?\" | sh`\n\n"

	SQL_CREATE_TABLE      = "CREATE TABLE IF NOT EXISTS conversations (id TEXT, ts TEXT, is_agent INT, body TEXT)"
	SQL_INSERT_ROW        = "INSERT INTO conversations (id, ts, is_agent, body) VALUES (?, ?, ?, ?)"
	SQL_GET_CONVERSATION  = "SELECT is_agent, body FROM conversations WHERE id = ?"
	SQL_GET_CURRENT_ID    = "SELECT id FROM conversations ORDER BY ts DESC LIMIT 1"
	SQL_GET_LAST_RESPONSE = "SELECT body FROM conversations WHERE is_agent = true ORDER BY ts DESC LIMIT 1"
)
