package main

import (
	"fmt"
	"gobbs/auth"
	"gobbs/auth/methods/dbauth"
	"gobbs/bbs"
	"gobbs/config"
	"gobbs/database"
	"gobbs/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Process command line flags
	args := config.ProcessFlags()

	// Load configuration
	cfg, err := config.LoadConfig(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize the logger with the loaded configuration
	log := logger.InitializeLogger(cfg)

	// Use the logger for any subsequent error handling
	db, err := database.Connect(cfg, log)
	if err != nil {
		log.WithError(err).Fatal("Failed to connect to database")
	}

	// Initialize the database tables and predefined users
	if err := database.InitializeDatabase(db, cfg); err != nil {
		log.WithError(err).Fatal("Failed to initialize database")
	}

	// Initialize the appropriate authenticator
	var authenticator auth.Authenticator
	switch cfg.Authentication.Method {
	case "db":
		authenticator = dbauth.NewDBAuth(db)
	// ... other authentication methods ...
	default:
		log.Fatal("Unsupported authentication method")
	}

	// Start the BBS server
	stopChan := make(chan struct{})
	go bbs.StartServer(cfg, authenticator, log, stopChan)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	// Signal received, notify the server to stop
	close(stopChan)

	// Perform any additional cleanup...

	log.Info("Application exiting")

}
