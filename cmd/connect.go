package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to the sample digital service",
	RunE:  runClient,
}

func init() {
	rootCmd.AddCommand(connectCmd)
}

func runClient(cmd *cobra.Command, args []string) error {
	fmt.Println("connecting now")
	return nil
}
