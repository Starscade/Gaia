package main

import (
	"github.com/Starscade/Gaia/internal/ai"
	"github.com/Starscade/Gaia/internal/env"
)

func main() {
	environment := env.Init()
	agent := ai.Init(environment)
	ai.Ask(environment, agent)
}
