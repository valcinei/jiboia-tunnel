package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/valcinei/jiboia-tunnel/client"
	"github.com/valcinei/jiboia-tunnel/shared"
)

var DefaultRelayURL = func() string {
	if val := os.Getenv("DEFAULT_RELAY_URL"); val != "" {
		return val
	}
	return "wss://relay.jiboia.cloud/ws"
}()

func main() {
	var local, relay, name, proto, hostname, authtoken, config, region, label, logLevel, baseDomain string
	var inspect bool

	// Initialize the in-memory store
	store := shared.NewInMemoryStore()

	cmd := &cobra.Command{
		Use:   "client [protocol] [port]",
		Short: "Starts the client that connects to the relay",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			protocol := args[0]
			port := args[1]
			local = fmt.Sprintf("%s://localhost:%s", protocol, port)

			if name == "" {
				name = client.GenerateRandomName()
				fmt.Println("Generated name:", name)
			}

			// Add tunnel information to the in-memory store
			store.AddTunnel(name, local, relay)

			if baseDomain == "" {
				baseDomain = os.Getenv("TUNNEL_DOMAIN")
				if baseDomain == "" {
					baseDomain = "jiboia.cloud"
				}
			}

			c := client.NewClient(local, relay, name, baseDomain)
			if err := c.Start(); err != nil {
				log.Fatal(err)
			}
		},
	}

	cmd.Flags().StringVar(&relay, "relay", DefaultRelayURL, "Relay URL")
	cmd.Flags().StringVar(&name, "name", "", "Tunnel name (subdomain)")
	cmd.Flags().StringVar(&proto, "proto", "http", "Protocol to expose (http, tcp)")
	cmd.Flags().StringVar(&hostname, "hostname", "", "Full custom domain (e.g., mywebsite.com)")
	cmd.Flags().BoolVar(&inspect, "inspect", false, "Shows detailed traffic (debug mode)")
	cmd.Flags().StringVar(&authtoken, "authtoken", "", "Authentication token with the server")
	cmd.Flags().StringVar(&config, "config", "", "Path to external configuration file")
	cmd.Flags().StringVar(&region, "region", "", "Relay region (e.g., us, sa-east)")
	cmd.Flags().StringVar(&label, "label", "", "Friendly tunnel identifier (used in logs/future API)")
	cmd.Flags().StringVar(&logLevel, "log-level", "info", "Log level (debug, info, warn, error)")
	cmd.Flags().StringVar(&baseDomain, "domain", "", "Base domain to generate tunnel URL (e.g., jiboia.cloud)")

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
