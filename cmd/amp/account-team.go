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
	createTeamCmd = &cobra.Command{
		Use:   "create ORGANIZATION-NAME TEAM-NAME",
		Short: "Create team within an organization",
		Long:  `The create team command creates a team within the specified organization.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return createTeam(AMP, cmd, args)
		},
	}

	editTeamCmd = &cobra.Command{
		Use:   "edit ORGANIZATION-NAME TEAM-NAME",
		Short: "Edit team information",
		Long:  `The edit team command updates information of team within the specified organization.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return editTeam(AMP, cmd, args)
		},
	}

	deleteTeamCmd = &cobra.Command{
		Use:   "delete ORGANIZATION-NAME TEAM-NAME",
		Short: "Delete team within an organization",
		Long:  `The delete team command deletes a team within the specified organization.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return deleteTeam(AMP, cmd, args)
		},
	}

	listTeamsCmd = &cobra.Command{
		Use:   "list ORGANIZATION-NAME",
		Short: "List teams by organization",
		Long:  `The list team command displays the available teams in a specified organization.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return listTeam(AMP, cmd, args)
		},
	}

	addTeamMembersCmd = &cobra.Command{
		Use:   "add ORGANIZATION-NAME TEAM-NAME MEMBERS...",
		Short: "Add users to a team",
		Long:  `The add-team command allows an owner team member to add new members to a team in an organization.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return addTeamMember(AMP, cmd, args)
		},
	}

	removeTeamMembersCmd = &cobra.Command{
		Use:   "remove ORGANIZATION-NAME TEAM-NAME MEMBERS...",
		Short: "Remove users from a team",
		Long:  `The remove-team command allows an owner team member to remove existing members from a team in the organization.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return removeTeamMember(AMP, cmd, args)
		},
	}

	teamCmd = &cobra.Command{
		Use:   "team",
		Short: "Team operations",
		Long:  `The team command manages all team-related operations.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return team(AMP, cmd, args)
		},
	}
)

func init() {
	AccountCmd.AddCommand(teamCmd)
	teamCmd.AddCommand(createTeamCmd)
	teamCmd.AddCommand(editTeamCmd)
	teamCmd.AddCommand(deleteTeamCmd)
	teamCmd.AddCommand(listTeamsCmd)
	teamCmd.AddCommand(addTeamMembersCmd)
	teamCmd.AddCommand(removeTeamMembersCmd)
}

func team(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
	color.Set(color.FgYellow, color.Bold)
	fmt.Println("Choose a command for team operation")
	fmt.Println("Use amp account team -h for help")
	color.Unset()
	return nil
}

// createTeam validates the input command line arguments and creates a team in an organization
// with name by invoking the corresponding rpc/storage method
func createTeam(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
	fmt.Println("This will create a team within the specified organization.")
	org, team, e := getTeamOrg(args)
	if e != nil {
		return fmt.Errorf("user error: %v", e)
	}
	request := &account.TeamRequest{
		Organization: org,
		Name:         team,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	_, err = client.CreateTeam(context.Background(), request)
	if err != nil {
		return fmt.Errorf("server error: %v", err)
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("Team '", team, "' created successfully.")
	color.Unset()
	return nil
}

// editTeam validates the input command line arguments and creates a team in an organization
// by invoking the corresponding rpc/storage method
func editTeam(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
	fmt.Println("This will edit team information.")
	org, team, e := getTeamOrg(args)
	if e != nil {
		return fmt.Errorf("user error: %v", e)
	}
	request := &account.TeamRequest{
		Organization: org,
		Name:         team,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	_, err = client.EditTeam(context.Background(), request)
	if err != nil {
		return fmt.Errorf("server error: %v", err)
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("Team '", team, "' updated successfully.")
	color.Unset()
	return nil
}

// deleteTeam validates the input command line arguments and deletes a team in an organization
// by invoking the corresponding rpc/storage method
func deleteTeam(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
	fmt.Println("This will deletes a team within the specified organization.")
	org, team, e := getTeamOrg(args)
	if e != nil {
		return fmt.Errorf("user error: %v", e)
	}
	request := &account.TeamRequest{
		Organization: org,
		Name:         team,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	_, err = client.DeleteTeam(context.Background(), request)
	if err != nil {
		return fmt.Errorf("server error: %v", err)
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("Team '", team, "' deleted successfully.")
	color.Unset()
	return nil
}

// listTeam validates the input command line arguments and lists the available teams
// within an organization by invoking the corresponding rpc/storage method
func listTeam(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
	fmt.Println("This will list all the teams available in a specific organization.")
	org, e := getOrgName(args)
	if e != nil {
		return fmt.Errorf("user error: %v", e)
	}
	request := &account.TeamRequest{
		Organization: org,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	response, er := client.ListTeams(context.Background(), request)
	if er != nil {
		return fmt.Errorf("server error: %v", er)
	}
	if response == nil || len(response.Teams) == 0 {
		color.Set(color.FgRed, color.Bold)
		fmt.Println("No information available!")
		color.Unset()
		return nil
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("")
	color.Unset()
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	fmt.Fprintln(w, "NAME\tDESCRIPTION\tMEMBERS\t")
	fmt.Fprintln(w, "----\t-----------\t-------\t")
	for _, info := range response.Teams {
		fmt.Fprintf(w, "%s\t%s\t%s\t\n", info.Name, info.Description, info.Members)
	}
	w.Flush()
	return nil
}

// addTeamMember validates the input command line arguments and adds new members to
// a team in an organization by invoking the corresponding rpc/storage method
func addTeamMember(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
	fmt.Println("This will add new members to the specified team within an organization.")
	org, team, e := getTeamOrg(args)
	if e != nil {
		return fmt.Errorf("user error: %v", e)
	}
	members := getMembers()
	request := &account.TeamMembershipsRequest{
		Organization: org,
		Name:         team,
		Members:      members,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	_, err = client.AddTeamMemberships(context.Background(), request)
	if err != nil {
		return fmt.Errorf("server error: %v", err)
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("Members successfully added to team '", team, "'.")
	color.Unset()
	return nil
}

// removeTeamMember validates the input command line arguments and removes members from
// a team in an organization by invoking the corresponding rpc/storage method
func removeTeamMember(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
	fmt.Println("This will remove members from the specified team within an organization.")
	org, team, e := getTeamOrg(args)
	if e != nil {
		return fmt.Errorf("user error: %v", e)
	}
	members := getMembers()
	request := &account.TeamMembershipsRequest{
		Organization: org,
		Name:         team,
		Members:      members,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	_, err = client.DeleteTeamMemberships(context.Background(), request)
	if err != nil {
		return fmt.Errorf("server error: %v", err)
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("Members successfully removed from team '", team, "'.")
	color.Unset()
	return nil
}

func getTeam() (team string) {
	fmt.Print("Team Name: ")
	color.Set(color.FgGreen, color.Bold)
	fmt.Scan(&team)
	color.Unset()
	err := account.CheckTeamName(team)
	if err != nil {
		color.Set(color.FgRed, color.Bold)
		fmt.Println("Team name is mandatory. Try again!")
		color.Unset()
		fmt.Println("")
		return getTeam()
	}
	return
}

func getTeamOrg(args []string) (org string, team string, err error) {
	switch len(args) {
	case 0:
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify organization")
		color.Unset()
		org = getOrganization()
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify team")
		color.Unset()
		team = getTeam()
	case 1:
		org = args[0]
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify team")
		color.Unset()
		team = getTeam()
	case 2:
		org = args[0]
		team = args[1]
	default:
		defer color.Set(color.FgRed, color.Bold)
		return "", "", errors.New("too many arguments - check again")
	}
	return
}
