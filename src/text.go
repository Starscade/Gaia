package main


const DEFAULT_AGENT_MODEL = "gemini-flash-lite-latest"
const DEFAULT_AGENT_NAME = "Gaia"
const DEFAULT_AGENT_PERSONA = "Unless a language is specified, assume that code must be written in POSIX-complient sh (or PostgreSQL if dealing with SQL). Always opt for the most portable syntax. "
const DEFAULT_AGENT_PERSONA_CODE = DEFAULT_AGENT_PERSONA + "You are a programmer. You respond exclusively in plaintext code snippets that can be executed as is. Never format your responses using markdown. "
const DEFAULT_AGENT_PERSONA_CODE_EXECUTION = DEFAULT_AGENT_PERSONA + DEFAULT_AGENT_PERSONA_CODE + "Never comment. Provide only the raw execution code itself."
const DEFAULT_AGENT_PERSONA_VERBOSE = DEFAULT_AGENT_PERSONA + "You are a technical writer. You always respond in concise, precise paragraphs that describe the topic and provide a concrete code example. Don't introduce yourself. If no language is specified, omit the code example."

const ENV_AGENT_MODEL = "GAIA_AGENT_MODEL"
const ENV_AGENT_NAME = "GAIA_AGENT_NAME"
const ENV_AGENT_PERSONA = "GAIA_AGENT_PERSONA"
const ENV_APIKEY = "GEMINI_API_KEY"

const ERR_NO_APIKEY = "No " + ENV_APIKEY + "!"
const ERR_NO_PROMPT = "No prompt provided!"

const FLAG_EXECUTE_HELP = "Execute the prompt as code."
const FLAG_EXECUTE_OPTION_SHORT = "x"
const FLAG_VERBOSE_HELP = "Speak casually."
const FLAG_VERBOSE_OPTION_SHORT = "v"
