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

	environment, err := env.Init()

	if err != nil {
		log.Fatal(err)
	}

	if environment.ApiKey == "" {
		log.Fatal(text.ERR_NO_API_KEY) // No key? Why continue?
	}

	if environment.Prompt == "" {
		os.Exit(1)
	}

	agent, _ := ai.Init(*environment)
	ai.Ask(ctx, *environment, *agent)
}
