package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

func exitOnErr(err error) {
	if err != nil {
		printErr(err.Error())
	}
}

func getHistory() string {
	return os.Getenv(ENV_CHAT_HISTORY)
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

func setHistory(body string) {
	history := getHistory()
	revised_history := history + "\n\n\n" + body
	os.Setenv(ENV_CHAT_HISTORY, revised_history)
	insertMessage(body, true)
}
