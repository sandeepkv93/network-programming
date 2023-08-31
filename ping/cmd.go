package ping

import (
	"fmt"

	"github.com/spf13/cobra"
)

var pingCommand = &cobra.Command{
	Use:   "ping [hostname]",
	Short: "Ping a specified host",
	Long:  `Ping is a simple tool to ping a specified host using ICMP.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		host := args[0]
		duration, err := Ping(host)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Printf("Ping to %s took %v\n", host, duration)
	},
}

func Execute() {
	if err := pingCommand.Execute(); err != nil {
		fmt.Println(err)
	}
}
