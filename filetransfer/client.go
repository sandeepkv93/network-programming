package filetransfer

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
)

// Client represents a file transfer client
type Client struct {
	address string
}

// NewClient creates a new file transfer client
func NewClient(address string) *Client {
	return &Client{
		address: address,
	}
}

// SendFile sends a file to the server
func (c *Client) SendFile(filePath string) error {
	// Open file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Get file info
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %v", err)
	}

	// Connect to server
	conn, err := net.Dial("tcp", c.address)
	if err != nil {
		return fmt.Errorf("failed to connect: %v", err)
	}
	defer conn.Close()

	// Send filename length and filename
	filename := filepath.Base(filePath)
	filenameLen := uint32(len(filename))
	if err := binary.Write(conn, binary.BigEndian, filenameLen); err != nil {
		return fmt.Errorf("failed to send filename length: %v", err)
	}

	if _, err := conn.Write([]byte(filename)); err != nil {
		return fmt.Errorf("failed to send filename: %v", err)
	}

	// Send file size
	fileSize := fileInfo.Size()
	if err := binary.Write(conn, binary.BigEndian, fileSize); err != nil {
		return fmt.Errorf("failed to send file size: %v", err)
	}

	// Send file data
	reader := bufio.NewReader(file)
	written, err := io.Copy(conn, reader)
	if err != nil {
		return fmt.Errorf("failed to send file data: %v", err)
	}

	if written != fileSize {
		return fmt.Errorf("incomplete transfer: sent %d of %d bytes", written, fileSize)
	}

	// Read acknowledgment
	ack := make([]byte, 2)
	_, err = conn.Read(ack)
	if err != nil {
		return fmt.Errorf("failed to read acknowledgment: %v", err)
	}

	if string(ack) != "OK" {
		return fmt.Errorf("server did not acknowledge transfer")
	}

	return nil
}

// ReceiveFile receives a file from a connection (for client-to-client transfers)
func (c *Client) ReceiveFile(conn net.Conn, savePath string) error {
	defer conn.Close()

	// Read filename length
	var filenameLen uint32
	if err := binary.Read(conn, binary.BigEndian, &filenameLen); err != nil {
		return fmt.Errorf("failed to read filename length: %v", err)
	}

	// Read filename
	filenameBytes := make([]byte, filenameLen)
	if _, err := io.ReadFull(conn, filenameBytes); err != nil {
		return fmt.Errorf("failed to read filename: %v", err)
	}

	// Read file size
	var fileSize int64
	if err := binary.Read(conn, binary.BigEndian, &fileSize); err != nil {
		return fmt.Errorf("failed to read file size: %v", err)
	}

	// Create file
	file, err := os.Create(savePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	// Receive file data
	_, err = io.CopyN(file, conn, fileSize)
	if err != nil {
		return fmt.Errorf("failed to receive file: %v", err)
	}

	// Send acknowledgment
	conn.Write([]byte("OK"))

	return nil
}
