package main

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/valcinei/jiboia-tunnel/server"
)

func main() {
	var addr string
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Mock local HTTP server (para testes)",
		Run: func(cmd *cobra.Command, args []string) {
			mockServer := server.NewServer(addr)
			if err := mockServer.Start(); err != nil {
				log.Fatal(err)
			}
		},
	}
	cmd.Flags().StringVar(&addr, "addr", ":3000", "Endereço para escutar")

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
