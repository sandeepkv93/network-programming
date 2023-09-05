package udp

import (
	"fmt"
	"net"
)

type UDPClient struct {
	serverAddr string
}

func NewUDPClient(serverAddr string) *UDPClient {
	return &UDPClient{
		serverAddr: serverAddr,
	}
}

func (c *UDPClient) SendMessage(msg string) {
	serverUDPAddr, err := net.ResolveUDPAddr("udp", c.serverAddr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, serverUDPAddr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer conn.Close()

	_, err = conn.Write([]byte(msg))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	buf := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	fmt.Println("Received from server:", string(buf[:n]))
}

// func main() {
// 	client := UDPClient{
// 		serverAddr: "127.0.0.1:8080",
// 	}

// 	for {
// 		client.SendMessage("Hello, server!")
// 		time.Sleep(2 * time.Second)
// 	}
// }
