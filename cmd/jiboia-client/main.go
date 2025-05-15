package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/valcinei/jiboia-tunnel/client"
)

func main() {
	var local, relay, name string
	cmd := &cobra.Command{
		Use:   "client [protocol] [port]",
		Short: "Inicia o cliente que conecta ao relay",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			protocol := args[0]
			port := args[1]
			local = fmt.Sprintf("%s://localhost:%s", protocol, port)

			if name == "" {
				name = client.GenerateRandomName()
				fmt.Println("Nome gerado:", name)
			}
			c := client.NewClient(local, relay, name)
			if err := c.Start(); err != nil {
				log.Fatal(err)
			}
		},
	}
	cmd.Flags().StringVar(&relay, "relay", "ws://localhost:80/ws", "URL do relay")
	cmd.Flags().StringVar(&name, "name", "", "Nome do túnel (subdomínio)")

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
