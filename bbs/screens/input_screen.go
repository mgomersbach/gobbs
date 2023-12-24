package screens

import (
	"fmt"
	"regexp"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type TextInputScreen struct {
	screenManager *ScreenManager
	input         textinput.Model
}

func NewTextInputScreen(screenManager *ScreenManager) *TextInputScreen {
	input := textinput.New()
	input.Placeholder = "Enter text here..."
	input.Focus()
	return &TextInputScreen{
		screenManager: screenManager,
		input:         input,
	}
}

func (s *TextInputScreen) Update(msg tea.Msg) (Screen, tea.Cmd) {
	var cmd tea.Cmd

	// First, update the text input component
	s.input, cmd = s.input.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			// Check if the input has content and is focused
			if s.input.Value() != "" && s.input.Focused() {
				fmt.Printf("You entered: %s\n", s.input.Value())

				// Unfocus the input and return to the previous screen
				s.input.Blur()
				return nil, tea.Batch(tea.EnterAltScreen, tea.ClearScreen) // This triggers the transition back
			}
		}

	case string:
		// Ignore cursor position messages like "{0 2}"
		if matched, _ := regexp.MatchString(`^\{\d+ \d+\}$`, msg); !matched {
			// log position
			fmt.Printf("Cursor position: %s\n", msg)
		}
	}

	return s, cmd
}

func (s *TextInputScreen) View() string {
	// render input field
	return "\033[H" + s.input.View()
}

func (s *TextInputScreen) Init() tea.Cmd {
	return textinput.Blink
}
