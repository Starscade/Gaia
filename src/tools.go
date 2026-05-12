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
	log.Fatal(" \033[1;31mERR:\033[0m ", err)
}

func printHelp() {
	fmt.Println(PRINT_GAIA)
	flag.PrintDefaults()
	fmt.Println(PRINT_HELP_USAGE)
}

func setHistory(body string, is_agent bool, is_topic bool) {
	insertMessage(body, is_agent, is_topic)
}
