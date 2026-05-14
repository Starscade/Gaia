package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"google.golang.org/genai"
)

func exitOnErr(err error) {
	if err != nil {
		printErr(err.Error())
	}
}

func getEnv(env_var_name string, default_value string) string {
	env_var := os.Getenv(env_var_name)

	if env_var == "" {
		env_var = default_value
	}

	os.Setenv(env_var_name, env_var)

	return env_var
}

func getHistory(chat_history *[]*genai.Content) {
	selectMessage(chat_history)
}

func getStdin() string {
	stat, _ := os.Stdin.Stat()
	var input []byte
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		input, _ = io.ReadAll(os.Stdin)
	}
	return string(input)
}

func printErr(err string) {
	log.Fatal(PRINT_ERR, err)
}

func printHelp() {
	fmt.Println(PRINT_GAIA)
	flag.PrintDefaults()
	fmt.Println(PRINT_HELP_USAGE)
}

func setHistory(body string, is_agent bool, is_topic bool) {
	insertMessage(body, is_agent, is_topic)
}
