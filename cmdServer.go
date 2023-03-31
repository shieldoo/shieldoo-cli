package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Manage servers",
}

func initServerCmd() *cobra.Command {
	serverEnsureCmd.Flags().String("name", "", "Name of the server to ensure (required)")
	serverEnsureCmd.Flags().String("groups", "", "Servers groups (comma separated) - list of group IDs, group names or group ObjectId expected)\n"+
		"	for IDs use format id=###, for name use format name=###, for objectId use format objectId=###")
	serverEnsureCmd.Flags().String("ip", "", "IP address of the server (optional)")
	serverEnsureCmd.Flags().String("listeners", "", "Listeners (comma separated) - list of listeners in format ListenerPort;Protocol;ForwardPort;ForwardHost;Description\n"+
		"	ListenerPort - port to listen on\n"+
		"	Protocol - protocol to listen on (tcp or udp)\n"+
		"	ForwardPort - port to forward to\n"+
		"	ForwardHost - host to forward to\n"+
		"	Example: 80;tcp;8080;myhost.example.com,443;tcp;8443;myhost.example.com;Any description")
	serverEnsureCmd.Flags().String("firewall-id", "", "Firewall ID (required)")
	serverEnsureCmd.Flags().String("firewall-name", "", "Firewall name (required)")
	serverEnsureCmd.Flags().String("description", "", "Description of the server (optional)")
	serverEnsureCmd.MarkFlagRequired("name")
	serverCmd.AddCommand(serverEnsureCmd)

	serverDeleteCmd.Flags().String("id", "", "ID of the server to delete (required)")
	serverDeleteCmd.MarkFlagRequired("id")
	serverCmd.AddCommand(serverDeleteCmd)

	serverCmd.AddCommand(serverListCmd)

	serverShowCmd.Flags().String("name", "", "Name of the server to show (required)")
	serverShowCmd.Flags().String("id", "", "Id of the server to show (required)")
	serverCmd.AddCommand(serverShowCmd)

	return serverCmd
}

var serverEnsureCmd = &cobra.Command{
	Use:   "ensure",
	Short: "Ensure a server (create or update)",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		firewallId, _ := cmd.Flags().GetString("firewall-id")
		firewallName, _ := cmd.Flags().GetString("firewall-name")
		listeners, _ := cmd.Flags().GetString("listeners")
		groups, _ := cmd.Flags().GetString("groups")
		ipAddr, _ := cmd.Flags().GetString("ip")
		description, _ := cmd.Flags().GetString("description")

		if firewallId == "" && firewallName == "" {
			fmt.Printf("Error: either firewall-id or firewall-id must be specified\n")
			os.Exit(1)
		}

		// get firewall id if name is given
		if firewallId == "" {
			ret, err := callApi("GET", "firewalls", firewallName, "", nil)
			if err != nil {
				fmt.Printf("ERROR: %s\n", err.Error(), ret)
				os.Exit(1)
			}
			var fws []Firewall
			err = json.Unmarshal([]byte(ret), &fws)
			if err != nil {
				fmt.Printf("ERROR: %s\n", err.Error())
				os.Exit(1)
			}
			if len(fws) == 0 {
				fmt.Printf("ERROR: no firewall found with name '%s'\n", firewallName)
				os.Exit(1)
			}
			firewallId = fws[0].Id
		}

		// parse groups
		groupsList := strings.Split(groups, ",")
		var serverGroups []Group
		for _, group := range groupsList {
			group = strings.TrimSpace(group)
			if group == "" {
				continue
			}
			g, err := parseGroup(group)
			if err != nil {
				fmt.Printf("ERROR: %s\n", err.Error())
				os.Exit(1)
			}
			serverGroups = append(serverGroups, g)
		}

		// convert server name to id
		ret, err := callApi("GET", "servers", name, "", nil)
		if err != nil {
			fmt.Printf("ERROR: %s %s\n", err.Error(), ret)
			os.Exit(1)
		}
		var servers []Server
		err = json.Unmarshal([]byte(ret), &servers)
		if err != nil {
			fmt.Printf("ERROR: %s %s\n", err.Error(), ret)
			os.Exit(1)
		}
		list, err := parseListeners(listeners)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
			os.Exit(1)
		}
		server := Server{
			Name:   name,
			Groups: serverGroups,
			Firewall: Firewall{
				Id: firewallId,
			},
			Listeners:   list,
			IpAddress:   ipAddr,
			Description: description,
		}
		if len(servers) > 0 {
			// server already exists
			server.Id = servers[0].Id
			ret, err = callApi("PUT", "servers", "", server.Id, server)
			if err != nil {
				fmt.Printf("ERROR: %s %s\n", err.Error(), ret)
				os.Exit(1)
			}
		} else {
			// create server
			ret, err = callApi("POST", "servers", "", "", server)
			if err != nil {
				fmt.Printf("ERROR: %s %s\n", err.Error(), ret)
				os.Exit(1)
			}
		}
		fmt.Println(ret)
	},
}

var serverDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a server",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("id")
		ret, err := callApi("DELETE", "servers", "", id, nil)
		if err != nil {
			fmt.Printf("ERROR: %s %s\n", err.Error(), ret)
			os.Exit(1)
		}
		fmt.Println("Server deleted")
	},
}

var serverListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all servers",
	Run: func(cmd *cobra.Command, args []string) {
		ret, err := callApi("GET", "servers", "", "", nil)
		if err != nil {
			fmt.Printf("ERROR: ", err.Error())
			os.Exit(1)
		}
		fmt.Println(ret)
	},
}

var serverShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show a server",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		id, _ := cmd.Flags().GetString("id")
		if name == "" && id == "" {
			fmt.Printf("Error: either name or id must be specified\n")
			os.Exit(1)
		}
		ret, err := callApi("GET", "servers", name, id, nil)
		if err != nil {
			fmt.Printf("ERROR: %s - %s\n", err.Error(), ret)
			os.Exit(1)
		}
		if ret == "" || strings.TrimSpace(ret) == "null" {
			fmt.Printf("Server not found\n")
			os.Exit(1)
		}
		// if it is array, clean it up
		if ret[0] == '[' {
			ret = strings.TrimPrefix(ret, "[")
			ret = strings.TrimSuffix(ret, "]")
		}
		fmt.Println(ret)
	},
}
