package lsp

import "encoding/json"

// Request represents a JSON-PRC 2.0 request message.
type Request struct {
	JsonRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
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

// Notification represents a JSON-RPC 2.0 notification message.
type Notification struct {
	JsonRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

// InitializeParams represents the parameters for the 'initialize' request.
type InitializeParams struct {
	ProcessID    int                `json:"processId"`
	RootPath     string             `json:"rootPath,omitempty"`
	RootURI      string             `json:"rootUri,omitempty"`
	Capabilities ClientCapabilities `json:"capabilities"`
}

// ClientCapabilities represents editor/client support flags.
type ClientCapabilities struct {
	TextDocument interface{} `json:"textDocument,omitempty"`
	Workspace    interface{} `json:"workspace,omitempty"`
}

// InitializeResult represents the server capabilities returned on handshake.
type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
}

// ServerCapabilities defines features supported by the language server.
type ServerCapabilities struct {
	TextDocumentSync          interface{} `json:"textDocumentSync,omitempty"`
	DefinitionProvider        bool        `json:"definitionProvider,omitempty"`
	HoverProvider             bool        `json:"hoverProvider,omitempty"`
	CompletionProvider        interface{} `json:"completionProvider,omitempty"`
	SignatureHelpProvider     interface{} `json:"signatureHelpProvider,omitempty"`
	DocumentHighlightProvider bool        `json:"documentHighlightProvider,omitempty"`
}
