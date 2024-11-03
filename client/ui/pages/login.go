package pages

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/johnkhk/cli_chat_app/client/app"
	"github.com/johnkhk/cli_chat_app/client/ui/ascii"
)

type loginModel struct {
	focusIndex int
	inputs     []textinput.Model
	buttons    []string
	cursorMode cursor.Mode
	rpcClient  *app.RpcClient
	errorMsg   string
}

// NewLoginModel initializes the login component
func NewLoginModel(rpcClient *app.RpcClient) loginModel {
	m := loginModel{
		inputs:     make([]textinput.Model, 2),
		buttons:    []string{"Submit", "Back"}, // Initialize buttons
		cursorMode: cursor.CursorBlink,
		rpcClient:  rpcClient,
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Username"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
			t.CharLimit = 64
		case 1:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
			t.CharLimit = 64
		}

		m.inputs[i] = t
	}

	return m
}

func (m loginModel) Init() tea.Cmd {
	return textinput.Blink
}
func (m loginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.rpcClient.Logger.Infoln("Received message of type:", fmt.Sprintf("%T", msg)) // Log the message type

	switch msg := msg.(type) {
	case errMsg:
		m.rpcClient.Logger.Errorln("Error logging in user:", msg.err)
		m.errorMsg = msg.err.Error()
		return m, nil
	case logInRespMsg:
		return NewFriendManagementModel(m.rpcClient), tea.WindowSize()
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Handle Enter key when a button is focused
			if s == "enter" {
				if m.focusIndex == len(m.inputs) {
					// Log in user
					return m, logInUserCmd(m.rpcClient, m.inputs[0].Value(), m.inputs[1].Value())
				} else if m.focusIndex == len(m.inputs)+1 {
					// Go back to the landing page
					return NewLandingModel(m.rpcClient), nil
				}
			}

			// Update focus index without cycling
			if s == "up" || s == "shift+tab" {
				if m.focusIndex > 0 {
					m.focusIndex--
				}
			} else if s == "down" || s == "tab" {
				if m.focusIndex < len(m.inputs)+len(m.buttons)-1 {
					m.focusIndex++
				}
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i < len(m.inputs); i++ {
				if i == m.focusIndex {
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
				} else {
					m.inputs[i].Blur()
					m.inputs[i].PromptStyle = noStyle
					m.inputs[i].TextStyle = noStyle
				}
			}

			return m, tea.Batch(cmds...)
		}
	}

	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *loginModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m loginModel) View() string {
	var b strings.Builder

	// Render the title
	b.WriteString(logoStyle.Render(ascii.Logo))
	b.WriteString("\n\n")
	b.WriteString(titleStyle.Render("Login"))
	b.WriteString("\n\n")

	// Render the inputs
	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	// Render the submit and back buttons
	submitButton := &blurredSubmitButton
	backButton := &blurredBackButton

	if m.focusIndex == len(m.inputs) {
		submitButton = &focusedSubmitButton
	} else if m.focusIndex == len(m.inputs)+1 {
		backButton = &focusedBackButton
	}

	fmt.Fprintf(&b, "\n\n%s\n\n", *submitButton)
	fmt.Fprintf(&b, "%s\n\n", *backButton)

	// Render error message if exists
	if m.errorMsg != "" {
		b.WriteString(errorMsgStyle.Render(m.errorMsg))
	}

	b.WriteString(helpStyle.Render("\nesc/ctrl+c: quit"))

	return b.String()
}

func logInUserCmd(rpcClient *app.RpcClient, username, password string) tea.Cmd {
	return func() tea.Msg {
		err, _ := rpcClient.AuthClient.LoginUser(username, password)
		if err != nil {
			return errMsg{err}
		}
		return logInRespMsg{}
	}
}

type logInRespMsg struct{}
