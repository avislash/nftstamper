package main

import (
	"context"
	"os"
	"os/signal"

	_ "github.com/avislash/nftstamper/ape"
	_ "github.com/avislash/nftstamper/cartel"
	"github.com/avislash/nftstamper/lib/log"
	"github.com/avislash/nftstamper/root"
)

func main() {
	log, _ := log.NewSugaredLogger()
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	root.Cmd.SetContext(ctx)

	if err := root.Cmd.Execute(); err != nil {
		log.Errorf("Exited with error: %s", err)
		os.Exit(1)
	}

	log.Info("Exited successfully")
}
