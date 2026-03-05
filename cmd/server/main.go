package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"scaffold/cmd/server/bootstrap"
	"scaffold/internal/pkg/logger"
)

func main() {
	h, err := bootstrap.Bootstrap()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to bootstrap server: %v\n", err)
		os.Exit(1)
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		logger.Info("Shutting down server...")
		bootstrap.Cleanup()
		os.Exit(0)
	}()

	ctx := context.Background()
	_ = ctx

	logger.Info("Starting scaffold server...")
	h.Spin()
}
