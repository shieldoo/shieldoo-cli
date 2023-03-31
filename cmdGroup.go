package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var groupCmd = &cobra.Command{
	Use:   "group",
	Short: "Manage groups",
}

func initGroupCmd() *cobra.Command {
	groupCmd.AddCommand(groupListCmd)

	groupShowCmd.Flags().String("name", "", "Name of the group to show")
	groupShowCmd.Flags().String("id", "", "Id of the group to show")
	groupCmd.AddCommand(groupShowCmd)

	return groupCmd
}

var groupListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all groups",
	Run: func(cmd *cobra.Command, args []string) {
		ret, err := callApi("GET", "groups", "", "", nil)
		if err != nil {
			fmt.Printf("ERROR: ", err.Error())
			os.Exit(1)
		}
		fmt.Println(ret)
	},
}

var groupShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show a group",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		id, _ := cmd.Flags().GetString("id")
		if name == "" && id == "" {
			fmt.Printf("Error: either name or id must be specified\n")
			os.Exit(1)
		}
		ret, err := callApi("GET", "groups", name, id, nil)
		if err != nil {
			fmt.Printf("ERROR: %s - %s\n", err.Error(), ret)
			os.Exit(1)
		}
		if ret == "" || strings.TrimSpace(ret) == "null" {
			fmt.Printf("Group not found\n")
			os.Exit(1)
		}
		fmt.Println(ret)
	},
}
