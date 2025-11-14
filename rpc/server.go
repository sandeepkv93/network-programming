package rpc

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"sync"
)

// MathService provides mathematical operations
type MathService struct {
	mu sync.Mutex
}

// Args represents arguments for mathematical operations
type Args struct {
	A, B int
}

// Result represents the result of an operation
type Result struct {
	Value int
}

// Add adds two numbers
func (m *MathService) Add(args *Args, result *Result) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	result.Value = args.A + args.B
	log.Printf("Add(%d, %d) = %d\n", args.A, args.B, result.Value)
	return nil
}

// Subtract subtracts two numbers
func (m *MathService) Subtract(args *Args, result *Result) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	result.Value = args.A - args.B
	log.Printf("Subtract(%d, %d) = %d\n", args.A, args.B, result.Value)
	return nil
}

// Multiply multiplies two numbers
func (m *MathService) Multiply(args *Args, result *Result) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	result.Value = args.A * args.B
	log.Printf("Multiply(%d, %d) = %d\n", args.A, args.B, result.Value)
	return nil
}

// Divide divides two numbers
func (m *MathService) Divide(args *Args, result *Result) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if args.B == 0 {
		return fmt.Errorf("division by zero")
	}
	result.Value = args.A / args.B
	log.Printf("Divide(%d, %d) = %d\n", args.A, args.B, result.Value)
	return nil
}

// Server represents an RPC server
type Server struct {
	address  string
	listener net.Listener
	rpcServer *rpc.Server
}

// NewServer creates a new RPC server
func NewServer(address string) *Server {
	return &Server{
		address:   address,
		rpcServer: rpc.NewServer(),
	}
}

// RegisterService registers a service with the RPC server
func (s *Server) RegisterService(name string, service interface{}) error {
	return s.rpcServer.RegisterName(name, service)
}

// Start starts the RPC server
func (s *Server) Start() error {
	// Register default math service
	mathService := new(MathService)
	s.RegisterService("MathService", mathService)

	var err error
	s.listener, err = net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("failed to start RPC server: %v", err)
	}

	log.Printf("RPC Server listening on %s\n", s.address)

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v\n", err)
			continue
		}

		log.Printf("Client connected from %s\n", conn.RemoteAddr())
		go s.rpcServer.ServeConn(conn)
	}
}

// Stop stops the RPC server
func (s *Server) Stop() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}
