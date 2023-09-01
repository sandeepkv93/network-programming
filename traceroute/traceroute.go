package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"golang.org/x/net/ipv4"
)

const (
	maxHops         = 64
	timeoutDuration = 2 * time.Second
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go-traceroute <hostname>")
		os.Exit(1)
	}

	host := os.Args[1]
	traceroute(host)
}

func traceroute(host string) {
	ipAddr, err := net.ResolveIPAddr("ip", host)
	if err != nil {
		fmt.Println("Failed to resolve hostname:", err)
		return
	}

	fmt.Printf("Tracing route to %s (%s)\n", host, ipAddr.String())

	var wg sync.WaitGroup

	for i := 1; i <= maxHops; i++ {
		wg.Add(1)
		go traceHop(host, i, &wg)
	}

	wg.Wait()
}

func traceHop(host string, ttl int, wg *sync.WaitGroup) {
	defer wg.Done()
	ipAddr, err := net.ResolveIPAddr("ip", host)
	if err != nil {
		fmt.Println("Failed to resolve hostname:", err)
		return
	}

	// Create a raw socket to send and receive ICMP packets
	conn, err := net.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		log.Printf("Failed to create socket: %s\n", err)
		return
	}
	defer conn.Close()

	// Set the TTL on the socket
	if err := conn.SetDeadline(time.Now().Add(timeoutDuration)); err != nil {
		log.Printf("Failed to set deadline: %s\n", err)
		return
	}
	if err := ipv4.NewPacketConn(conn).SetTTL(ttl); err != nil {
		log.Printf("Failed to set TTL: %s\n", err)
		return
	}

	// Create an ICMP packet
	packet := make([]byte, 8)
	packet[0] = 8 // ICMP type Echo Request

	// Send the packet
	if _, err := conn.WriteTo(packet, ipAddr.IP); err != nil {
		log.Printf("Failed to send packet: %s\n", err)
		return
	}

	// Receive the response packet
	recvBuf := make([]byte, 1024)
	n, addr, err := conn.ReadFrom(recvBuf)
	if err != nil {
		log.Printf("Failed to read response packet: %s\n", err)
		return
	}

	peer := addr.String()

	if recvBuf[0] == 11 { // Time Exceeded
		fmt.Printf("%d %s\n", ttl, peer)
	} else if recvBuf[0] == 0 { // Echo Reply
		fmt.Printf("%d %s\n", ttl, peer)
		if n >= 8 {
			fmt.Printf("Reached %s\n", host)
		}
	}
}
