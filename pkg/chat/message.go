package chat

const SYSTEM_MESSAGE = "system"

type Message struct {
	Type   int    `json:"type"`
	Sender string `json:"sender"`
	Body   string `json:"body"`
}
