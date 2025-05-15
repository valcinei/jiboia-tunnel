package main

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/valcinei/jiboia-tunnel/relay"
)

func main() {
	var addr string
	cmd := &cobra.Command{
		Use:   "relay",
		Short: "Inicia o relay WebSocket e proxy HTTP",
		Run: func(cmd *cobra.Command, args []string) {
			relayServer := relay.NewServer()
			if err := relayServer.Start(addr); err != nil {
				log.Fatal(err)
			}
		},
	}
	cmd.Flags().StringVar(&addr, "addr", ":80", "Endereço para escutar")

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
