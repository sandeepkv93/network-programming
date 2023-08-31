package ping

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"time"
)

const (
	icmpEchoReply   = 0
	icmpEchoRequest = 8
	defaultTimeout  = 3 * time.Second
)

type icmpMessage struct {
	Type        uint8  // message type
	Code        uint8  // type sub-code
	Checksum    uint16 // checksum for header and payload
	Identifier  uint16 // unique identifier
	SequenceNum uint16 // sequence number
}

// Ping sends an ICMP request to the specified host and waits for a reply.
// It returns the round-trip time duration and any error encountered.
func Ping(host string) (time.Duration, error) {
	// Resolve the IP address of the host
	ipAddr, err := net.ResolveIPAddr("ip", host)
	if err != nil {
		return 0, err
	}

	// Print the ping message
	fmt.Println("Pinging", ipAddr.String(), "with 32 bytes of data:")

	// Open a connection to the host
	conn, err := net.DialIP("ip4:icmp", nil, ipAddr)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	// Create the ICMP request message
	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, icmpMessage{
		Type: icmpEchoRequest,
	})

	// Compute the checksum for the message
	data := buffer.Bytes()
	checksum := computeChecksum(data)
	data[2] = byte(checksum >> 8)
	data[3] = byte(checksum & 0xff)

	// Buffer to store the reply message
	reply := make([]byte, 1024)

	// Channel to receive the ping result
	resultChan := make(chan time.Duration)

	// Goroutine to send ICMP request and wait for a reply
	go func() {
		startTime := time.Now()
		_, err = conn.Write(data)
		if err != nil {
			resultChan <- 0
			return
		}

		conn.SetReadDeadline(time.Now().Add(defaultTimeout))
		_, err = conn.Read(reply)
		if err != nil {
			resultChan <- 0
			return
		}

		resultChan <- time.Since(startTime)
	}()

	// Wait for a reply or timeout
	select {
	case duration := <-resultChan:
		if duration == 0 {
			return 0, errors.New("no response from host")
		}
		return duration, nil
	case <-time.After(defaultTimeout):
		return 0, errors.New("request timed out")
	}
}

// computeChecksum computes the checksum for the given data.
func computeChecksum(data []byte) uint16 {
	if len(data)%2 != 0 {
		data = append(data, 0)
	}

	var sum uint32
	for i := 0; i < len(data); i += 2 {
		sum += uint32(data[i])<<8 | uint32(data[i+1])
	}

	sum = (sum >> 16) + (sum & 0xffff)
	sum += sum >> 16

	return uint16(^sum)
}
