package bbs

import (
	"gobbs/auth"
	"gobbs/config"

	"github.com/reiver/go-oi"
	"github.com/reiver/go-telnet"

	"github.com/sirupsen/logrus"
)

// bbsHandler is a custom TELNET handler for your BBS server.
type bbsHandler struct {
	cfg        *config.Config
	authMethod auth.Authenticator
	log        *logrus.Logger
}

// ServeTELNET handles the Telnet connection for the BBS.
func (handler bbsHandler) ServeTELNET(ctx telnet.Context, writer telnet.Writer, reader telnet.Reader) {
	if !frontDoor(writer, reader, handler.authMethod, handler.log) {
		handler.log.Info("FrontDoor access denied for connection")
		return
	}

	// User is authenticated; proceed with BBS main logic
	oi.LongWrite(writer, []byte("\nWelcome to the BBS! Type 'exit' to disconnect.\n"))

	for {
		oi.LongWrite(writer, []byte("\n> ")) // Prompt
		command := readLine(reader)

		// Process commands
		if command == "exit" {
			break
		}

		// Handle other commands here
		oi.LongWrite(writer, []byte("You typed: "+command+"\n"))
	}
}

// StartServer starts the BBS/Telnet server
func StartServer(cfg *config.Config, authMethod auth.Authenticator, log *logrus.Logger, stopChan <-chan struct{}) {
	// Define the server address
	serverAddr := cfg.BBS.Address
	log.Printf("BBS/Telnet server listening on %s", serverAddr)

	// Create and start the server
	srv := &telnet.Server{
		Addr:    serverAddr,
		Handler: bbsHandler{cfg, authMethod, log},
	}

	// Start the server in a separate goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Errorf("Failed to start Telnet server: %s", err)
		}
	}()

	// Wait for stop signal to gracefully shut down the server
	<-stopChan
	log.Info("Server is stopping")
}

// frontDoor handles the initial interaction with the user
func frontDoor(writer telnet.Writer, reader telnet.Reader, authMethod auth.Authenticator, log *logrus.Logger) bool {
	oi.LongWrite(writer, []byte("Welcome to GoBBS!\r\nLogin\r\nUsername: "))

	username := readLine(reader)
	oi.LongWrite(writer, []byte("Password: "))
	password := readLine(reader)

	authenticated, err := authMethod.Authenticate(username, password)
	if err != nil {
		log.Error("Authentication error:", err)
		return false
	}

	return authenticated
}

// readLine reads a line of input from the Telnet reader and trims CR and LF characters
func readLine(reader telnet.Reader) string {
	var buffer [1]byte
	var line []byte
	newlineEncountered := false

	for {
		_, err := reader.Read(buffer[:])
		if err != nil {
			break
		}

		// Check for line feed or carriage return
		if buffer[0] == '\n' || buffer[0] == '\r' {
			if newlineEncountered {
				// If we've already encountered a newline, break on subsequent newlines
				break
			}
			newlineEncountered = true
			continue
		}

		// If not newline, reset the flag and append to line
		newlineEncountered = false
		line = append(line, buffer[0])
	}

	return string(line)
}
