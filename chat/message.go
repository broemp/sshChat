package chat

import "time"

// How a message is saved/send
type Message struct {
	User       string
	Timestampt time.Time
	Text       string
}

// Currently not needed, but should later become an interface to have different storage backends
type MessageHandler struct{}

func NewMessageHandler() *MessageHandler {
	return &MessageHandler{}
}
