package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/joho/godotenv"
	"google.golang.org/genai"
)

var DB_FILE = os.Getenv(ENV_DB_PATH)

func main() {

	// FLAGS & ENVIRONMENT

	flag_env := flag.String(FLAG_ENV_OPTION_LONG, "", FLAG_ENV_HELP)
	flag_execute := flag.Bool(FLAG_EXECUTE_OPTION_SHORT, false, FLAG_EXECUTE_HELP)
	flag_help := flag.Bool(FLAG_HELP_OPTION_SHORT, false, FLAG_HELP_HELP)
	flag_help_long := flag.Bool(FLAG_HELP_OPTION_LONG, false, FLAG_HELP_HELP)
	flag_nsfw := flag.Bool(FLAG_NSFW_OPTION_LONG, false, FLAG_NSFW_HELP)
	flag_recall := flag.Bool(FLAG_RECALL_OPTION_LONG, false, FLAG_RECALL_HELP)
	flag_topic := flag.Bool(FLAG_TOPIC_OPTION_SHORT, false, FLAG_TOPIC_HELP)
	flag_topic_long := flag.Bool(FLAG_TOPIC_OPTION_LONG, false, FLAG_TOPIC_HELP)
	flag_verbose := flag.Bool(FLAG_VERBOSE_OPTION_SHORT, false, FLAG_VERBOSE_HELP)

	flag.Parse()

	env_file := os.Getenv(ENV_DOTENV_PATH)
	if env_file == "" {
		env_file = DEFAULT_ENV_FILE
	}

	_, err := os.Stat(env_file)
	if err == nil {
		err := godotenv.Overload(env_file)
		exitOnErr(err)
	}

	if *flag_env != "" {
		err := godotenv.Overload(*flag_env)
		exitOnErr(err)
	}

	api_key := os.Getenv(ENV_APIKEY)

	if api_key == "" {
		log.Fatal(ERR_NO_APIKEY) // No key? Why continue?
	}

	if *flag_help_long || *flag_help {
		printHelp()
		os.Exit(0)
	}

	if DB_FILE == "" {
		DB_FILE = DEFAULT_DB_PATH
	}

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
	user_prompt := strings.Join(args, " ") + getStdin()

	if user_prompt == "" {
		printHelp()
		printErr(ERR_NO_PROMPT)
	}

	agent_intellect := os.Getenv(ENV_AGENT_INTELLECT)
	if agent_intellect == "" {
		agent_intellect = DEFAULT_AGENT_INTELLECT
	}

	var prompt_history []*genai.Content

	is_topic := false
	if *flag_topic || *flag_topic_long {
		is_topic = true
		getHistory(&prompt_history)
	}

	prompt_history = append(prompt_history, genai.NewContentFromText(user_prompt, genai.RoleUser))

	setHistory(user_prompt, false, is_topic)

	agent_model := os.Getenv(ENV_AGENT_MODEL)

	if agent_model == "" {
		agent_model = DEFAULT_AGENT_MODEL
	}

	// SET PERSONA

	agent_name := os.Getenv(ENV_AGENT_NAME)

	if agent_name == "" {
		agent_name = DEFAULT_AGENT_NAME
	}

	agent_persona := os.Getenv(ENV_AGENT_PERSONA)

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

	// CONFIG

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		Backend: genai.BackendGeminiAPI,
		APIKey:  api_key,
	})

	exitOnErr(err)

	config := &genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{
				&genai.Part{
					Text: agent_persona,
				},
			},
		},
		ThinkingConfig: &genai.ThinkingConfig{
			ThinkingLevel: genai.ThinkingLevel(agent_intellect),
		},
		Tools: []*genai.Tool{
			{
				GoogleSearch: &genai.GoogleSearch{},
			},
		},
	}

	if *flag_nsfw {
		config.SafetySettings = []*genai.SafetySetting{
			{
				Category:  genai.HarmCategoryDangerousContent,
				Threshold: genai.HarmBlockThresholdBlockNone,
			},
			{
				Category:  genai.HarmCategoryHarassment,
				Threshold: genai.HarmBlockThresholdBlockNone,
			},
			{
				Category:  genai.HarmCategoryHateSpeech,
				Threshold: genai.HarmBlockThresholdBlockNone,
			},
			{
				Category:  genai.HarmCategorySexuallyExplicit,
				Threshold: genai.HarmBlockThresholdBlockNone,
			},
		}
	}

	// ASK AI

	if *flag_execute {
		answer, _ := client.Models.GenerateContent(
			ctx,
			agent_model,
			prompt_history,
			config,
		)

		raw_code := answer.Text()

		setHistory(raw_code, true, true)

		cmd := exec.Command("sh", "-c", raw_code)
		cmd.Env = os.Environ() // Inherit the user's environment.
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin // Enable Ctrl-C in the subshell.
		cmd.Stdout = os.Stdout

		err := cmd.Run()

		exitOnErr(err)

		os.Exit(0) // Quit before running a regular query.
	}

	stream := client.Models.GenerateContentStream(
		ctx,
		agent_model,
		prompt_history,
		config,
	)

	// Print response as it comes ...

	complete_response := ""

	for chunk, err := range stream {

		exitOnErr(err)

		if len(chunk.Candidates) > 0 {
			part := chunk.Candidates[0].Content.Parts[0]
			complete_response = complete_response + part.Text
			fmt.Print(part.Text)
		}
	}

	setHistory(complete_response, true, true)

	fmt.Println() // Ensures terminal starts on a new line.

}
