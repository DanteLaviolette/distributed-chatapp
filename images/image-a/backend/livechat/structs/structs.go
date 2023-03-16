package structs

// Struct used in socket communication
type Message struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}
