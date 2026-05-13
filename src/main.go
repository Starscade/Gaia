package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"google.golang.org/genai"
)

func main() {

	initEnv()

	// CONFIG

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		Backend: genai.BackendGeminiAPI,
		APIKey:  api_key,
	})

	exitOnErr(err)

	config := &genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{
				&genai.Part{
					Text: agent_persona,
				},
			},
		},
		ThinkingConfig: &genai.ThinkingConfig{
			ThinkingLevel: genai.ThinkingLevel(agent_intellect),
		},
		Tools: []*genai.Tool{
			{
				GoogleSearch: &genai.GoogleSearch{},
			},
		},
	}

	if *flag_nsfw {
		config.SafetySettings = []*genai.SafetySetting{
			{
				Category:  genai.HarmCategoryDangerousContent,
				Threshold: genai.HarmBlockThresholdBlockNone,
			},
			{
				Category:  genai.HarmCategoryHarassment,
				Threshold: genai.HarmBlockThresholdBlockNone,
			},
			{
				Category:  genai.HarmCategoryHateSpeech,
				Threshold: genai.HarmBlockThresholdBlockNone,
			},
			{
				Category:  genai.HarmCategorySexuallyExplicit,
				Threshold: genai.HarmBlockThresholdBlockNone,
			},
		}
	}

	// ASK AI

	if *flag_execute {
		answer, _ := client.Models.GenerateContent(
			ctx,
			agent_model,
			prompt_history,
			config,
		)

		raw_code := answer.Text()

		setHistory(raw_code, true, true)

		cmd := exec.Command("sh", "-c", raw_code)
		cmd.Env = os.Environ() // Inherit the user's environment.
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin // Enable Ctrl-C in the subshell.
		cmd.Stdout = os.Stdout

		err := cmd.Run()

		exitOnErr(err)

		os.Exit(0) // Quit before running a regular query.
	}

	stream := client.Models.GenerateContentStream(
		ctx,
		agent_model,
		prompt_history,
		config,
	)

	// Print response as it comes ...

	complete_response := ""

	for chunk, err := range stream {

		exitOnErr(err)

		if len(chunk.Candidates) > 0 {
			part := chunk.Candidates[0].Content.Parts[0]
			complete_response = complete_response + part.Text
			fmt.Print(part.Text)
		}
	}

	setHistory(complete_response, true, true)

	fmt.Println() // Ensures terminal starts on a new line.

}
