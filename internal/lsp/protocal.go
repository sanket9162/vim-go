package lsp

import "encoding/json"

// Request represents a JSON-PRC 2.0 request message.
type Request struct {
	JsonRPC string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

// Response represents a JSON-RPC 2.0 response message.
type Response struct {
	JsonRPC string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *ResponseError  `json:"error,omitempty"`
}

// ReponseError represents the error fiel in a JSON-RPC 2.0 response.
type ResponseError struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}
