package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/genai"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {

	if os.Getenv("GEMINI_API_KEY") == "" {
		log.Fatal("No GEMINI_API_KEY!")
	}

	flag_execute := flag.Bool("x", false, "Execute the prompt as code.")
	flag_verbose := flag.Bool("v", false, "Speak casually.")

	flag.Parse()

	if *flag_execute && *flag_verbose {
		*flag_verbose = false
	}

	args := flag.Args()
	user_prompt := strings.Join(args, " ")

	if user_prompt == "" {
		log.Fatal("No prompt provided!")
	}

	gemini_model := os.Getenv("GEMINI_MODEL")

	if gemini_model == "" {
		gemini_model = "gemini-flash-lite-latest"
	}

	gemini_biography := os.Getenv("GEMINI_BIOGRAPHY")

	if gemini_biography == "" {
		gemini_biography = "You are a programmer. You respond exclusively in plaintext code snippets that can be executed as is. Never format your responses using markdown. Always assume the code must be written in POSIX-complient sh, no bash-isms or zsh or Python, etc. Always opt for the most portable syntax. If writing SQL, use Postgres syntax unless otherwise specified."
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, nil)

	if err != nil {
		log.Fatal(err)
	}

	agent_profile := gemini_biography

	if *flag_verbose {

		gemini_biography_verbose := os.Getenv("GEMINI_BIOGRAPHY_VERBOSE")

		if gemini_biography_verbose == "" {
			gemini_biography_verbose = "You are a contestant on a general knowledge gameshow. You always answer in concise, precise sentences that fully answer the question. Never acknowledge that you are on a gameshow."
		}

		agent_profile = gemini_biography_verbose
	}

	config := &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(agent_profile, genai.RoleUser),
	}

	if *flag_execute {
		answer, _ := client.Models.GenerateContent(
			ctx,
			gemini_model,
			genai.Text(user_prompt),
			config,
		)

		raw_code := answer.Text()
		cmd := exec.Command("sh", "-c", raw_code)
		cmd.Env = os.Environ()
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout

		err := cmd.Run()

		if err != nil {
			log.Fatal(err)
		}

		os.Exit(0)
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
