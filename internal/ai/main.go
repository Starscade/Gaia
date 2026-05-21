package ai

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/genai"

	"github.com/Starscade/Gaia/internal/env"
	"github.com/Starscade/Gaia/internal/sql"
	"github.com/Starscade/Gaia/internal/text"
)

type Agent struct {
	Client *genai.Client
	Config *genai.GenerateContentConfig
}

func Ask(ctx context.Context, environment env.Environment, agent Agent) {
	stream := agent.Client.Models.GenerateContentStream(
		ctx,
		environment.AgentModel,
		environment.PromptHistory,
		agent.Config,
	)

	// Print response as it comes ...

	var response_buffer strings.Builder

	for chunk, err := range stream {
		if err != nil {
			sql.InsertMessage(environment.DbFile, response_buffer.String(), true, true)
			break
		}

		if chunk != nil && len(chunk.Candidates) > 0 {
			for _, candidate := range chunk.Candidates {

				if candidate.FinishReason != "" && candidate.FinishReason != genai.FinishReasonStop {
					switch candidate.FinishReason {
					case genai.FinishReasonMaxTokens:
						fmt.Print(text.PRINT_TOKENS_EXHAUSTED)
					case genai.FinishReasonRecitation:
						fmt.Print(text.PRINT_COPYRIGHTED)
					case genai.FinishReasonSafety:
						fmt.Print(text.PRINT_CENSORED)
					default:
						fmt.Print(text.PRINT_ERR)
					}
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

	sql.InsertMessage(environment.DbFile, response_buffer.String(), true, true)

	fmt.Println() // Ensure terminal returns on a new line.
}

func Init(environment env.Environment) (*Agent, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		Backend: genai.BackendGeminiAPI,
		APIKey:  environment.ApiKey,
	})

	if err != nil {
		return nil, err
	}

	config := &genai.GenerateContentConfig{
		SafetySettings: []*genai.SafetySetting{
			{
				Category:  genai.HarmCategoryDangerousContent,
				Threshold: genai.HarmBlockThreshold(environment.CensorRating),
			},
			{
				Category:  genai.HarmCategoryHarassment,
				Threshold: genai.HarmBlockThreshold(environment.CensorRating),
			},
			{
				Category:  genai.HarmCategoryHateSpeech,
				Threshold: genai.HarmBlockThreshold(environment.CensorRating),
			},
			{
				Category:  genai.HarmCategorySexuallyExplicit,
				Threshold: genai.HarmBlockThreshold(environment.CensorRating),
			},
		},
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{
				&genai.Part{
					Text: environment.AgentPersona,
				},
			},
		},
		ThinkingConfig: &genai.ThinkingConfig{
			ThinkingLevel: genai.ThinkingLevel(environment.AgentIntellect),
		},
		Tools: []*genai.Tool{
			{
				GoogleSearch: &genai.GoogleSearch{},
			},
		},
	}

	a := Agent{
		Client: client,
		Config: config,
	}

	return &a, nil

}
