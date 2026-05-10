package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

func getHistory() {
	jsonStr := os.Getenv(ENV_CHAT_HISTORY)
	if jsonStr == "" {
		jsonStr = "[]"
	}

	var data []map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		printErr(err.Error())
	}

	newData := map[string]interface{}{"id": 1, "status": "active"}
	data = append(data, newData)

	updatedJson, err := json.Marshal(data)
	if err != nil {
		printErr(err.Error())
	}

	os.Setenv(ENV_CHAT_HISTORY, string(updatedJson))
	fmt.Println(os.Getenv(ENV_CHAT_HISTORY))
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
