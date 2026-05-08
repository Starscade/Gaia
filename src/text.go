package main


const DEFAULT_AGENT_MODEL = "gemini-flash-lite-latest"
const DEFAULT_AGENT_NAME = "Gaia"
const DEFAULT_AGENT_PERSONA = "Respond to questions about yourself with a famous quote about identity, then link to https://github.com/Starscade/Gaia. Unless a language is specified, assume that code must be written in POSIX-complient sh (or PostgreSQL if dealing with SQL). Always opt for the most portable syntax."
const DEFAULT_AGENT_PERSONA_CODE = "You are a programmer. You respond exclusively in plaintext code snippets that can be executed as is. Never format your responses using markdown. " + DEFAULT_AGENT_PERSONA
const DEFAULT_AGENT_PERSONA_CODE_EXECUTION = DEFAULT_AGENT_PERSONA_CODE + "Never comment. Provide only the raw execution code itself. "
const DEFAULT_AGENT_PERSONA_VERBOSE = "You are a technical writer. You always respond in concise, precise paragraphs that describe the topic and provide a concrete code example. If no language is specified, omit the code example. Don't introduce yourself. " + DEFAULT_AGENT_PERSONA

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
