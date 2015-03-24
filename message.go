package main

import (
	"time"
)

const (
	TextMessage MessageType = iota
	RoomAction
	ChatError
)

type MessageType int

type Message struct {
	MessageType
	RoomId    int
	Content   string
	CreatedAt time.Time
}

type Messager interface {
	ReadChatMessage(*Message) error
	WriteChatMessage(*Message) error
	Close() error
}
