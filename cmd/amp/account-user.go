package main

import (
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/appcelerator/amp/api/client"
	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

var (
	switchRoleCmd = &cobra.Command{
		Use:   "switch ORGANIZATION-NAME",
		Short: "Switch primary organization",
		Long:  `The switch command changes the current login from a user account to the specified organization.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return switchRole(AMP, cmd, args)
		},
	}

	infoCmd = &cobra.Command{
		Use:   "info ACCOUNT-NAME",
		Short: "Display account information",
		Long: `The info command displays information about the specified account name.
	If the input account name belongs to the user who is currently logged-in, the following information is displayed :
	Account Name, Email, Organization Name, Team Name, Billing Information, Other Settings
	If the input account name belongs to a different user, only the following information can be viewed :
	Account Name, Email, Organization Name, Team Name.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return getAccount(AMP, cmd, args)
		},
	}

	editCmd = &cobra.Command{
		Use:   "edit ACCOUNT-NAME",
		Short: "Edit account information",
		Long:  `The edit command allows an account owner to modify the specified account information.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return editAccount(AMP, cmd, args)
		},
	}

	deleteCmd = &cobra.Command{
		Use:   "delete ACCOUNT-NAME",
		Short: "Delete an account",
		Long:  `The delete command allows an account owner to delete the specified account.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return deleteAccount(AMP, cmd, args)
		},
	}

	//TODO
	listUsersCmd = &cobra.Command{
		Use:   "list [ORGANIZATION-NAME] [TEAM-NAME]",
		Short: "List users, optionally filter by organization and team",
		Long:  `The list user command displays information about all users currently on the system, which can be filtered by team or organization.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return listUser(AMP, cmd, args)
		},
	}

	userCmd = &cobra.Command{
		Use:   "user",
		Short: "User operations",
		Long:  `The user command manager all user-related operations.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return userOp(AMP, cmd, args)
		},
	}
)

func init() {
	AccountCmd.AddCommand(userCmd)
	userCmd.AddCommand(switchRoleCmd)
	userCmd.AddCommand(infoCmd)
	userCmd.AddCommand(editCmd)
	userCmd.AddCommand(deleteCmd)
	userCmd.AddCommand(listUsersCmd)
}

func userOp(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
	color.Set(color.FgYellow, color.Bold)
	fmt.Println("Choose a command for user operation")
	fmt.Println("Use amp account user -h for help")
	color.Unset()
	return nil
}

// switchRole validates the input command line arguments and switches to an organization
// to the input value by invoking the corresponding rpc/storage method
func switchRole(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
	fmt.Println("This will switch between the primary account and an organization account that a user is part of.")
	org, er := getOrgName(args)
	if er != nil {
		return fmt.Errorf("user error: %v", er)
	}
	request := &account.TeamRequest{
		Organization: org,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	_, err = client.Switch(context.Background(), request)
	if err != nil {
		return fmt.Errorf("server error: %v", err)
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("Switch organization successful - ", org)
	color.Unset()
	return nil
}

// getAccount validates the input command line arguments and displays the account
// information by invoking the corresponding rpc/storage method
func getAccount(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
	fmt.Println("This will display the information of a specific account.")
	name, er := getName(args)
	if er != nil {
		return fmt.Errorf("user error: %v", er)
	}
	request := &account.AccountRequest{
		Name: name,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	response, e := client.GetAccountDetails(context.Background(), request)
	if er != nil {
		return fmt.Errorf("server error: %v", e)
	}
	if response == nil {
		color.Set(color.FgRed, color.Bold)
		fmt.Println("No information available!")
		color.Unset()
		return nil
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("")
	color.Unset()
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	fmt.Fprintln(w, "ACCOUNT NAME\tEMAIL\tIS EMAIL VERIFIED?\tACCOUNT TYPE\t")
	fmt.Fprintln(w, "------------\t-----\t------------------\t------------\t")
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t\n", response.Account.Name, response.Account.Email, response.Account.EmailVerified, response.Account.AccountType)
	w.Flush()
	return nil
}

// editAccount validates the input command line arguments and modifies the account
// information by invoking the corresponding rpc/storage method
func editAccount(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
	fmt.Println("This will modify the information of a specific account.")
	name, er := getName(args)
	if er != nil {
		return fmt.Errorf("user error: %v", er)
	}
	request := &account.EditRequest{
		Name: name,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	response, e := client.EditAccount(context.Background(), request)
	if e != nil {
		return fmt.Errorf("server error: %v", e)
	}
	if response == nil {
		color.Set(color.FgRed, color.Bold)
		fmt.Println("No information available!")
		color.Unset()
		return nil
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("")
	color.Unset()
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	fmt.Fprintln(w, "ACCOUNT NAME\tEMAIL\tIS EMAIL VERIFIED?\tACCOUNT TYPE\t")
	fmt.Fprintln(w, "------------\t-----\t------------------\t------------\t")
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t\n", response.Account.Name, response.Account.Email, response.Account.EmailVerified, response.Account.AccountType)
	w.Flush()
	return nil
}

// deleteAccount validates the input command line arguments and deletes the account
// by invoking the corresponding rpc/storage method
func deleteAccount(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
	fmt.Println("This will delete the specified account.")
	name, er := getName(args)
	if er != nil {
		return fmt.Errorf("user error: %v", er)
	}
	request := &account.AccountRequest{
		Name: name,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	_, err = client.DeleteAccount(context.Background(), request)
	if err != nil {
		return fmt.Errorf("server error: %v", err)
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("Account '", name, "' deleted successfully.")
	color.Unset()
	return nil
}

// listUser validates the input command line arguments and lists the users available on the system
// by invoking the corresponding rpc/storage method
func listUser(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
	fmt.Println("This will list all users available on the system.")
	org, team, e := getTeamOrg(args)
	if e != nil {
		return fmt.Errorf("user error: %v", e)
	}
	request := &account.AccountsRequest{
		Organization: org,
		Team:         team,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	response, er := client.ListAccounts(context.Background(), request)
	if er != nil {
		return fmt.Errorf("server error: %v", er)
	}
	if response == nil || len(response.Accounts) == 0 {
		color.Set(color.FgRed, color.Bold)
		fmt.Println("No information available!")
		color.Unset()
		return nil
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("")
	color.Unset()
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	fmt.Fprintln(w, "NAME\tEMAIL\tIS EMAIL VERIFIED?\tACCOUNT TYPE\tORGANIZATION\tTEAM\t")
	fmt.Fprintln(w, "----\t-----\t------------------\t------------\t------------\t----\t")
	for _, info := range response.Accounts {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t\n", info.Name, info.Email, info.EmailVerified, info.AccountType, org, team)
	}
	w.Flush()
	return nil
}

func getName(args []string) (name string, err error) {
	switch len(args) {
	case 0:
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify account name")
		color.Unset()
		name = getUserName()
	case 1:
		name = args[0]
	default:
		defer color.Set(color.FgRed, color.Bold)
		return "", errors.New("too many arguments - check again")
	}
	return
}
