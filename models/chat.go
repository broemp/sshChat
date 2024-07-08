package models

import (
	"fmt"
	"strings"

	"github.com/broemp/sshChat/chat"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	errMsg error
)

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	textinputStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()
)

type chatModel struct {
	app         *App
	quit        bool
	viewport    viewport.Model
	textinput   textinput.Model
	senderStyle lipgloss.Style
	help        help.Model
	keys        *listKeyMap

	messages       []string
	messageHandler *chat.MessageHandler
	user           string

	ready bool
	err   error
}

func InitializeChatModel(quit bool) chatModel {
	ti := textinput.New()
	ti.Placeholder = "Send a message..."
	ti.Focus()
	ti.CharLimit = 280
	ti.Width = 100

	return chatModel{
		textinput:   ti,
		keys:        newListKeyMap(),
		help:        help.New(),
		quit:        quit,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err:         nil,
	}
}

func (m chatModel) Init() tea.Cmd {
	if m.quit {
		return tea.Quit
	}
	return textarea.Blink
}

func (m chatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textinput, tiCmd = m.textinput.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		textinputHeight := lipgloss.Height(m.textinputView())
		footerHeight := lipgloss.Height(m.helpView())
		verticalMarginHeight := headerHeight + footerHeight + textinputHeight

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.ready = true
			m.viewport.YPosition = headerHeight + 1
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.sendMessage):
			if m.textinput.Value() != "" {
				chatMsg := chat.Message{
					User: m.user,
					Text: m.textinput.Value(),
				}
				m.app.send(chatMsg)
				m.textinput, tiCmd = m.textinput.Update(msg)
			}
			m.textinput.Reset()
		}

		// If new Message is received, render it
	case chat.Message:
		m.messages = append(m.messages, m.senderStyle.Render(msg.User)+": "+msg.Text)
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
	if !m.ready {
		return "\n  Initializing..."
	}
	return fmt.Sprintf(
		"%s\n%s\n%s\n%s",
		m.headerView(),
		m.viewport.View(),
		m.textinputView(),
		m.helpView(),
	)
}

func (m chatModel) headerView() string {
	title := titleStyle.Render("KN2 Chat")
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m chatModel) helpView() string {
	return "\n" + m.help.ShortHelpView([]key.Binding{
		m.keys.sendMessage,
		m.keys.quit,
	})
}

func (m chatModel) textinputView() string {
	line := strings.Repeat("─", max(0, m.viewport.Width))
	return lipgloss.JoinVertical(lipgloss.Left, line, m.textinput.View())
}

func (m chatModel) handleCommand(cmd string) string {
	switch cmd {
	case "/msg":
	}
	return ""
}

// Helpers
type listKeyMap struct {
	sendMessage key.Binding
	quit        key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		sendMessage: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("Enter", "Send a message"),
		),
		quit: key.NewBinding(
			key.WithKeys("esc", "ctrl+c"),
			key.WithHelp("Ctrl+c/Esc", "Quit")),
	}
}
