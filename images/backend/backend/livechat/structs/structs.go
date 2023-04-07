package structs

import (
	"go.violettedev.com/eecs4222/shared/database/schema"
)

// Struct used in socket communication
type Message struct {
	Type    string `json:"type"`
	Content string `json:"content"`
	Subject string `json:"subject"`
}

// Struct representing a chat message
type ChatMessage struct {
	Type string `json:"type"`
	schema.MessageSchema
}

type UserCountMessage struct {
	Type            string `json:"type"`
	AuthorizedUsers int    `json:"authorizedUsers"`
	AnonymousUsers  int    `json:"anonymousUsers"`
}

type AuthContext struct {
	Email    string
	Name     string
	SocketId string
	UserId   string
}
