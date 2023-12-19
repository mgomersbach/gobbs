package bbs

import (
	"net"

	"gobbs/auth"
	"gobbs/config"

	"github.com/sirupsen/logrus"
)

// StartServer starts the BBS/Telnet server
func StartServer(cfg *config.Config, authMethod auth.Authenticator, log *logrus.Logger, stopChan <-chan struct{}) {
	// Use the address from the configuration
	address := cfg.BBS.Address
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}
	defer listener.Close()

	log.Printf("BBS/Telnet server listening on %s", address)

	// Continuously listen for incoming connections
	for {
		select {
		case <-stopChan:
			// Perform cleanup and stop the server
			log.Info("Shutting down server...")
			return
		default:
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("Failed to accept connection: %s", err)
				continue
			}

			// Handle each connection in a separate goroutine
			go handleConnection(conn, authMethod, log)
		}
	}
}

// handleConnection deals with an individual client connection
func handleConnection(conn net.Conn, authMethod auth.Authenticator, log *logrus.Logger) {
	defer conn.Close()

	if !frontDoor(conn, authMethod, log) {
		log.Info("FrontDoor access denied for connection")
		return
	}

	// Proceed with the rest of the BBS logic after successful FrontDoor validation
	// For example, navigating the menu, accessing forums, etc.
	// ...
}

// frontDoor handles the initial interaction with the user
func frontDoor(conn net.Conn, authMethod auth.Authenticator, log *logrus.Logger) bool {
	conn.Write([]byte("Welcome to the BBS! Please login.\n"))
	return true
}
