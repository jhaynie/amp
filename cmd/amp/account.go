package main

import (
	"github.com/spf13/cobra"
)

//
var	AccountCmd = &cobra.Command{
	Use:   "account",
	Short: "Account operations",
	Long:  `The account command manages all account-related operations.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return AMP.Connect()
	},
}

func init() {
	RootCmd.AddCommand(AccountCmd)
}
