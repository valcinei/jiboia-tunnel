package main

import (
	"os"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "jiboia",
		Short: "Jiboia Tunnel CLI",
		Long:  `Jiboia é uma ferramenta de tunelamento HTTP reverso baseada em WebSocket.`,
	}

	rootCmd.AddCommand(cmdRelay())
	rootCmd.AddCommand(cmdClient())
	rootCmd.AddCommand(cmdServer())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
} 