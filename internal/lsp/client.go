package lsp

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
