package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

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

func getFile(glob_path string) string {
	files, err := filepath.Glob(glob_path)
	exitOnErr(err)
	var bodies strings.Builder
	for _, file_path := range files {
		file, err := os.ReadFile(file_path)
		exitOnErr(err)
		eof_start := "EOF " + file_path + "\n\n"
		eof_end := "\n\nEOF " + file_path
		bodies.WriteString(eof_start + string(file) + eof_end)
	}
	return bodies.String()
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
	fmt.Print(PRINT_GAIA)
	flag.PrintDefaults()
	fmt.Print(PRINT_HELP_USAGE)
}

func setHistory(body string, is_agent bool, is_topic bool) {
	insertMessage(body, is_agent, is_topic)
}
