package main

const (
	DEFAULT_AGENT_INTELLECT = "MINIMAL"
	DEFAULT_AGENT_MODEL     = "gemini-flash-lite-latest"
	DEFAULT_AGENT_PERSONA   = "Your name is Gaia. You are a CLI tool. You respond exclusively in plaintext code snippets that can be executed (or compiled) as is. Never format your responses using markdown. If no language is specified, write code in POSIX-compliant sh (or PostgreSQL if dealing with SQL). Otherwise, write the code in the language that the user mentions. Always use the most portable shell syntax (e.g. the oldest, most widely supported). Prefer tab indentation to spaces. Never introduce yourself."
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

	FLAG_ENV_HELP                        = "Path to dotenv file."
	FLAG_ENV_OPTION_LONG                 = "env"
	FLAG_HELP_HELP                       = "Show this message."
	FLAG_HELP_OPTION_LONG                = "help"
	FLAG_HELP_OPTION_SHORT               = "h"
	FLAG_NSFW_HELP                       = "Disable safeguards regardless of ENV settings."
	FLAG_NSFW_OPTION_LONG                = "nsfw"
	FLAG_PRESERVE_CONTEXT_HELP           = "Preserve conversation context."
	FLAG_PRESERVE_CONTEXT_OPTION_LONG    = "related"
	FLAG_PRESERVE_CONTEXT_OPTION_SHORT   = "r"
	FLAG_PRINT_ENV_HELP                  = "Print Gaia's current settings."
	FLAG_PRINT_ENV_OPTION_LONG           = "print-env"
	FLAG_PRINT_LAST_RESPONSE_HELP        = "Print the last response verbatim."
	FLAG_PRINT_LAST_RESPONSE_OPTION_LONG = "echo"

	PRINT_CENSORED         = " \033[1;31;40m CENSORED \033[0m "
	PRINT_COPYRIGHTED      = " \033[1;31;40m COPYRIGHT \033[0m "
	PRINT_ERR              = " \033[1;31;40m ERR \033[0m "
	PRINT_GAIA             = "\n \033[1mGAIA - Gaia AI Agent\033[0m\n\n"
	PRINT_HELP_USAGE       = "\n \033[1mUsage\033[0m: `gaia \"What's my IP address?\" | sh`\n\n\n"
	PRINT_TOKENS_EXHAUSTED = " \033[1;31;40m TOKENS EXHAUSTED \033[0m "

	SQL_CREATE_TABLE      = "CREATE TABLE IF NOT EXISTS transcript (id TEXT, ts TEXT, is_agent INT, body TEXT)"
	SQL_INSERT_ROW        = "INSERT INTO transcript (id, ts, is_agent, body) VALUES (?, ?, ?, ?)"
	SQL_GET_CONVERSATION  = "SELECT is_agent, body FROM transcript WHERE id = ?"
	SQL_GET_CURRENT_ID    = "SELECT id FROM transcript ORDER BY ts DESC LIMIT 1"
	SQL_GET_LAST_RESPONSE = "SELECT body FROM transcript WHERE is_agent = true ORDER BY ts DESC LIMIT 1"
)
