package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Starscade/Gaia/internal/ai"
	"github.com/Starscade/Gaia/internal/env"
	"github.com/Starscade/Gaia/internal/text"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	environment := env.Init()

	if environment.ApiKey == "" {
		log.Fatal(text.ERR_NO_API_KEY) // No key? Why continue?
	}

	agent := ai.Init(environment)
	ai.Ask(ctx, environment, agent)
}
