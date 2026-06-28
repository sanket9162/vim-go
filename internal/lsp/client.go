package lsp

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
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
