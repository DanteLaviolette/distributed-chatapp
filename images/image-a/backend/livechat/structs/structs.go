package structs

// Struct used in socket communication
type Message struct {
	Type    string `json:"type"`
	Content string `json:"content"`
	Subject string `json:"subject"`
}

type ChatMessage struct {
	Type    string `json:"type"`
	Subject string `json:"subject"`
	Message string `json:"message"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Ts      int64  `json:"ts"`
}

type AuthContext struct {
	Email    string
	Name     string
	SocketId string
}
