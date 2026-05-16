package main

import (
	"context"
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
var agent_persona string
var censor_rating string
var client *genai.Client
var config *genai.GenerateContentConfig
var ctx context.Context
var db_file string
var env_file string
var flag_env *string
var flag_help *bool
var flag_help_long *bool
var flag_nsfw *bool
var flag_preserve_context *bool
var flag_preserve_context_long *bool
var flag_print_env *bool
var flag_print_last_response *bool
var prompt_history []*genai.Content
var user_prompt string

func initEnv() {

	flag_env = flag.String(FLAG_ENV_OPTION_LONG, "", FLAG_ENV_HELP)
	flag_help = flag.Bool(FLAG_HELP_OPTION_SHORT, false, FLAG_HELP_HELP)
	flag_help_long = flag.Bool(FLAG_HELP_OPTION_LONG, false, FLAG_HELP_HELP)
	flag_nsfw = flag.Bool(FLAG_NSFW_OPTION_LONG, false, FLAG_NSFW_HELP)
	flag_print_last_response = flag.Bool(FLAG_PRINT_LAST_RESPONSE_OPTION_LONG, false, FLAG_PRINT_LAST_RESPONSE_HELP)
	flag_preserve_context = flag.Bool(FLAG_PRESERVE_CONTEXT_OPTION_SHORT, false, FLAG_PRESERVE_CONTEXT_HELP)
	flag_preserve_context_long = flag.Bool(FLAG_PRESERVE_CONTEXT_OPTION_LONG, false, FLAG_PRESERVE_CONTEXT_HELP)
	flag_print_env = flag.Bool(FLAG_PRINT_ENV_OPTION_LONG, false, FLAG_PRINT_ENV_HELP)

	flag.Parse()

	env_file = getEnv(ENV_DOTENV_PATH, DEFAULT_ENV_PATH)

	_, err := os.Stat(env_file)
	if err == nil {
		err := godotenv.Load(env_file)
		exitOnErr(err)
	}

	if *flag_env != "" {
		err := godotenv.Overload(*flag_env)
		exitOnErr(err)
	}

	api_key = os.Getenv(ENV_API_KEY)

	if api_key == "" {
		log.Fatal(ERR_NO_API_KEY) // No key? Why continue?
	}

	if *flag_help_long || *flag_help {
		printHelp()
		os.Exit(0)
	}

	censor_rating = getEnv(ENV_CENSOR_RATING, DEFAULT_CENSOR_RATING)
	if *flag_nsfw {
		censor_rating = string(genai.HarmBlockThresholdBlockNone)
	}

	db_file = getEnv(ENV_DB_PATH, DEFAULT_DB_PATH)

	initDb()

	if *flag_print_last_response {
		last_response, err := getLastBody()
		if err == nil {
			fmt.Println(last_response)
		}
		os.Exit(0)
	}

	agent_intellect = getEnv(ENV_AGENT_INTELLECT, DEFAULT_AGENT_INTELLECT)

	agent_model = getEnv(ENV_AGENT_MODEL, DEFAULT_AGENT_MODEL)

	// SET PERSONA

	agent_persona = getEnv(ENV_AGENT_PERSONA, DEFAULT_AGENT_PERSONA)

	if *flag_print_env {
		for _, env := range os.Environ() {
			if strings.HasPrefix(env, "GAIA_") {
				fmt.Println(env)
			}
		}
		os.Exit(0)
	}

	args := flag.Args()
	user_prompt = strings.Join(args, " ") + getStdin()

	if user_prompt == "" {
		os.Exit(1)
	}

	is_topic := false
	if *flag_preserve_context || *flag_preserve_context_long {
		is_topic = true
		getHistory(&prompt_history)
	}

	prompt_history = append(prompt_history, genai.NewContentFromText(user_prompt, genai.RoleUser))

	setHistory(user_prompt, false, is_topic)

}
