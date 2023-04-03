package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var firewallCmd = &cobra.Command{
	Use:   "firewall",
	Short: "Manage firewall settings",
}

func initFirewallCmd() *cobra.Command {
	firewallEnsureCmd.Flags().String("name", "", "Name of the firewall rule (required)")
	firewallEnsureCmd.Flags().String("rules-in", "", "Array of input firewall rules in format protocol;port;host;group-ids. \n"+
		"Protocol value: any, icmp, tcp or udp\n"+
		"Port value: any, number (example: 80) or range (example: 16000-16999)\n"+
		"Host value: any or group (in case that is used group than list of group IDs, group names or group ObjectId expected),\n"+
		"	for IDs use format id=###, for name use format name=###, for objectId use format objectId=###\n"+
		"Example:\n"+
		"	any;any;any,tcp;22;group;id=demo.shieldoo.net:groups:1,udp;53;group;objectId=e7549a43-f3c2-4d0d-9cd1-6811a107cdc4;objectId=a3e4ead5-ffb7-4d94-ba71-0185b5466426")
	firewallEnsureCmd.Flags().String("rules-out", "", "Array of output firewall rules in format protocol;port;host;group-ids. \n"+
		"!IMPORTANT! if no value is provided, default output rule any;any;any is created (which is usually fine)\n"+
		"Protocol value: any, icmp, tcp or udp\n"+
		"Port value: any, number (example: 80) or range (example: 16000-16999)\n"+
		"Host value: any or group (in case that is used group than list of group IDs, group names or group ObjectId expected)\n"+
		"	for IDs use format id=###, for name use format name=###, for objectId use format objectId=###\n"+
		"Example:\n"+
		"	any;any;any,tcp;22;group;demo.shieldoo.net:groups:1,udp;53;group:e7549a43-f3c2-4d0d-9cd1-6811a107cdc4;a3e4ead5-ffb7-4d94-ba71-0185b5466426")
	firewallEnsureCmd.MarkFlagRequired("name")
	firewallCmd.AddCommand(firewallEnsureCmd)

	firewallDeleteCmd.Flags().String("id", "", "ID of the firewall rule to delete (required)")
	firewallDeleteCmd.MarkFlagRequired("id")
	firewallCmd.AddCommand(firewallDeleteCmd)

	firewallCmd.AddCommand(firewallListCmd)

	firewallShowCmd.Flags().String("id", "", "ID of the firewall rule to show (required)")
	firewallShowCmd.Flags().String("name", "", "Name of the firewall rule to show (required)")
	firewallCmd.AddCommand(firewallShowCmd)

	return firewallCmd
}

var firewallEnsureCmd = &cobra.Command{
	Use:   "ensure",
	Short: "Ensure a firewall (create or update)",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		rulesIn, _ := cmd.Flags().GetString("rules-in")
		rulesOut, _ := cmd.Flags().GetString("rules-out")

		// parse rules
		rin, err := parseFirewallRules(rulesIn)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			os.Exit(1)
		}
		rout, err := parseFirewallRules(rulesOut)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			os.Exit(1)
		}
		// if out rules empty, create default
		if len(rout) == 0 {
			rout = append(rout, FirewallRule{
				Protocol: "any",
				Port:     "any",
				Host:     "any",
			})
		}
		fw := Firewall{
			Name:     name,
			RulesIn:  rin,
			RulesOut: rout,
		}
		// convert FW name to ID
		fwDetailData, err := callApi("GET", "/firewalls", name, "", nil)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			os.Exit(1)
		}
		var fwDetail []Firewall
		err = json.Unmarshal([]byte(fwDetailData), &fwDetail)
		if err != nil {
			fmt.Printf("ERROR: %s (%s)\n", err, fwDetailData)
			os.Exit(1)
		}
		if len(fwDetail) != 0 {
			// update
			fwDetailData, err = callApi("PUT", "firewalls", "", fwDetail[0].Id, &fw)
			if err != nil {
				fmt.Printf("ERROR: %s, %s\n", err, fwDetailData)
				os.Exit(1)
			}
		} else {
			// create
			fwDetailData, err = callApi("POST", "firewalls", "", "", &fw)
			if err != nil {
				fmt.Printf("ERROR: %s, %s\n", err, fwDetailData)
				os.Exit(1)
			}
		}
		fmt.Printf("Firewall detail: %s\n", fwDetailData)
	},
}

var firewallDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a firewall",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("id")
		ret, err := callApi("DELETE", "firewalls", "", id, nil)
		if err != nil {
			fmt.Printf("ERROR: %s - %s\n", err.Error(), ret)
			os.Exit(1)
		}
		fmt.Println("Firewall deleted")
	},
}

var firewallListCmd = &cobra.Command{
	Use:   "list",
	Short: "List firewalls",
	Run: func(cmd *cobra.Command, args []string) {
		ret, err := callApi("GET", "firewalls", "", "", nil)
		if err != nil {
			fmt.Printf("ERROR: ", err.Error())
			os.Exit(1)
		}
		fmt.Println(ret)
	},
}

var firewallShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show a firewall",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		id, _ := cmd.Flags().GetString("id")
		if name == "" && id == "" {
			fmt.Printf("Error: either name or id must be specified\n")
			os.Exit(1)
		}
		ret, err := callApi("GET", "firewalls", name, id, nil)
		if err != nil {
			fmt.Printf("ERROR: %s - %s\n", err.Error(), ret)
			os.Exit(1)
		}
		if ret == "" || strings.TrimSpace(ret) == "null" {
			fmt.Printf("Firewall not found\n")
			os.Exit(1)
		}
		// if it is json array than clean it up
		if strings.HasPrefix(ret, "[") {
			ret = strings.TrimPrefix(ret, "[")
			ret = strings.TrimSuffix(ret, "]")
		}
		fmt.Println(ret)
	},
}
