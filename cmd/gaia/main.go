package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Starscade/Gaia/internal/ai"
	"github.com/Starscade/Gaia/internal/env"
	"github.com/Starscade/Gaia/internal/text"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	environment, err := env.Init()
	if err != nil {
		return err
	}
	defer environment.Db.Close()

	if environment.ApiKey == "" {
		return fmt.Errorf("%s", text.ErrNoApiKey)
	}

	if environment.Prompt == "" {
		return nil
	}

	agent, err := ai.Init(*environment)
	if err != nil {
		return err
	}

	return ai.Ask(ctx, *environment, *agent)
}
