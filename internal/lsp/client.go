package lsp

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

type Client struct {
	cmd       *exec.Cmd
	stdin     io.WriteCloser
	stdout    io.ReadCloser
	reader    *bufio.Reader
	nextID    int
	idMu      sync.Mutex
	pending   map[int]chan *Response
	pendingMu sync.Mutex
	rootURI   string
}

// StartClient spawns the server binary and completes the initialization handshake.
func StartClient(binaryPath string, agrs []string, workspacePath string) (*Client, error) {
	cmd := exec.Command(binaryPath, agrs...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stdin pipe: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start server subprocess: %w", err)
	}

	c := &Client{
		cmd:     cmd,
		stdin:   stdin,
		stdout:  stdout,
		reader:  bufio.NewReader(stdout),
		pending: make(map[int]chan *Response),
		rootURI: "file://" + workspacePath,
	}

	// Start reading stdout in background goroutine
	go c.readLoop()

	// Perform initialize handshake
	if err := c.initialize(); err != nil {
		c.Close()
		return nil, fmt.Errorf("LSP initialization handshake failed: %w", err)
	}

	return c, nil

}

func (c *Client) Close() {
	_ = c.stdin.Close()
	_ = c.stdout.Close()
	_ = c.cmd.Process.Kill()
}

func (c *Client) getNextID() int {
	c.idMu.Lock()
	defer c.idMu.Unlock()

	c.nextID++
	return c.nextID
}

// sendRequest sends a JSON-RPC request and returns a channel to await the reponse.
func (c *Client) sendRequest(method string, params interface{}) (chan *Response, int, error) {
	id := c.getNextID()
	req := Request{
		JsonRPC: "2.0",
		ID:      id,
		Method:  method,
		Params:  params,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, 0, err
	}

	ch := make(chan *Response, 1)
	c.pendingMu.Lock()
	c.pending[id] = ch
	c.pendingMu.Unlock()

	if err := c.writeMessage(body); err != nil {
		c.pendingMu.Lock()
		delete(c.pending, id)
		c.pendingMu.Unlock()
		return nil, 0, err
	}

	return ch, 0, nil

}

func (c *Client) writeMessage(body []byte) error {
	var buf bytes.Buffer
	buf.WriteString("Content-Length: ")
	buf.WriteString(strconv.Itoa(len(body)))
	buf.WriteString("\r\n\r\n")
	buf.Write(body)

	_, err := c.stdin.Write(buf.Bytes())
	return err
}

// sendNotification sends an asynchronoys JSON-RPC notification
func (c *Client) sendNotification(method string, params interface{}) error {
	noti := Notification{
		JsonRPC: "2.0",
		Method:  method,
		Params:  params,
	}

	body, err := json.Marshal(noti)
	if err != nil {
		return err
	}

	return c.writeMessage(body)
}

func (c *Client) readLoop() {

	for {
		// Read Headers
		var contentLength int
		for {
			line, err := c.reader.ReadString('\n')
			if err != nil {
				return
			}
			line = strings.TrimSpace(line)
			if line == "" {
				// End of headers delimiter
				break
			}

			if strings.HasPrefix(line, "Content-Length:") {
				parts := strings.Split(line, ":")
				if len(parts) == 2 {
					contentLength, _ = strconv.Atoi(strings.TrimSpace(parts[1]))
				}
			}

		}

		if contentLength == 0 {
			continue
		}

		// Read body
		body := make([]byte, contentLength)
		if _, err := io.ReadFull(c.reader, body); err != nil {
			return
		}

		var msg Response
		if err := json.Unmarshal(body, &msg); err == nil && msg.ID != 0 {
			// Resolve request callback
			c.pendingMu.Lock()
			ch, ok := c.pending[msg.ID]
			if ok {
				ch <- &msg
				delete(c.pending, msg.ID)
			}
			c.pendingMu.Unlock()
		}
	}
}

func (c *Client) initialize() error {
	parms := InitializeParams{
		ProcessID: os.Getpid(),
		RootURI:   c.rootURI,
		Capabilities: ClientCapabilities{
			TextDocument: map[string]interface{}{
				"defintion": map[string]interface{}{
					"dynamicRegistration": true,
				},
				"hover": map[string]interface{}{
					"contentFormat": []string{"markdown", "plaintext"},
				},
			},
		},
	}

	ch, _, err := c.sendRequest("initialize", parms)
	if err != nil {
		return err
	}

	resp := <-ch
	if resp.Error != nil {
		return fmt.Errorf("server initialization error: %s", resp.Error.Message)
	}

	// Send initialized notification confirmation
	return c.sendNotification("initialized", map[string]interface{}{})

}

// NotifyDidOpen notifies the server that a file has been opened.
func (c *Client) NotifyDidOpen(uri, languageID, text string) error {
	params := DidOpenTextDocumentParams{
		TextDocument: TextDocumentItem{
			URI:        uri,
			LanguageID: languageID,
			Version:    1,
			Text:       text,
		},
	}
	return c.sendNotification("textDocument/didOpen", params)
}

// NotifyDidChange notifies the server of changes. We send full text changes for simplicity and safety.
func (c *Client) NotifyDidChange(uri string, version int, text string) error {
	params := DidChangeTextDocumentParams{
		TextDocument: VersionedTextDocumentIdentifier{
			URI:     uri,
			Version: version,
		},
		ContentChanges: []TextDocumentContentChangeEvent{
			{Text: text},
		},
	}
	return c.sendNotification("textDocument/didChange", params)
}

// NotifyDidSave notifies the server that the file was written to disk.
func (c *Client) NotifyDidSave(uri string) error {
	params := DidSaveTextDocumentParams{
		TextDocument: TextDocumentIdentifier{
			URI: uri,
		},
	}
	return c.sendNotification("textDocument/didSave", params)
}

type ServerConfig struct {
	Binary string
	Args   []string
}

// ServerRegistry maps file extensions to their corresponding language server binaries and options.
var ServerRegistry = map[string]ServerConfig{
	".go":   {Binary: "gopls", Args: []string{"serve"}},
	".rs":   {Binary: "rust-analyzer", Args: []string{}},
	".html": {Binary: "vscode-html-language-server", Args: []string{"--stdio"}},
	".css":  {Binary: "vscode-css-language-server", Args: []string{"--stdio"}},
	".js":   {Binary: "typescript-language-server", Args: []string{"--stdio"}},
}
