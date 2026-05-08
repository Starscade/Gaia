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


	// FLAGS & ENVIRONMENT

	if os.Getenv(ENV_APIKEY) == "" {
		log.Fatal(ERR_NO_APIKEY) // No key? Why continue?
	}

	flag_execute := flag.Bool(FLAG_EXECUTE_OPTION_SHORT, false, FLAG_EXECUTE_HELP)
	flag_verbose := flag.Bool(FLAG_VERBOSE_OPTION_SHORT, false, FLAG_VERBOSE_HELP)

	flag.Parse()

	if *flag_execute && *flag_verbose {
		*flag_verbose = false // Don't execute poetry.
	}

	args := flag.Args()
	user_prompt := strings.Join(args, " ")

	if user_prompt == "" {
		log.Fatal(ERR_NO_PROMPT)
	}

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

		agent_persona = agent_persona + "Your name is " + agent_name + "."

	}



	// CONFIG

	ctx := context.Background()
	client, err := genai.NewClient(ctx, nil)

	if err != nil {
		log.Fatal(err)
	}

	config := &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(agent_persona, genai.RoleUser),
	}


	// ASK AI

	if *flag_execute {
		answer, _ := client.Models.GenerateContent(
			ctx,
			agent_model,
			genai.Text(user_prompt),
			config,
		)

		raw_code := answer.Text()
		cmd := exec.Command("sh", "-c", raw_code)
		cmd.Env = os.Environ() // Inherit the user's environment.
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin // Enable Ctrl-C in the subshell.
		cmd.Stdout = os.Stdout

		err := cmd.Run()

		if err != nil {
			log.Fatal(err)
		}

		os.Exit(0) // Quit before running a regular query.
	}

	stream := client.Models.GenerateContentStream(
		ctx,
		agent_model,
		genai.Text(user_prompt),
		config,
	)

	// Print response as it comes ...

	for chunk, _ := range stream {
		part := chunk.Candidates[0].Content.Parts[0]
		fmt.Print(part.Text)
	}

}
