package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"flower/api/internal/app"
	"flower/api/internal/platform/config"
)

func main() {
	command, operand, err := parseCommand(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	switch command {
	case "migrate":
		app.RunMigrate()
		return
	case "rollback":
		app.RunRollback()
		return
	case "force":
		app.RunForce(operand)
		return
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create application: %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigChan
		log.Printf("Received shutdown signal (%s), beginning graceful shutdown...", sig)
		application.Cancel()
	}()

	if err := application.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func parseCommand(args []string) (string, string, error) {
	if len(args) < 2 {
		return "", "", nil
	}

	switch args[1] {
	case "migrate", "rollback":
		return args[1], "", nil
	case "force":
		if len(args) < 3 || args[2] == "" {
			return "", "", fmt.Errorf("usage: server force <version>")
		}
		return "force", args[2], nil
	default:
		return "", "", fmt.Errorf("unknown command %q", args[1])
	}
}
