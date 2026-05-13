package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"google.golang.org/genai"
)

var api_key string
var agent_intellect string
var agent_model string
var agent_name string
var agent_persona string
var censor_rating string
var db_file string
var env_file string
var flag_env *string
var flag_execute *bool
var flag_help *bool
var flag_help_long *bool
var flag_nsfw *bool
var flag_recall *bool
var flag_topic *bool
var flag_topic_long *bool
var flag_verbose *bool
var user_prompt string
var prompt_history []*genai.Content

func initEnv() {

	flag_env = flag.String(FLAG_ENV_OPTION_LONG, "", FLAG_ENV_HELP)
	flag_execute = flag.Bool(FLAG_EXECUTE_OPTION_SHORT, false, FLAG_EXECUTE_HELP)
	flag_help = flag.Bool(FLAG_HELP_OPTION_SHORT, false, FLAG_HELP_HELP)
	flag_help_long = flag.Bool(FLAG_HELP_OPTION_LONG, false, FLAG_HELP_HELP)
	flag_nsfw = flag.Bool(FLAG_NSFW_OPTION_LONG, false, FLAG_NSFW_HELP)
	flag_recall = flag.Bool(FLAG_RECALL_OPTION_LONG, false, FLAG_RECALL_HELP)
	flag_topic = flag.Bool(FLAG_TOPIC_OPTION_SHORT, false, FLAG_TOPIC_HELP)
	flag_topic_long = flag.Bool(FLAG_TOPIC_OPTION_LONG, false, FLAG_TOPIC_HELP)
	flag_verbose = flag.Bool(FLAG_VERBOSE_OPTION_SHORT, false, FLAG_VERBOSE_HELP)

	flag.Parse()

	env_file = getEnv(ENV_DOTENV_PATH, DEFAULT_ENV_PATH)

	_, err := os.Stat(env_file)
	if err == nil {
		err := godotenv.Overload(env_file)
		exitOnErr(err)
	}

	if *flag_env != "" {
		err := godotenv.Overload(*flag_env)
		exitOnErr(err)
	}

	api_key = os.Getenv(ENV_APIKEY)

	if api_key == "" {
		log.Fatal(ERR_NO_APIKEY) // No key? Why continue?
	}

	if *flag_help_long || *flag_help {
		printHelp()
		os.Exit(0)
	}

	censor_rating = getEnv(ENV_CENSOR_RATING, DEFAULT_CENSOR_RATING)

	db_file = getEnv(ENV_DB_PATH, DEFAULT_DB_PATH)

	initDb()

	if *flag_recall {
		fmt.Println(getLastBody())
		os.Exit(0)
	}

	if *flag_execute {
		// Don't build a rootkit or execute poetry.
		*flag_nsfw = false
		*flag_verbose = false
	}

	args := flag.Args()
	user_prompt = strings.Join(args, " ") + getStdin()

	if user_prompt == "" {
		printHelp()
		printErr(ERR_NO_PROMPT)
	}

	agent_intellect = getEnv(ENV_AGENT_INTELLECT, DEFAULT_AGENT_INTELLECT)

	is_topic := false
	if *flag_topic || *flag_topic_long {
		is_topic = true
		getHistory(&prompt_history)
	}

	prompt_history = append(prompt_history, genai.NewContentFromText(user_prompt, genai.RoleUser))

	setHistory(user_prompt, false, is_topic)

	agent_model = getEnv(ENV_AGENT_MODEL, DEFAULT_AGENT_MODEL)

	// SET PERSONA

	agent_name = getEnv(ENV_AGENT_NAME, DEFAULT_AGENT_NAME)

	agent_persona = os.Getenv(ENV_AGENT_PERSONA)

	if agent_persona == "" {

		agent_persona = DEFAULT_AGENT_PERSONA_CODE

		if *flag_execute {
			agent_persona = DEFAULT_AGENT_PERSONA_CODE_EXECUTION
		}

		if *flag_verbose {
			agent_persona = DEFAULT_AGENT_PERSONA_VERBOSE
		}

		agent_persona = "Your name is " + agent_name + ". " + agent_persona

	}

}
