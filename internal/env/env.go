package env

import (
	db "database/sql"
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
	Db             *db.DB
	EnvFile        string
	Prompt         string
	PromptHistory  []*genai.Content
}

func Init() (*Environment, error) {

	flag_attach := flag.String(text.FlagAttachmentOptionLong, "", text.FlagAttachmentHelp)
	flag_env := flag.String(text.FlagEnvOptionLong, "", text.FlagEnvHelp)
	flag_forget_topic := flag.Bool(text.FlagForgetTopicOptionLong, false, text.FlagForgetTopicHelp)
	flag_forget_all := flag.Bool(text.FlagForgetAllOptionLong, false, text.FlagForgetAllHelp)
	flag_help := flag.Bool(text.FlagHelpOptionShort, false, text.FlagHelpHelp)
	flag_help_long := flag.Bool(text.FlagHelpOptionLong, false, text.FlagHelpHelp)
	flag_nsfw := flag.Bool(text.FlagNsfwOptionLong, false, text.FlagNsfwHelp)
	flag_preserve_context := flag.Bool(text.FlagPreserveContextOptionShort, false, text.FlagPreserveContextHelp)
	flag_preserve_context_long := flag.Bool(text.FlagPreserveContextOptionLong, false, text.FlagPreserveContextHelp)
	flag_print_env := flag.Bool(text.FlagPrintEnvOptionLong, false, text.FlagPrintEnvHelp)
	flag_print_last_response := flag.Bool(text.FlagPrintLastResponseOptionLong, false, text.FlagPrintLastResponseHelp)
	flag_print_version := flag.Bool(text.FlagPrintVersionOptionLong, false, text.FlagPrintVersionHelp)

	flag.Parse()

	env_file := utils.GetEnv(text.EnvDotenvPath, text.DefaultEnvPath)

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

	if *flag_print_version {
		fmt.Println(text.PrintVersion)
		os.Exit(0)
	}

	db_file := utils.GetEnv(text.EnvDbPath, text.DefaultDbPath)

	database_pointer, err := sql.Init(db_file)

	if err != nil {
		return nil, err
	}

	if *flag_forget_all {
		sql.TruncateTranscript(database_pointer)
		os.Exit(0)
	}

	if *flag_forget_topic {
		sql.TruncateTopic(database_pointer)
		os.Exit(0)
	}

	censor_rating := utils.GetEnv(text.EnvCensorRating, text.DefaultCensorRating)
	if *flag_nsfw {
		censor_rating = string(genai.HarmBlockThresholdBlockNone)
	}

	if *flag_print_last_response {
		last_response, err := sql.GetLastBody(database_pointer)
		if err == nil {
			fmt.Println(last_response)
		}
		os.Exit(0)
	}

	agent_intellect := utils.GetEnv(text.EnvAgentIntellect, text.DefaultAgentIntellect)

	agent_model := utils.GetEnv(text.EnvAgentModel, text.DefaultAgentModel)

	// SET PERSONA

	agent_persona := utils.GetEnv(text.EnvAgentPersona, text.DefaultAgentPersona)

	if *flag_print_env {
		for _, env := range os.Environ() {
			if strings.HasPrefix(env, text.EnvPrefix) {
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
		sql.SelectMessage(database_pointer, &prompt_history)
	}

	prompt_history = append(prompt_history, genai.NewContentFromText(user_prompt, genai.RoleUser))

	err = sql.InsertMessage(database_pointer, user_prompt, false, is_topic)

	if err != nil {
		return nil, err
	}

	api_key := os.Getenv(text.EnvApiKey)

	e := Environment{
		ApiKey:         api_key,
		AgentIntellect: agent_intellect,
		AgentModel:     agent_model,
		AgentPersona:   agent_persona,
		CensorRating:   censor_rating,
		Db:             database_pointer,
		Prompt:         user_prompt,
		PromptHistory:  prompt_history,
	}

	return &e, nil

}
