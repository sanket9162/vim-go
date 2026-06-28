package lsp

import (
	"bufio"
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
