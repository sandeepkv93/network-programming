package cmd

import (
	"fmt"
	"os"
	"os/user"

	"networkprogramming/ping"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(pingCmd)
}

var pingCmd = &cobra.Command{
	Use:   "ping [hostname]",
	Short: "Ping a specified host",
	Long:  `Ping is a simple tool to ping a specified host using ICMP.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if u, err := user.Current(); err != nil || u.Uid != "0" {
			fmt.Println("Error: This command must be run with sudo")
			os.Exit(1)
		}
		host := args[0]
		duration, err := ping.Ping(host)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Printf("Ping to %s took %v\n", host, duration)
	},
}
