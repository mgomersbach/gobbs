package screens

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Screen represents different screens/views in your BBS application.
type Screen interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (Screen, tea.Cmd)
	View() string
}

type ScreenManager struct {
	screens       map[string]Screen
	currentScreen Screen
	prevScreen    Screen
}

func NewScreenManager() *ScreenManager {
	return &ScreenManager{
		screens: make(map[string]Screen),
	}
}

func (sm *ScreenManager) AddScreen(name string, screen Screen) {
	sm.screens[name] = screen
}

// SetScreen changes the current screen
func (sm *ScreenManager) SetScreen(screen Screen) {
	if screen != nil {
		sm.prevScreen = sm.currentScreen
		sm.currentScreen = screen
	}
}

// GoBack switches back to the previous screen
func (sm *ScreenManager) GoBack() {
	if sm.prevScreen != nil {
		sm.currentScreen, sm.prevScreen = sm.prevScreen, nil
	}
}

// Update delegates the update to the current screen
func (sm *ScreenManager) Update(msg tea.Msg) (Screen, tea.Cmd) {
	if sm.currentScreen == nil {
		// Set default screen or handle error

	}

	newScreen, cmd := sm.currentScreen.Update(msg)

	// Check for nil to go back to the previous screen
	if newScreen == nil {
		sm.GoBack()
		return sm.currentScreen, nil
	}

	// Handle transition to a new screen
	if newScreen != sm.currentScreen {
		sm.SetScreen(newScreen)
	}

	return newScreen, cmd
}

// View renders the view of the current screen
func (sm *ScreenManager) View() string {
	if sm.currentScreen == nil {
		// Handle nil currentScreen, perhaps log an error or return a default view
		return "Error: No current screen is set."
	}
	return sm.currentScreen.View()
}
