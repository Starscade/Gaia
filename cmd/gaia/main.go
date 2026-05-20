package main

import (
	"log"

	"github.com/Starscade/Gaia/internal/ai"
	"github.com/Starscade/Gaia/internal/env"
	"github.com/Starscade/Gaia/internal/text"
)

func main() {
	environment := env.Init()

	if environment.ApiKey == "" {
		log.Fatal(text.ERR_NO_API_KEY) // No key? Why continue?
	}

	agent := ai.Init(environment)
	ai.Ask(environment, agent)
}
