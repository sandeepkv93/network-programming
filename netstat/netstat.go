package netstat

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

// ConnectionState represents the state of a connection
type ConnectionState string

const (
	StateEstablished ConnectionState = "ESTABLISHED"
	StateSynSent     ConnectionState = "SYN_SENT"
	StateSynRecv     ConnectionState = "SYN_RECV"
	StateFinWait1    ConnectionState = "FIN_WAIT1"
	StateFinWait2    ConnectionState = "FIN_WAIT2"
	StateTimeWait    ConnectionState = "TIME_WAIT"
	StateClose       ConnectionState = "CLOSE"
	StateCloseWait   ConnectionState = "CLOSE_WAIT"
	StateLastAck     ConnectionState = "LAST_ACK"
	StateListen      ConnectionState = "LISTEN"
	StateClosing     ConnectionState = "CLOSING"
)

// Connection represents a network connection
type Connection struct {
	Protocol    string
	LocalAddr   string
	LocalPort   uint16
	RemoteAddr  string
	RemotePort  uint16
	State       ConnectionState
	UID         uint32
	Inode       uint64
}

var stateMap = map[string]ConnectionState{
	"01": StateEstablished,
	"02": StateSynSent,
	"03": StateSynRecv,
	"04": StateFinWait1,
	"05": StateFinWait2,
	"06": StateTimeWait,
	"07": StateClose,
	"08": StateCloseWait,
	"09": StateLastAck,
	"0A": StateListen,
	"0B": StateClosing,
}

// GetConnections retrieves all network connections
func GetConnections() ([]Connection, error) {
	var connections []Connection

	// Get TCP connections
	tcpConns, err := getTCPConnections()
	if err == nil {
		connections = append(connections, tcpConns...)
	}

	// Get UDP connections
	udpConns, err := getUDPConnections()
	if err == nil {
		connections = append(connections, udpConns...)
	}

	return connections, nil
}

// GetTCPConnections retrieves only TCP connections
func GetTCPConnections() ([]Connection, error) {
	return getTCPConnections()
}

// GetUDPConnections retrieves only UDP connections
func GetUDPConnections() ([]Connection, error) {
	return getUDPConnections()
}

func getTCPConnections() ([]Connection, error) {
	return parseConnections("/proc/net/tcp", "tcp")
}

func getUDPConnections() ([]Connection, error) {
	return parseConnections("/proc/net/udp", "udp")
}

func parseConnections(filename, protocol string) ([]Connection, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %v", filename, err)
	}
	defer file.Close()

	return parseConnectionsFromReader(file, protocol)
}

func parseConnectionsFromReader(r io.Reader, protocol string) ([]Connection, error) {
	var connections []Connection
	scanner := bufio.NewScanner(r)

	// Skip header line
	scanner.Scan()

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		if len(fields) < 10 {
			continue
		}

		conn := Connection{
			Protocol: protocol,
		}

		// Parse local address and port
		localAddrPort := strings.Split(fields[1], ":")
		if len(localAddrPort) == 2 {
			conn.LocalAddr = parseIPAddress(localAddrPort[0])
			conn.LocalPort = parsePort(localAddrPort[1])
		}

		// Parse remote address and port
		remoteAddrPort := strings.Split(fields[2], ":")
		if len(remoteAddrPort) == 2 {
			conn.RemoteAddr = parseIPAddress(remoteAddrPort[0])
			conn.RemotePort = parsePort(remoteAddrPort[1])
		}

		// Parse state
		if state, ok := stateMap[fields[3]]; ok {
			conn.State = state
		}

		// Parse UID
		if uid, err := strconv.ParseUint(fields[7], 10, 32); err == nil {
			conn.UID = uint32(uid)
		}

		// Parse inode
		if inode, err := strconv.ParseUint(fields[9], 10, 64); err == nil {
			conn.Inode = inode
		}

		connections = append(connections, conn)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return connections, nil
}

func parseIPAddress(hexAddr string) string {
	// Parse hex IP address (little-endian)
	if len(hexAddr) != 8 {
		return hexAddr
	}

	var ip [4]byte
	for i := 0; i < 4; i++ {
		val, _ := strconv.ParseUint(hexAddr[i*2:i*2+2], 16, 8)
		ip[3-i] = byte(val)
	}

	return net.IPv4(ip[0], ip[1], ip[2], ip[3]).String()
}

func parsePort(hexPort string) uint16 {
	port, _ := strconv.ParseUint(hexPort, 16, 16)
	return uint16(port)
}

// FormatConnection formats a connection for display
func FormatConnection(conn Connection) string {
	localAddr := fmt.Sprintf("%s:%d", conn.LocalAddr, conn.LocalPort)
	remoteAddr := fmt.Sprintf("%s:%d", conn.RemoteAddr, conn.RemotePort)

	return fmt.Sprintf("%-6s %-23s %-23s %-15s",
		conn.Protocol,
		localAddr,
		remoteAddr,
		conn.State)
}

// FormatConnections formats all connections for display
func FormatConnections(connections []Connection) string {
	output := fmt.Sprintf("%-6s %-23s %-23s %-15s\n",
		"Proto", "Local Address", "Foreign Address", "State")

	for _, conn := range connections {
		output += FormatConnection(conn) + "\n"
	}

	return output
}

// GetListeningPorts returns all listening ports
func GetListeningPorts() ([]Connection, error) {
	connections, err := GetConnections()
	if err != nil {
		return nil, err
	}

	var listening []Connection
	for _, conn := range connections {
		if conn.State == StateListen || conn.Protocol == "udp" {
			listening = append(listening, conn)
		}
	}

	return listening, nil
}
