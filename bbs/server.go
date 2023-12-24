package bbs

import (
	"gobbs/auth"
	"gobbs/config"
	"net"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mgomersbach/go-telnet"
	"github.com/reiver/go-oi"

	"github.com/sirupsen/logrus"
)

type bbsHandler struct {
	cfg        *config.Config
	authMethod auth.Authenticator
	log        *logrus.Logger
}

// bbsHandler is a custom TELNET handler for your BBS server.
func StartServer(cfg *config.Config, authenticator auth.Authenticator, log *logrus.Logger, stopChan chan struct{}) {
	telnetHandler := func(ctx telnet.Context, conn *net.Conn, writer telnet.Writer, reader telnet.Reader) (tea.Model, []tea.ProgramOption) {
		model := NewBBSModel(authenticator, log) // Initialize your Bubble Tea model here
		opts := []tea.ProgramOption{}
		return model, opts
	}

	customHandler := NewCustomTelnetHandler(telnetHandler, nil) // Use 'nil' if there's no next handler

	server := &telnet.Server{
		Addr:    cfg.BBS.Address,
		Handler: customHandler,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal("Server failed:", err)
		}
	}()

	<-stopChan // Wait for a signal to stop the server
	// Perform server shutdown logic if necessary
}

// frontDoor handles the initial interaction with the user
func frontDoor(writer telnet.Writer, reader telnet.Reader, authMethod auth.Authenticator, log *logrus.Logger, p *tea.Program) bool {
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
