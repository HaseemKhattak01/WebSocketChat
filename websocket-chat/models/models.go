package models

type Message struct {
	Username string `json:"username"`
	Type     string `json:"type"`
	Content  string `json:"content"`
}
