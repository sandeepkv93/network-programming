package smtp

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

// Client represents an SMTP client
type Client struct {
	address string
	conn    net.Conn
	reader  *bufio.Reader
}

// NewClient creates a new SMTP client
func NewClient(address string) *Client {
	return &Client{
		address: address,
	}
}

// Connect connects to the SMTP server
func (c *Client) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, 10*time.Second)
	if err != nil {
		return fmt.Errorf("failed to connect: %v", err)
	}
	c.conn = conn
	c.reader = bufio.NewReader(conn)

	// Read greeting
	response, err := c.readResponse()
	if err != nil {
		return fmt.Errorf("failed to read greeting: %v", err)
	}

	if !strings.HasPrefix(response, "220") {
		return fmt.Errorf("unexpected greeting: %s", response)
	}

	return nil
}

// Hello sends HELO command
func (c *Client) Hello(hostname string) error {
	if err := c.sendCommand(fmt.Sprintf("HELO %s", hostname)); err != nil {
		return err
	}

	response, err := c.readResponse()
	if err != nil {
		return err
	}

	if !strings.HasPrefix(response, "250") {
		return fmt.Errorf("HELO failed: %s", response)
	}

	return nil
}

// MailFrom sends MAIL FROM command
func (c *Client) MailFrom(from string) error {
	if err := c.sendCommand(fmt.Sprintf("MAIL FROM:<%s>", from)); err != nil {
		return err
	}

	response, err := c.readResponse()
	if err != nil {
		return err
	}

	if !strings.HasPrefix(response, "250") {
		return fmt.Errorf("MAIL FROM failed: %s", response)
	}

	return nil
}

// RcptTo sends RCPT TO command
func (c *Client) RcptTo(to string) error {
	if err := c.sendCommand(fmt.Sprintf("RCPT TO:<%s>", to)); err != nil {
		return err
	}

	response, err := c.readResponse()
	if err != nil {
		return err
	}

	if !strings.HasPrefix(response, "250") {
		return fmt.Errorf("RCPT TO failed: %s", response)
	}

	return nil
}

// Data sends the DATA command and email content
func (c *Client) Data(message string) error {
	if err := c.sendCommand("DATA"); err != nil {
		return err
	}

	response, err := c.readResponse()
	if err != nil {
		return err
	}

	if !strings.HasPrefix(response, "354") {
		return fmt.Errorf("DATA command failed: %s", response)
	}

	// Send message
	if err := c.sendCommand(message); err != nil {
		return err
	}

	// Send terminator
	if err := c.sendCommand("."); err != nil {
		return err
	}

	response, err = c.readResponse()
	if err != nil {
		return err
	}

	if !strings.HasPrefix(response, "250") {
		return fmt.Errorf("message not accepted: %s", response)
	}

	return nil
}

// SendMail sends a complete email
func (c *Client) SendMail(from string, to []string, subject, body string) error {
	// HELO
	if err := c.Hello("client"); err != nil {
		return err
	}

	// MAIL FROM
	if err := c.MailFrom(from); err != nil {
		return err
	}

	// RCPT TO (for each recipient)
	for _, recipient := range to {
		if err := c.RcptTo(recipient); err != nil {
			return err
		}
	}

	// Build message
	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		from,
		strings.Join(to, ", "),
		subject,
		body)

	// Send DATA
	if err := c.Data(message); err != nil {
		return err
	}

	return nil
}

// Quit sends QUIT command and closes connection
func (c *Client) Quit() error {
	if c.conn == nil {
		return nil
	}

	c.sendCommand("QUIT")
	c.readResponse()

	return c.conn.Close()
}

// sendCommand sends a command to the server
func (c *Client) sendCommand(command string) error {
	if c.conn == nil {
		return fmt.Errorf("not connected")
	}

	_, err := c.conn.Write([]byte(command + "\r\n"))
	return err
}

// readResponse reads a response from the server
func (c *Client) readResponse() (string, error) {
	line, err := c.reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(line), nil
}
