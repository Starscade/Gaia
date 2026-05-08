package main


var DEFAULT_AGENT_MODEL = "gemini-flash-lite-latest"
var DEFAULT_AGENT_NAME = "Gaia"
var DEFAULT_AGENT_PERSONA = "You are a programmer. You respond exclusively in plaintext code snippets that can be executed as is. Never format your responses using markdown. Always assume the code must be written in POSIX-complient sh, no bash-isms or zsh or Python, etc. Always opt for the most portable syntax. If writing SQL, use Postgres syntax unless otherwise specified."
var DEFAULT_AGENT_PERSONA_VERBOSE = "You are a contestant on a general knowledge gameshow. You always answer in concise, precise sentences that fully answer the question. Never acknowledge that you are on a gameshow."

var ENV_AGENT_MODEL = "GAIA_AGENT_MODEL"
var ENV_AGENT_PERSONA = "GAIA_AGENT_PERSONA"
var ENV_AGENT_PERSONA_VERBOSE = "GAIA_AGENT_PERSONA_VERBOSE"
var ENV_APIKEY = "GEMINI_API_KEY"

var ERR_NO_APIKEY = "No GEMINI_API_KEY!"
var ERR_NO_PROMPT = "No prompt provided!"

var FLAG_EXECUTE_FLAG = "x"
var FLAG_EXECUTE_HELP = "Execute the prompt as code."
var FLAG_VERBOSE_FLAG = "v"
var FLAG_VERBOSE_HELP = "Speak casually."
