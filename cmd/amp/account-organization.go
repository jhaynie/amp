package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/appcelerator/amp/api/client"
	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

var (
	createOrganizationCmd = &cobra.Command{
		Use:   "create ORGANIZATION-NAME EMAIL",
		Short: "Create organization",
		Long:  `The create organization command creates an organization with a name and email address.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return createOrganization(AMP, cmd, args)
		},
	}

	editOrganizationCmd = &cobra.Command{
		Use:   "edit ORGANIZATION-NAME",
		Short: "Edit organization",
		Long:  `The edit organization command updates an existing organization.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return editOrganization(AMP, cmd, args)
		},
	}

	deleteOrganizationCmd = &cobra.Command{
		Use:   "delete ORGANIZATION-NAME",
		Short: "Delete organization",
		Long:  `The delete organization command deletes an existing organization and all related information.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return deleteOrganization(AMP, cmd, args)
		},
	}

	listOrganizationsCmd = &cobra.Command{
		Use:   "list",
		Short: "List organization",
		Long:  `The list organization command displays all the available organizations.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return listOrganization(AMP, cmd, args)
		},
	}

	addOrganizationMembersCmd = &cobra.Command{
		Use:   "add ORGANIZATION-NAME MEMBERS...",
		Short: "Add users to organization",
		Long:  `The add-organization command allows an owner team member to add new members in an organization.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return addOrgMember(AMP, cmd, args)
		},
	}

	removeOrganizationMembersCmd = &cobra.Command{
		Use:   "remove ORGANIZATION-NAME MEMBERS...",
		Short: "Remove users from organization",
		Long:  `The remove-organization command allows an owner team member to remove existing members from an organization.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return removeOrgMember(AMP, cmd, args)
		},
	}

	transferOwnershipCmd = &cobra.Command{
		Use:   "transfer RESOURCE-ID ORGANIZATION-NAME",
		Short: "Transfer ownership to organization",
		Long:  `The transfer command allows a resource owner to transfer a particular resource to a different organization.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return transferOwnership(AMP, cmd, args)
		},
	}

	orgCmd = &cobra.Command{
		Use:   "organization",
		Short: "Organization operations",
		Long:  `The organization command manages all organization-related operations.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return organization(AMP, cmd, args)
		},
	}
)

func init() {
	AccountCmd.AddCommand(orgCmd)
	orgCmd.AddCommand(createOrganizationCmd)
	orgCmd.AddCommand(editOrganizationCmd)
	orgCmd.AddCommand(deleteOrganizationCmd)
	orgCmd.AddCommand(listOrganizationsCmd)
	orgCmd.AddCommand(addOrganizationMembersCmd)
	orgCmd.AddCommand(removeOrganizationMembersCmd)
	orgCmd.AddCommand(transferOwnershipCmd)
}

func organization(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
	color.Set(color.FgYellow, color.Bold)
	fmt.Println("Choose a command for organization operation")
	fmt.Println("Use amp account organization -h for help")
	color.Unset()
	return nil
}

// createOrganization validates the input command line arguments and creates an organization
// with name and email by invoking the corresponding rpc/storage method
func createOrganization(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
	fmt.Println("This will create an organization account with the specified account name and email address.")
	org, email, er := getOrgNameEmail(args)
	if er != nil {
		return fmt.Errorf("user error: %v", er)
	}
	request := &account.OrganizationRequest{
		Name:  org,
		Email: email,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	_, err = client.CreateOrganization(context.Background(), request)
	if err != nil {
		return fmt.Errorf("server error: %v", err)
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("Organization '", org, "' created successfully.")
	color.Unset()
	return nil
}

// editOrganization validates the input command line arguments and creates an organization
// with name and email by invoking the corresponding rpc/storage method
func editOrganization(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
	fmt.Println("This will edit an organization account.")
	org, email, er := getOrgNameEmail(args)
	if er != nil {
		return fmt.Errorf("user error: %v", er)
	}
	request := &account.OrganizationRequest{
		Name:  org,
		Email: email,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	_, err = client.EditOrganization(context.Background(), request)
	if err != nil {
		return fmt.Errorf("server error: %v", err)
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("Organization '", org, "' edited successfully.")
	color.Unset()
	return nil
}

// deleteOrganization validates the input command line arguments and deletes an organization
// by invoking the corresponding rpc/storage method
func deleteOrganization(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
	fmt.Println("This will delete an organization account.")
	org, er := getOrgName(args)
	if er != nil {
		return fmt.Errorf("user error: %v", er)
	}
	request := &account.OrganizationRequest{
		Name: org,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	_, err = client.DeleteOrganization(context.Background(), request)
	if err != nil {
		return fmt.Errorf("server error: %v", err)
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("Organization '", org, "' deleted successfully.")
	color.Unset()
	return nil
}

// listOrganization validates the input command line arguments and lists the available
// organizations by invoking the corresponding rpc/storage method
func listOrganization(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
	fmt.Println("This will list all the organizations available.")
	if len(args) > 0 {
		defer color.Set(color.FgRed, color.Bold)
		return errors.New("too many arguments - check again")
	}
	request := &account.AccountsRequest{}
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
	fmt.Fprintln(w, "NAME\tEMAIL\t")
	fmt.Fprintln(w, "----\t----\t")
	for _, info := range response.Accounts {
		fmt.Fprintf(w, "%s\t%s\t\n", info.Name, info.Email)
	}
	w.Flush()
	return nil
}

// addOrgMember validates the input command line arguments and adds new members to
// an organization by invoking the corresponding rpc/storage method
func addOrgMember(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
	fmt.Println("This will add new members to the specified organization.")
	org, er := getOrgName(args)
	if er != nil {
		return fmt.Errorf("user error: %v", er)
	}
	members := getMembers()
	request := &account.OrganizationMembershipsRequest{
		Name:    org,
		Members: members,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	_, err = client.AddOrganizationMemberships(context.Background(), request)
	if err != nil {
		return fmt.Errorf("server error: %v", err)
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("Members successfully added to organization '", org, "'.")
	color.Unset()
	return nil
}

// removeOrgMember validates the input command line arguments and removes members from
// an organization by invoking the corresponding rpc/storage method
func removeOrgMember(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
	fmt.Println("This will remove members from the specified organization.")
	org, er := getOrgName(args)
	if er != nil {
		return fmt.Errorf("user error: %v", er)
	}
	members := getMembers()
	request := &account.OrganizationMembershipsRequest{
		Name:    org,
		Members: members,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	_, err = client.DeleteOrganizationMemberships(context.Background(), request)
	if err != nil {
		return fmt.Errorf("server error: %v", err)
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("Members successfully removed from organization '", org, "'.")
	color.Unset()
	return nil
}

// transferOwnership validates the input command line arguments and transfer ownership of a resource
// by invoking the corresponding rpc/storage method
func transferOwnership(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
	fmt.Println("This will transfer ownership of a resource to the specified organization.")
	var resId, org string
	switch len(args) {
	case 0:
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify resource id")
		color.Unset()
		resId = getResourceID()
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify organization")
		color.Unset()
		org = getOrganization()
	case 1:
		resId = args[0]
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify organization")
		color.Unset()
		org = getOrganization()
	case 2:
		resId = args[0]
		org = args[1]
	default:
		defer color.Set(color.FgRed, color.Bold)
		return errors.New("too many arguments - check again")
	}
	request := &account.PermissionRequest{
		ResourceId:   resId,
		Organization: org,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	_, err = client.TransferOwnership(context.Background(), request)
	if err != nil {
		return fmt.Errorf("server error: %v", err)
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("Ownership Transfer successful for organization '", org, "' for resource '", resId, "'.")
	color.Unset()
	return nil
}

func getOrganization() (org string) {
	fmt.Print("Organization: ")
	color.Set(color.FgGreen, color.Bold)
	fmt.Scan(&org)
	color.Unset()
	err := account.CheckOrganizationName(org)
	if err != nil {
		color.Set(color.FgRed, color.Bold)
		fmt.Println("Organization Name is mandatory. Try again!")
		color.Unset()
		fmt.Println("")
		return getOrganization()
	}
	return
}

func getMembers() (memArr []string) {
	fmt.Print("Member name(s): ")
	color.Set(color.FgGreen, color.Bold)
	reader := bufio.NewReader(os.Stdin)
	members, _ := reader.ReadString('\n')
	memArr = strings.Fields(members)
	color.Unset()
	err := account.CheckMembers(memArr)
	if err != nil {
		color.Set(color.FgRed, color.Bold)
		fmt.Println("At least one member is mandatory. Try again!")
		color.Unset()
		fmt.Println("")
		return getMembers()
	}
	return
}

func getResourceID() (id string) {
	fmt.Print("Resource ID: ")
	color.Set(color.FgGreen, color.Bold)
	fmt.Scan(&id)
	color.Unset()
	err := account.CheckResourceID(id)
	if err != nil {
		color.Set(color.FgRed, color.Bold)
		fmt.Println("Resource ID is mandatory. Try again!")
		color.Unset()
		fmt.Println("")
		return getResourceID()
	}
	return
}

func getOrgName(args []string) (org string, err error) {
	switch len(args) {
	case 0:
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify organization")
		color.Unset()
		org = getOrganization()
	case 1:
		org = args[0]
	default:
		defer color.Set(color.FgRed, color.Bold)
		return "", errors.New("too many arguments - check again")
	}
	return
}

func getOrgNameEmail(args []string) (org string, email string, err error) {
	switch len(args) {
	case 0:
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify organization")
		color.Unset()
		org = getOrganization()
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify email")
		color.Unset()
		email, err = getEmailAddress()
	case 1:
		org = args[0]
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify email")
		color.Unset()
		email, err = getEmailAddress()
	case 2:
		org = args[0]
		email = args[1]
		mail, er := account.CheckEmailAddress(email)
		if er != nil {
			return "", "", fmt.Errorf("user error: %v", er)
		} else {
			email = mail
		}
	default:
		defer color.Set(color.FgRed, color.Bold)
		return "", "", errors.New("too many arguments - check again")
	}
	return
}
