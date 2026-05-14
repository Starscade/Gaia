package main

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/genai"
)

func askAi() {
	stream := client.Models.GenerateContentStream(
		ctx,
		agent_model,
		prompt_history,
		config,
	)

	// Print response as it comes ...

	var response_buffer strings.Builder

	for chunk, err := range stream {
		exitOnErr(err)

		if chunk != nil && len(chunk.Candidates) > 0 {
			for _, candidate := range chunk.Candidates {

				if candidate.FinishReason == genai.FinishReasonSafety {
					fmt.Print(PRINT_CENSORED)
				}

				if candidate.Content != nil && len(candidate.Content.Parts) > 0 {
					for _, part := range candidate.Content.Parts {

						text := part.Text
						response_buffer.WriteString(text)
						fmt.Print(text)

					}
				}

			}
		}
	}

	setHistory(response_buffer.String(), true, true)

	fmt.Println() // Ensure terminal returns on a new line.
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
