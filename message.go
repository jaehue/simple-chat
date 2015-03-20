package main

import (
	"time"
)

const (
	TextMessage MessageType = iota
	RoomAction
)

type MessageType int

type Message struct {
	MessageType
	RoomId    int
	Content   string
	CreatedAt time.Time
}
