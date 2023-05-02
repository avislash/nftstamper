package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	_ "github.com/avislash/nftstamper/ape"
	_ "github.com/avislash/nftstamper/cartel"
	"github.com/avislash/nftstamper/root"
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	root.Cmd.SetContext(ctx)

	if err := root.Cmd.Execute(); err != nil {
		log.Println("Exited with error: %w")
		os.Exit(1)
	}

	log.Println("Exited successfully")
}
