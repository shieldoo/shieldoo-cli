package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var shieldooUri = ""
var shieldooApiKey = ""

func init() {
	rootCmd.AddCommand(initServerCmd())
	rootCmd.AddCommand(initFirewallCmd())
	rootCmd.AddCommand(initGroupCmd())
	// load env variables
	shieldooUri = os.Getenv("SHIELDOO_URI")
	if shieldooUri == "" {
		fmt.Println("Error: SHIELDOO_URI environment variable not set")
		os.Exit(1)
	}
	shieldooApiKey = os.Getenv("SHIELDOO_APIKEY")
	if shieldooApiKey == "" {
		fmt.Println("Error: SHIELDOO_APIKEY environment variable not set")
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "shieldoo",
	Short: "A simple CLI tool",
	Long:  "A simple CLI tool to manage shieldoo servers and firewalls.",
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

//tokens
// https://blog.canopas.com/jwt-in-golang-how-to-implement-token-based-authentication-298c89a26ffd
// asdfghjklpoiuztrewq
