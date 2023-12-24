package bbs

import (
	"gobbs/auth"
	"gobbs/bbs/screens"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sirupsen/logrus"
)

type bbsModel struct {
	Username      string
	Password      string
	AuthStep      int
	ErrMsg        string
	AuthMethod    auth.Authenticator
	Log           *logrus.Logger
	ScreenManager *screens.ScreenManager
	UserSession   UserSession
}

// UserSession represents a user's session data.
type UserSession struct {
	// Add fields like username, isAuthenticated, etc.
}

func NewBBSModel(authMethod auth.Authenticator, log *logrus.Logger) *bbsModel {
	screenManager := screens.NewScreenManager()
	mainMenuScreen := screens.NewMainMenuScreen(screenManager)
	screenManager.AddScreen("mainMenu", mainMenuScreen)
	screenManager.SetScreen(mainMenuScreen)

	return &bbsModel{
		ScreenManager: screenManager,
		UserSession:   UserSession{}, // Initialize user session
		AuthMethod:    authMethod,
		Log:           log,
	}
}

// Init is called once when the program starts.
func (m *bbsModel) Init() tea.Cmd {
	// Initialize any resources or return an initial command here.
	return nil
}

// Update is called in response to messages and returns a modified model.
func (m *bbsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	if m.ScreenManager != nil {
		m.ScreenManager.Update(msg)
	}
	return m, cmd
}

// View renders the current screen of the UI.
func (m *bbsModel) View() string {
	if m.ScreenManager != nil {
		return m.ScreenManager.View()
	}
	return "Error: ScreenManager not initialized."
}
