package models

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/broemp/sshChat/chat"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	"github.com/muesli/termenv"
)

const (
	host = "localhost"
	port = "23234"
)

var users = map[string]string{
	"broemp": "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIPneKBhbjx1rVlhaNehDGsOBAh5r0vupuyTyQ+luSOVw ",
}

func NewApp() *App {
	a := new(App)
	a.messageHandler = chat.NewMessageHandler()
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithPublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
			return key.Type() == "ssh-ed25519"
		}),
		wish.WithMiddleware(
			bubbletea.MiddlewareWithProgramHandler(a.ProgramHandler, termenv.ANSI256),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Error("Could not start server", "error", err)
	}

	a.Server = s
	return a
}

// app contains a wish server and the list of running programs.
type App struct {
	*ssh.Server
	progs          []*tea.Program
	messageHandler *chat.MessageHandler
}

// send dispatches a message to all running programs.
func (a *App) send(msg tea.Msg) {
	for _, p := range a.progs {
		go p.Send(msg)
	}
}

func (a *App) Start() {
	var err error
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting SSH server", "host", host, "port", port)
	go func() {
		if err = a.ListenAndServe(); err != nil {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := a.Shutdown(ctx); err != nil {
		log.Error("Could not stop server", "error", err)
	}
}

func (a *App) ProgramHandler(s ssh.Session) *tea.Program {
	quit := false
	if _, ok := users[s.User()]; ok {
		pk := string(s.PublicKey().Marshal())
		println(pk)
		for _, pubkey := range users {
			parsed, _, _, _, _ := ssh.ParseAuthorizedKey(
				[]byte(pubkey),
			)
			if !ssh.KeysEqual(s.PublicKey(), parsed) {
				wish.Fatalln(s, "Hey, I don't know who you are!")
				quit = true
			}
		}
	}

	model := InitializeChatModel(quit)
	model.app = a
	model.user = s.User()
	model.messageHandler = chat.NewMessageHandler()
	pOptions := bubbletea.MakeOptions(s)
	pOptions = append(pOptions, tea.WithAltScreen())
	pOptions = append(pOptions, tea.WithMouseCellMotion())

	p := tea.NewProgram(model, pOptions...)
	a.progs = append(a.progs, p)

	return p
}
