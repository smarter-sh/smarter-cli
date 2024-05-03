/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package account

import (
	"fmt"

	"github.com/QueriumCorp/smarter-cli/cmd/delete"

	"github.com/spf13/cobra"
)

// DeleteCmd represents the delete command
var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "Account management commands",
	Long: `Account management commands:

smarter account <subcommand>

Subcommands:
  user: Delete a user from your account`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("delete called")
	},
}

func init() {
	delete.DeleteCmd.AddCommand(accountCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// DeleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// DeleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
