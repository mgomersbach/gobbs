package bbs

import (
	"fmt"
	"log"
	"net"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mgomersbach/go-telnet"
)

const (
	DO           = 253
	DONT         = 254
	Echo         = 1 // Echo
	IAC          = 255
	SB           = 250
	SE           = 240
	SuppressGA   = 3 // Suppress Go Ahead
	TerminalType = 24
	WILL         = 251
	WONT         = 252
)

// TelnetHandler defines a function that initializes the Bubble Tea model and returns it
// along with any additional program options for each Telnet session.
type TelnetHandler func(telnet.Context, *net.Conn, telnet.Writer, telnet.Reader) (tea.Model, []tea.ProgramOption)

// CustomTelnetHandler wraps a TelnetHandler to integrate it with the telnet.Handler interface.
type CustomTelnetHandler struct {
	handler TelnetHandler
	next    telnet.Handler
}

// negotiateCharacterMode sends commands to switch the client to character mode.
func negotiateCharacterMode(conn net.Conn) {
	conn.Write([]byte{IAC, WONT, SuppressGA})
	conn.Write([]byte{IAC, DONT, Echo})
}

func NewCustomTelnetHandler(handler TelnetHandler, next telnet.Handler) *CustomTelnetHandler {
	return &CustomTelnetHandler{
		handler: handler,
		next:    next,
	}
}

func (cth *CustomTelnetHandler) ServeTELNET(ctx telnet.Context, conn *net.Conn, writer telnet.Writer, reader telnet.Reader) {
	// Switch to character mode
	negotiateCharacterMode(*conn)
	// Pass the conn argument to the handler
	model, opts := cth.handler(ctx, conn, writer, reader)
	if model == nil {
		log.Println("No model provided for Bubble Tea program")
		return
	}

	opts = append(opts, tea.WithOutput(writer), tea.WithInput(reader), tea.WithAltScreen())

	p := tea.NewProgram(model, opts...)

	done := make(chan struct{})
	inputDone := make(chan struct{})

	// Run Bubble Tea program in a separate goroutine
	go func() {
		defer close(done)
		if _, err := p.Run(); err != nil {
			log.Printf("Bubble Tea program exited with error: %v", err)
		}
	}()

	// Read and pass input to Bubble Tea program
	go func() {
		buffer := make([]byte, 0, 2) // Allocate buffer with initial capacity
		for {
			byteBuf := make([]byte, 1)
			n, err := reader.Read(byteBuf)
			if err != nil {
				log.Printf("Read error: %v", err)
				break
			}

			if n > 0 {
				log.Printf("Read %d bytes: %q", n, byteBuf[:n])
			}

			buffer = append(buffer, byteBuf[:n]...)

			for len(buffer) > 0 {
				fmt.Printf("buffer larger than 0\n")
				keyMsg, bytesRead, ok := parseKeyFromBuffer(buffer)
				if !ok {
					break
				}
				p.Send(keyMsg)
				buffer = buffer[bytesRead:]
			}
		}
	}()
	// Wait for either the Bubble Tea program or the input loop to finish
	select {
	case <-done:
	case <-inputDone:
	}
}

func parseKeyFromBuffer(buffer []byte) (tea.KeyMsg, int, bool) {
	if len(buffer) < 1 {
		return tea.KeyMsg{}, 0, false
	}

	// Check for CR (Carriage Return)
	if buffer[0] == '\r' {
		// Check if next byte is LF or NUL (if available)
		if len(buffer) > 1 {
			if buffer[1] == '\n' || buffer[1] == '\x00' {
				return tea.KeyMsg{Type: tea.KeyEnter}, 2, true
			}
		} else {
			// Need one more byte to decide
			return tea.KeyMsg{}, 0, false
		}
	}

	// Check for Enter, Backspace, and ANSI cursor keys
	if keyMsg, bytesRead, ok := parseSpecialKeys(buffer); ok {
		return keyMsg, bytesRead, ok
	}

	// Default case: treat as a regular rune
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{rune(buffer[0])}}, 1, true
}

func parseSpecialKeys(buffer []byte) (tea.KeyMsg, int, bool) {
	// Enter key
	fmt.Printf("checking for special keys\n")
	if buffer[0] == '\n' {
		fmt.Printf("found normal enter key")
		return tea.KeyMsg{Type: tea.KeyEnter}, 1, true
	}
	if len(buffer) >= 2 && buffer[0] == '\r' && buffer[1] == '\n' {
		fmt.Printf("found windows enter key\n")
		return tea.KeyMsg{Type: tea.KeyEnter}, 2, true
	}

	// Backspace key
	if buffer[0] == '\b' || buffer[0] == '\x7f' {
		fmt.Printf("found backspace key\n")
		return tea.KeyMsg{Type: tea.KeyBackspace}, 1, true
	}

	// ANSI cursor keys
	if len(buffer) >= 3 && buffer[0] == '\x1b' && buffer[1] == '[' {
		fmt.Printf("found ANSI cursor key\n")
		switch buffer[2] {
		case 'A':
			return tea.KeyMsg{Type: tea.KeyUp}, 3, true
		case 'B':
			return tea.KeyMsg{Type: tea.KeyDown}, 3, true
		case 'C':
			return tea.KeyMsg{Type: tea.KeyRight}, 3, true
		case 'D':
			return tea.KeyMsg{Type: tea.KeyLeft}, 3, true
		}
	}

	return tea.KeyMsg{}, 0, false
}
