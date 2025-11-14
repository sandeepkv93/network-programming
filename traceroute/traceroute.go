package traceroute

import (
	"fmt"
	"net"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

// Hop represents a single hop in the route
type Hop struct {
	Number  int
	Address string
	RTT     time.Duration
	Timeout bool
}

// Result represents the traceroute result
type Result struct {
	Destination string
	Hops        []Hop
}

// Traceroute performs a traceroute to the specified destination
func Traceroute(destination string, maxHops int, timeout time.Duration) (*Result, error) {
	// Resolve destination
	destAddr, err := net.ResolveIPAddr("ip4", destination)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve destination: %v", err)
	}

	result := &Result{
		Destination: destAddr.String(),
		Hops:        make([]Hop, 0),
	}

	// Create ICMP connection
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return nil, fmt.Errorf("failed to create ICMP connection: %v", err)
	}
	defer conn.Close()

	// Perform traceroute
	for ttl := 1; ttl <= maxHops; ttl++ {
		hop := Hop{
			Number: ttl,
		}

		// Set TTL
		if err := conn.IPv4PacketConn().SetTTL(ttl); err != nil {
			return nil, fmt.Errorf("failed to set TTL: %v", err)
		}

		// Send ICMP Echo Request
		msg := icmp.Message{
			Type: ipv4.ICMPTypeEcho,
			Code: 0,
			Body: &icmp.Echo{
				ID:   1,
				Seq:  ttl,
				Data: []byte("TRACEROUTE"),
			},
		}

		msgBytes, err := msg.Marshal(nil)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal ICMP message: %v", err)
		}

		start := time.Now()
		_, err = conn.WriteTo(msgBytes, destAddr)
		if err != nil {
			return nil, fmt.Errorf("failed to send ICMP message: %v", err)
		}

		// Set read timeout
		conn.SetReadDeadline(time.Now().Add(timeout))

		// Read response
		reply := make([]byte, 1500)
		n, peer, err := conn.ReadFrom(reply)
		rtt := time.Since(start)

		if err != nil {
			// Timeout
			hop.Timeout = true
			result.Hops = append(result.Hops, hop)
			continue
		}

		hop.RTT = rtt
		hop.Address = peer.String()

		// Parse ICMP message
		rm, err := icmp.ParseMessage(1, reply[:n])
		if err != nil {
			result.Hops = append(result.Hops, hop)
			continue
		}

		result.Hops = append(result.Hops, hop)

		// Check if we reached destination
		if rm.Type == ipv4.ICMPTypeEchoReply {
			break
		}
	}

	return result, nil
}

// FormatResult formats the traceroute result for display
func FormatResult(result *Result) string {
	output := fmt.Sprintf("traceroute to %s, %d hops max\n", result.Destination, len(result.Hops))

	for _, hop := range result.Hops {
		if hop.Timeout {
			output += fmt.Sprintf("%2d  * * * Request timeout\n", hop.Number)
		} else {
			output += fmt.Sprintf("%2d  %s  %v\n", hop.Number, hop.Address, hop.RTT)
		}
	}

	return output
}

// TraceWithRetries performs traceroute with multiple attempts per hop
func TraceWithRetries(destination string, maxHops int, timeout time.Duration, retries int) (*Result, error) {
	// Resolve destination
	destAddr, err := net.ResolveIPAddr("ip4", destination)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve destination: %v", err)
	}

	result := &Result{
		Destination: destAddr.String(),
		Hops:        make([]Hop, 0),
	}

	// Create ICMP connection
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return nil, fmt.Errorf("failed to create ICMP connection: %v", err)
	}
	defer conn.Close()

	// Perform traceroute
	for ttl := 1; ttl <= maxHops; ttl++ {
		hop := Hop{
			Number:  ttl,
			Timeout: true,
		}

		// Try multiple times
		for attempt := 0; attempt < retries; attempt++ {
			// Set TTL
			if err := conn.IPv4PacketConn().SetTTL(ttl); err != nil {
				continue
			}

			// Send ICMP Echo Request
			msg := icmp.Message{
				Type: ipv4.ICMPTypeEcho,
				Code: 0,
				Body: &icmp.Echo{
					ID:   1,
					Seq:  ttl*100 + attempt,
					Data: []byte("TRACEROUTE"),
				},
			}

			msgBytes, err := msg.Marshal(nil)
			if err != nil {
				continue
			}

			start := time.Now()
			_, err = conn.WriteTo(msgBytes, destAddr)
			if err != nil {
				continue
			}

			// Set read timeout
			conn.SetReadDeadline(time.Now().Add(timeout))

			// Read response
			reply := make([]byte, 1500)
			n, peer, err := conn.ReadFrom(reply)
			rtt := time.Since(start)

			if err == nil {
				hop.RTT = rtt
				hop.Address = peer.String()
				hop.Timeout = false

				// Parse ICMP message to check if destination reached
				rm, err := icmp.ParseMessage(1, reply[:n])
				if err == nil && rm.Type == ipv4.ICMPTypeEchoReply {
					result.Hops = append(result.Hops, hop)
					return result, nil
				}
				break
			}
		}

		result.Hops = append(result.Hops, hop)
	}

	return result, nil
}
