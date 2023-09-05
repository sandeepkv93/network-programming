package cmd

import (
	"fmt"

	"networkprogramming/udp"

	"github.com/spf13/cobra"
)

var serverType string
var port string

// udpCmd represents the udp command
var udpCmd = &cobra.Command{
	Use:   "udp",
	Short: "UDP client/server",
	Run: func(cmd *cobra.Command, args []string) {
		switch serverType {
		case "server":
			server := udp.NewUDPServer(":" + port)
			server.Start()
		case "client":
			client := udp.NewUDPClient("127.0.0.1:" + port)
			client.SendMessage("Hello, server!")
		default:
			fmt.Println("Invalid type. Choose either 'server' or 'client'.")
		}
	},
}

func init() {
	rootCmd.AddCommand(udpCmd)

	udpCmd.Flags().StringVar(&serverType, "type", "", "Type of the UDP instance to run (server or client)")
	// Add an optional flag port to the UDP server
	udpCmd.Flags().StringVar(&port, "port", "8080", "Port to listen on")
}
