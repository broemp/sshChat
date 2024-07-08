package chat

import "time"

type Message struct {
	User       string
	Timestampt time.Time
	Text       string
}

type MessageHandler struct{}

func NewMessageHandler() *MessageHandler {
	return &MessageHandler{}
}
