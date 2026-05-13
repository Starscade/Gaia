package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"google.golang.org/genai"
)

func askAi() {
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

	var builder strings.Builder

	for {
		chunk, err := stream.Next()
		if err != nil {
			// Check if stream is done
			break
		}

		if chunk != nil && len(chunk.Candidates) > 0 {
			candidate := chunk.Candidates[0]
			if candidate.Content != nil && len(candidate.Content.Parts) > 0 {
				text := candidate.Content.Parts[0].Text
				builder.WriteString(text)
				print(text)
			}
		}
	}

	setHistory(builder.String(), true, true)

	fmt.Println() // Ensures terminal starts on a new line.
}

func initAi() {
	ctx = context.Background()
	var err error
	client, err = genai.NewClient(ctx, &genai.ClientConfig{
		Backend: genai.BackendGeminiAPI,
		APIKey:  api_key,
	})

	exitOnErr(err)

	config = &genai.GenerateContentConfig{
		SafetySettings: []*genai.SafetySetting{
			{
				Category:  genai.HarmCategoryDangerousContent,
				Threshold: genai.HarmBlockThreshold(censor_rating),
			},
			{
				Category:  genai.HarmCategoryHarassment,
				Threshold: genai.HarmBlockThreshold(censor_rating),
			},
			{
				Category:  genai.HarmCategoryHateSpeech,
				Threshold: genai.HarmBlockThreshold(censor_rating),
			},
			{
				Category:  genai.HarmCategorySexuallyExplicit,
				Threshold: genai.HarmBlockThreshold(censor_rating),
			},
		},
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
}
