package pages

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/johnkhk/cli_chat_app/client/app"
)

type registerModel struct {
	focusIndex int
	inputs     []textinput.Model
	cursorMode cursor.Mode
	rpcClient  *app.RpcClient
}

// NewRegisterModel initializes the register component
func NewRegisterModel(rpcClient *app.RpcClient) registerModel {
	m := registerModel{
		inputs:     make([]textinput.Model, 2),
		rpcClient:  rpcClient,
		cursorMode: cursor.CursorBlink,
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

func (m registerModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m registerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case errMsg:
		m.rpcClient.Logger.Errorln("Error registering user:", msg.err)
		// TODO show some error message to the user
		return m, nil
	case registerRespMsg:
		// Go to log in page
		return NewLoginModel(m.rpcClient), nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		// Set focus to next input or button
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Handle Enter key when a button is focused
			if s == "enter" {
				if m.focusIndex == len(m.inputs) {
					// Submit button logic here
					fmt.Println("Form submitted!")
					username, pasword := m.inputs[0].Value(), m.inputs[1].Value()
					return m, registerUserCmd(m.rpcClient, username, pasword)
					// return m, tea.Quit
				} else if m.focusIndex == len(m.inputs)+1 {
					// Back button logic here
					return NewLandingModel(m.rpcClient), nil
				}
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs)+1 {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs) + 1
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *registerModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m registerModel) View() string {
	var b strings.Builder

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

	// Render help text
	// b.WriteString(helpStyle.Render("cursor mode is "))
	// b.WriteString(cursorModeHelpStyle.Render(m.cursorMode.String()))
	// b.WriteString(helpStyle.Render("\ntab: focus next • n: new boy • q: exit"))
	b.WriteString(helpStyle.Render("\nesc/ctrl+c: quit"))

	return b.String()
}

func registerUserCmd(rpcClient *app.RpcClient, username, password string) tea.Cmd {
	return func() tea.Msg {
		err := rpcClient.AuthClient.RegisterUser(username, password)
		if err != nil {
			return errMsg{err}
		}
		return registerRespMsg{}
	}
}

type registerRespMsg struct{}
