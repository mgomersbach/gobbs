package screens

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type MainMenuScreen struct {
	screenManager *ScreenManager
	menuItems     []string
	selectedIdx   int
}

// NewMainMenuScreen creates a new instance of MainMenuScreen.
func NewMainMenuScreen(screenManager *ScreenManager) *MainMenuScreen {
	return &MainMenuScreen{
		screenManager: screenManager,
		menuItems:     []string{"Menu Item 1", "Menu Item 2", "Menu Item 3"},
		selectedIdx:   0,
	}
}

// Init initializes the MainMenuScreen.
func (s *MainMenuScreen) Init() tea.Cmd {
	// Initialization logic, if any
	return nil
}

// Update processes key events and updates the screen state.
func (s *MainMenuScreen) Update(msg tea.Msg) (Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			s.selectedIdx--
			if s.selectedIdx < 0 {
				s.selectedIdx = len(s.menuItems) - 1
			}
		case "down", "j":
			s.selectedIdx++
			if s.selectedIdx >= len(s.menuItems) {
				s.selectedIdx = 0
			}
		case "enter":
			if s.selectedIdx == 0 {
				textInputScreen := NewTextInputScreen(s.screenManager)
				s.screenManager.SetScreen(textInputScreen)
				return textInputScreen, nil
			}
			// Handle other menu items similarly
		}
	}
	return s, nil
}

// View renders the main menu screen.
func (s *MainMenuScreen) View() string {
	var view strings.Builder

	for i, item := range s.menuItems {
		if i == s.selectedIdx {
			view.WriteString(fmt.Sprintf("> %s\n", item)) // Highlight the selected item
		} else {
			view.WriteString(fmt.Sprintf("  %s\n", item))
		}
	}

	return view.String()
}
