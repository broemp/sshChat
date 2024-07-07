package models

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	errMsg  error
	chatMsg struct {
		id   string
		text string
	}
)

type chatModel struct {
	app         *App
	viewport    viewport.Model
	messages    []string
	id          string
	textarea    textarea.Model
	senderStyle lipgloss.Style
	err         error
}

func InitializeChatModel() chatModel {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "| "
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(3)

	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = false

	vp := viewport.New(30, 5)
	vp.SetContent("Welcome to the KN2 Chat!")

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return chatModel{
		textarea:    ta,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err:         nil,
	}
}

func (m chatModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m chatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			m.app.send(chatMsg{
				id:   m.id,
				text: m.textarea.Value(),
			})
			m.textarea.Reset()
		}

	case chatMsg:
		m.messages = append(m.messages, m.senderStyle.Render(msg.id)+": "+msg.text)
		m.viewport.SetContent(strings.Join(m.messages, "\n"))
		m.viewport.GotoBottom()

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m chatModel) View() string {
	return fmt.Sprintf(
		"%s\n\n%s",
		m.viewport.View(),
		m.textarea.View(),
	) + "\n\n"
}
