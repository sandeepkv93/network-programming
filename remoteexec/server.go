package remoteexec

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// Server represents a remote execution server
type Server struct {
	address      string
	listener     net.Listener
	allowedCmds  map[string]bool
	allowAll     bool
	activeExecs  map[string]*ExecutionContext
	mu           sync.RWMutex
	authRequired bool
	authToken    string
}

// ExecutionContext represents an active command execution
type ExecutionContext struct {
	ID        string
	Command   string
	StartTime time.Time
	cmd       *exec.Cmd
	mu        sync.Mutex
}

// CommandRequest represents a command execution request
type CommandRequest struct {
	ID      string   `json:"id"`
	Command string   `json:"command"`
	Args    []string `json:"args"`
	Token   string   `json:"token,omitempty"`
}

// CommandResponse represents a command execution response
type CommandResponse struct {
	ID       string `json:"id"`
	Success  bool   `json:"success"`
	Output   string `json:"output"`
	Error    string `json:"error,omitempty"`
	ExitCode int    `json:"exit_code"`
	Duration string `json:"duration"`
}

// NewServer creates a new remote execution server
func NewServer(address string, authToken string, allowedCmds []string) *Server {
	server := &Server{
		address:      address,
		allowedCmds:  make(map[string]bool),
		activeExecs:  make(map[string]*ExecutionContext),
		authRequired: authToken != "",
		authToken:    authToken,
	}

	if len(allowedCmds) == 0 {
		server.allowAll = true
	} else {
		for _, cmd := range allowedCmds {
			server.allowedCmds[cmd] = true
		}
	}

	return server
}

// Start starts the remote execution server
func (s *Server) Start() error {
	var err error
	s.listener, err = net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}

	log.Printf("Remote Execution Server listening on %s\n", s.address)
	if s.authRequired {
		log.Println("Authentication: ENABLED")
	} else {
		log.Println("Authentication: DISABLED (WARNING: Insecure)")
	}

	if s.allowAll {
		log.Println("Command filtering: DISABLED (WARNING: All commands allowed)")
	} else {
		log.Printf("Allowed commands: %v\n", s.getAllowedCommands())
	}

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v\n", err)
			continue
		}

		go s.handleConnection(conn)
	}
}

// handleConnection handles a client connection
func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Printf("Client connected from %s\n", conn.RemoteAddr())

	reader := bufio.NewReader(conn)

	for {
		// Read command request
		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Printf("Error reading request: %v\n", err)
			}
			break
		}

		var req CommandRequest
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			log.Printf("Invalid request format: %v\n", err)
			s.sendError(conn, "", "Invalid request format")
			continue
		}

		// Authenticate if required
		if s.authRequired && req.Token != s.authToken {
			log.Printf("Authentication failed for request %s\n", req.ID)
			s.sendError(conn, req.ID, "Authentication failed")
			continue
		}

		// Check if command is allowed
		if !s.allowAll && !s.allowedCmds[req.Command] {
			log.Printf("Command not allowed: %s\n", req.Command)
			s.sendError(conn, req.ID, fmt.Sprintf("Command not allowed: %s", req.Command))
			continue
		}

		// Execute command
		log.Printf("Executing command [%s]: %s %v\n", req.ID, req.Command, req.Args)
		response := s.executeCommand(req)

		// Send response
		respJSON, _ := json.Marshal(response)
		conn.Write(append(respJSON, '\n'))
	}

	log.Printf("Client disconnected from %s\n", conn.RemoteAddr())
}

// executeCommand executes a command and returns the response
func (s *Server) executeCommand(req CommandRequest) CommandResponse {
	startTime := time.Now()

	// Create execution context
	ctx := &ExecutionContext{
		ID:        req.ID,
		Command:   req.Command,
		StartTime: startTime,
	}

	s.mu.Lock()
	s.activeExecs[req.ID] = ctx
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.activeExecs, req.ID)
		s.mu.Unlock()
	}()

	// Execute command
	cmd := exec.Command(req.Command, req.Args...)
	ctx.cmd = cmd

	output, err := cmd.CombinedOutput()
	duration := time.Since(startTime)

	response := CommandResponse{
		ID:       req.ID,
		Success:  err == nil,
		Output:   string(output),
		Duration: duration.String(),
	}

	if err != nil {
		response.Error = err.Error()
		if exitErr, ok := err.(*exec.ExitError); ok {
			response.ExitCode = exitErr.ExitCode()
		} else {
			response.ExitCode = -1
		}
	} else {
		response.ExitCode = 0
	}

	return response
}

// sendError sends an error response to the client
func (s *Server) sendError(conn net.Conn, id, errMsg string) {
	response := CommandResponse{
		ID:       id,
		Success:  false,
		Error:    errMsg,
		ExitCode: -1,
	}

	respJSON, _ := json.Marshal(response)
	conn.Write(append(respJSON, '\n'))
}

// Stop stops the remote execution server
func (s *Server) Stop() error {
	// Kill all active executions
	s.mu.Lock()
	for _, ctx := range s.activeExecs {
		ctx.mu.Lock()
		if ctx.cmd != nil && ctx.cmd.Process != nil {
			ctx.cmd.Process.Kill()
		}
		ctx.mu.Unlock()
	}
	s.mu.Unlock()

	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

// GetActiveExecutions returns the number of active command executions
func (s *Server) GetActiveExecutions() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.activeExecs)
}

// getAllowedCommands returns a list of allowed commands
func (s *Server) getAllowedCommands() []string {
	var cmds []string
	for cmd := range s.allowedCmds {
		cmds = append(cmds, cmd)
	}
	return cmds
}

// KillExecution kills a running command by ID
func (s *Server) KillExecution(id string) error {
	s.mu.RLock()
	ctx, exists := s.activeExecs[id]
	s.mu.RUnlock()

	if !exists {
		return fmt.Errorf("execution not found: %s", id)
	}

	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	if ctx.cmd != nil && ctx.cmd.Process != nil {
		return ctx.cmd.Process.Kill()
	}

	return fmt.Errorf("no process to kill")
}

// AddAllowedCommand adds a command to the allowed list
func (s *Server) AddAllowedCommand(cmd string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.allowedCmds[cmd] = true
	s.allowAll = false
}

// RemoveAllowedCommand removes a command from the allowed list
func (s *Server) RemoveAllowedCommand(cmd string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.allowedCmds, cmd)
}

// IsCommandAllowed checks if a command is allowed
func (s *Server) IsCommandAllowed(cmd string) bool {
	if s.allowAll {
		return true
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.allowedCmds[cmd]
}
