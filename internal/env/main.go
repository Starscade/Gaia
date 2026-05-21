package env

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"google.golang.org/genai"

	"github.com/Starscade/Gaia/internal/sql"
	"github.com/Starscade/Gaia/internal/text"
	"github.com/Starscade/Gaia/internal/utils"
)

type Environment struct {
	ApiKey         string
	AgentIntellect string
	AgentModel     string
	AgentPersona   string
	CensorRating   string
	Client         *genai.Client
	Config         *genai.GenerateContentConfig
	DbFile         string
	EnvFile        string
	Prompt         string
	PromptHistory  []*genai.Content
}

func Init() (*Environment, error) {

	flag_attach := flag.String(text.FLAG_ATTACHMENT_OPTION_LONG, "", text.FLAG_ATTACHMENT_HELP)
	flag_env := flag.String(text.FLAG_ENV_OPTION_LONG, "", text.FLAG_ENV_HELP)
	flag_forget := flag.Bool(text.FLAG_FORGET_OPTION_LONG, false, text.FLAG_FORGET_HELP)
	flag_help := flag.Bool(text.FLAG_HELP_OPTION_SHORT, false, text.FLAG_HELP_HELP)
	flag_help_long := flag.Bool(text.FLAG_HELP_OPTION_LONG, false, text.FLAG_HELP_HELP)
	flag_nsfw := flag.Bool(text.FLAG_NSFW_OPTION_LONG, false, text.FLAG_NSFW_HELP)
	flag_print_last_response := flag.Bool(text.FLAG_PRINT_LAST_RESPONSE_OPTION_LONG, false, text.FLAG_PRINT_LAST_RESPONSE_HELP)
	flag_preserve_context := flag.Bool(text.FLAG_PRESERVE_CONTEXT_OPTION_SHORT, false, text.FLAG_PRESERVE_CONTEXT_HELP)
	flag_preserve_context_long := flag.Bool(text.FLAG_PRESERVE_CONTEXT_OPTION_LONG, false, text.FLAG_PRESERVE_CONTEXT_HELP)
	flag_print_env := flag.Bool(text.FLAG_PRINT_ENV_OPTION_LONG, false, text.FLAG_PRINT_ENV_HELP)

	flag.Parse()

	env_file := utils.GetEnv(text.ENV_DOTENV_PATH, text.DEFAULT_ENV_PATH)

	_, err := os.Stat(env_file)
	if err == nil {
		err := godotenv.Load(env_file)
		if err != nil {
			return nil, err
		}
	}

	if *flag_env != "" {
		err := godotenv.Overload(*flag_env)
		if err != nil {
			return nil, err
		}
	}

	if *flag_help_long || *flag_help {
		utils.PrintHelp()
		os.Exit(0)
	}

	db_file := utils.GetEnv(text.ENV_DB_PATH, text.DEFAULT_DB_PATH)

	sql.Init(db_file)

	if *flag_forget {
		sql.TruncateTranscript(db_file)
		os.Exit(0)
	}

	censor_rating := utils.GetEnv(text.ENV_CENSOR_RATING, text.DEFAULT_CENSOR_RATING)
	if *flag_nsfw {
		censor_rating = string(genai.HarmBlockThresholdBlockNone)
	}

	if *flag_print_last_response {
		last_response, err := sql.GetLastBody(db_file)
		if err == nil {
			fmt.Println(last_response)
		}
		os.Exit(0)
	}

	agent_intellect := utils.GetEnv(text.ENV_AGENT_INTELLECT, text.DEFAULT_AGENT_INTELLECT)

	agent_model := utils.GetEnv(text.ENV_AGENT_MODEL, text.DEFAULT_AGENT_MODEL)

	// SET PERSONA

	agent_persona := utils.GetEnv(text.ENV_AGENT_PERSONA, text.DEFAULT_AGENT_PERSONA)

	if *flag_print_env {
		for _, env := range os.Environ() {
			if strings.HasPrefix(env, text.ENV_PREFIX) {
				fmt.Println(env)
			}
		}
		os.Exit(0)
	}

	args := flag.Args()
	user_prompt := strings.Join(args, " ") + utils.GetStdin()

	if *flag_attach != "" {
		attached_text, err := utils.GetFile(*flag_attach)
		if err != nil {
			return nil, err
		}
		user_prompt = user_prompt + attached_text
	}

	var prompt_history []*genai.Content

	is_topic := false
	if *flag_preserve_context || *flag_preserve_context_long {
		is_topic = true
		sql.SelectMessage(db_file, &prompt_history)
	}

	prompt_history = append(prompt_history, genai.NewContentFromText(user_prompt, genai.RoleUser))

	sql.InsertMessage(db_file, user_prompt, false, is_topic)

	api_key := os.Getenv(text.ENV_API_KEY)

	e := Environment{
		ApiKey:         api_key,
		AgentIntellect: agent_intellect,
		AgentModel:     agent_model,
		AgentPersona:   agent_persona,
		CensorRating:   censor_rating,
		DbFile:         db_file,
		Prompt:         user_prompt,
		PromptHistory:  prompt_history,
	}

	return &e, nil

}
