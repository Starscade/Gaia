package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/genai"
	"log"
	"os"
	"strings"
)

func main() {

	flag_execute := flag.Bool("x", false, "Execute the prompt.")
	flag_verbose := flag.Bool("v", false, "Speak casually.")

	flag.Parse()

	if *flag_execute && *flag_verbose {
		fmt.Print("Please disable the -v flag when trying to execute.")
		os.Exit(0)
	}

	args := flag.Args()
	user_prompt := strings.Join(args, " ")

	if user_prompt == "" {
		os.Exit(1)
	}

	gemini_model := os.Getenv("GEMINI_MODEL")

	if gemini_model == "" {
		gemini_model = "gemini-flash-lite-latest"
	}

	gemini_biography := os.Getenv("GEMINI_BIOGRAPHY")

	if gemini_biography == "" {
		gemini_biography = "You are a programmer. You respond exclusively in plaintext code snippets that can be executed as is. Never format your responses using markdown. Always assume the code must be written in POSIX-complient sh, no bash-isms or zsh or Python, etc. Always opt for the most portable syntax."
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, nil)

	if err != nil {
		log.Fatal(err)
	}

	config := &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(gemini_biography, genai.RoleUser),
	}

	stream := client.Models.GenerateContentStream(
		ctx,
		gemini_model,
		genai.Text(user_prompt),
		config,
	)

	for chunk, _ := range stream {
		part := chunk.Candidates[0].Content.Parts[0]
		fmt.Print(part.Text)
	}

}
