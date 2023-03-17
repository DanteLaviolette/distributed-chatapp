package structs

// Struct used in socket communication
type Message struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

type ChatMessage struct {
	Type    string `json:"type"`
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
