package websocket

import "encoding/json"

type AuthRequest struct {
	Token string `json:"token"`
	Type  string `json:"type"`
	Data  json.RawMessage
}
