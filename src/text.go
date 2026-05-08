package main


const DEFAULT_AGENT_MODEL = "gemini-flash-lite-latest"
const DEFAULT_AGENT_NAME = "Gaia"
const DEFAULT_AGENT_PERSONA = "You are a programmer. You respond exclusively in plaintext code snippets that can be executed as is. Never format your responses using markdown. Always assume the code must be written in POSIX-complient sh, no bash-isms or zsh or Python, etc. Always opt for the most portable syntax. If writing SQL, use Postgres syntax unless otherwise specified."
const DEFAULT_AGENT_PERSONA_VERBOSE = "You are a technical writer. You always respond in concise, precise paragraphs that could be drop-in replacements for code documentation."

const ENV_AGENT_MODEL = "GAIA_AGENT_MODEL"
const ENV_AGENT_PERSONA = "GAIA_AGENT_PERSONA"
const ENV_AGENT_PERSONA_VERBOSE = "GAIA_AGENT_PERSONA_VERBOSE"
const ENV_APIKEY = "GEMINI_API_KEY"

const ERR_NO_APIKEY = "No GEMINI_API_KEY!"
const ERR_NO_PROMPT = "No prompt provided!"

const FLAG_EXECUTE_OPTION_SHORT = "x"
const FLAG_EXECUTE_HELP = "Execute the prompt as code."
const FLAG_VERBOSE_OPTION_SHORT = "v"
const FLAG_VERBOSE_HELP = "Speak casually."
