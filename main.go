package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/togethercomputer/together-kubelogin/pkg/di"
)

var version = "HEAD"

func main() {
	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()
	os.Exit(di.NewCmd().Run(ctx, os.Args, version))
}
