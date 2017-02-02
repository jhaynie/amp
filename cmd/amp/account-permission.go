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
	grantPermissionCmd = &cobra.Command{
		Use:   "grant RESOURCE-ID PERMISSION-LEVEL ORGANIZATION-NAME TEAM-NAME",
		Short: "Grant permission",
		Long: `The grant command permits an account owner to grant permissions to a team having access to specified resource.
	The permissions levels can be Read, Write, Read/Write and Delete.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return grantPermission(AMP, cmd, args)
		},
	}

	editPermissionCmd = &cobra.Command{
		Use:   "edit RESOURCE-ID PERMISSION-LEVEL ORGANIZATION-NAME TEAM-NAME",
		Short: "Edit permission",
		Long: `The edit command permits an account owner to edit permissions of a team having access to specified resource.
	The permissions levels can be Read, Write, Read/Write and Delete.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return editPermission(AMP, cmd, args)
		},
	}

	revokePermissionCmd = &cobra.Command{
		Use:   "revoke RESOURCE-ID ORGANIZATION-NAME TEAM-NAME",
		Short: "Revoke permission",
		Long:  `The revoke command allows an account owner to revoke permissions of a team from the specified resource.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return revokePermission(AMP, cmd, args)
		},
	}

	listPermissionsCmd = &cobra.Command{
		Use:   "list ORGANIZATION-NAME TEAM-NAME",
		Short: "List permission",
		Long:  `The list permission command displays the permissions, filtered by teams in an organization.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return listPermission(AMP, cmd, args)
		},
	}

	permissionCmd = &cobra.Command{
		Use:   "permission",
		Short: "Permission operations",
		Long:  `The permission command manages all permission-related operations.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return permission(AMP, cmd, args)
		},
	}
)

func init() {
	AccountCmd.AddCommand(permissionCmd)
	permissionCmd.AddCommand(grantPermissionCmd)
	permissionCmd.AddCommand(editPermissionCmd)
	permissionCmd.AddCommand(revokePermissionCmd)
	permissionCmd.AddCommand(listPermissionsCmd)
}

func permission(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
	color.Set(color.FgYellow, color.Bold)
	fmt.Println("Choose a command for permission operation")
	fmt.Println("Use amp account permission -h for help")
	color.Unset()
	return nil
}

// grantPermission validates the input command line arguments and grants permissions
// to a team in an organization by invoking the corresponding rpc/storage method
func grantPermission(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
	fmt.Println("This will grant permissions to the specified resource.")
	resId, level, org, team, er := getArgs(args)
	if er != nil {
		return fmt.Errorf("user error: %v", er)
	}
	request := &account.PermissionRequest{
		ResourceId:   resId,
		Level:        level,
		Organization: org,
		Team:         team,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	_, err = client.GrantPermission(context.Background(), request)
	if err != nil {
		return fmt.Errorf("server error: %v", err)
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("Grant Permission successful for team '", team, "' for resource '", resId, "'.")
	color.Unset()
	return nil
}

// editPermission validates the input command line arguments and edits permissions
// of a team in an organization by invoking the corresponding rpc/storage method
func editPermission(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
	fmt.Println("This will edit permissions of the specified resource.")
	resId, level, org, team, er := getArgs(args)
	if er != nil {
		return fmt.Errorf("user error: %v", er)
	}
	request := &account.PermissionRequest{
		ResourceId:   resId,
		Level:        level,
		Organization: org,
		Team:         team,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	_, err = client.EditPermission(context.Background(), request)
	if err != nil {
		return fmt.Errorf("server error: %v", err)
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("Edit Permission successful for team '", team, "' for resource '", resId, "'.")
	color.Unset()
	return nil
}

// revokePermission validates the input command line arguments and revokes permissions
// of a team in an organization by invoking the corresponding rpc/storage method
func revokePermission(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
	fmt.Println("This will revoke permissions from the specified resource.")
	var resId, org, team string
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
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify team")
		color.Unset()
		team = getTeam()
	case 1:
		resId = args[0]
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify organization")
		color.Unset()
		org = getOrganization()
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify team")
		color.Unset()
		team = getTeam()
	case 2:
		resId = args[0]
		org = args[1]
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify team")
		color.Unset()
		team = getTeam()
	case 3:
		resId = args[0]
		org = args[1]
		team = args[2]
	default:
		defer color.Set(color.FgRed, color.Bold)
		return errors.New("too many arguments - check again")
	}
	request := &account.PermissionRequest{
		ResourceId:   resId,
		Organization: org,
		Team:         team,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	_, err = client.RevokePermission(context.Background(), request)
	if err != nil {
		return fmt.Errorf("server error: %v", err)
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("Revoke Permission successful for team '", team, "' for resource '", resId, "'.")
	color.Unset()
	return nil
}

// listPermission validates the input command line arguments and lists the permissions
// by invoking the corresponding rpc/storage method
func listPermission(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
	fmt.Println("This will list the permissions of a team within an organization.")
	org, team, e := getTeamOrg(args)
	if e != nil {
		return fmt.Errorf("user error: %v", e)
	}
	request := &account.PermissionRequest{
		Organization: org,
		Team:         team,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	response, er := client.ListPermissions(context.Background(), request)
	if er != nil {
		return fmt.Errorf("server error: %v", er)
	}
	if response == nil || len(response.Permissions) == 0 {
		color.Set(color.FgRed, color.Bold)
		fmt.Println("No information available!")
		color.Unset()
		return nil
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("")
	color.Unset()
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	fmt.Fprintln(w, "RESOURCE ID\tLEVEL\tTEAM\tORGANIZATION\t")
	fmt.Fprintln(w, "-----------\t-----\t----\t------------\t")
	for _, info := range response.Permissions {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t\n", info.ResourceId, info.Level, info.Team, info.Organization)
	}
	w.Flush()
	return nil
}

func getPermissionLevel() (level string) {
	fmt.Print("Permission Level: ")
	color.Set(color.FgGreen, color.Bold)
	fmt.Scan(&level)
	color.Unset()
	err := account.CheckPermissionLevel(level)
	if err != nil {
		color.Set(color.FgRed, color.Bold)
		fmt.Println("Permission Level is mandatory. Try again!")
		color.Unset()
		fmt.Println("")
		return getPermissionLevel()
	}
	return
}

func getArgs(args []string) (resId string, level string, org string, team string, err error) {
	switch len(args) {
	case 0:
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify resource id")
		color.Unset()
		resId = getResourceID()
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify permission level")
		color.Unset()
		level = getPermissionLevel()
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify organization")
		color.Unset()
		org = getOrganization()
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify team")
		color.Unset()
		team = getTeam()
	case 1:
		resId = args[0]
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify permission level")
		color.Unset()
		level = getPermissionLevel()
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify organization")
		color.Unset()
		org = getOrganization()
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify team")
		color.Unset()
		team = getTeam()
	case 2:
		resId = args[0]
		level = args[1]
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify organization")
		color.Unset()
		org = getOrganization()
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify team")
		color.Unset()
		team = getTeam()
	case 3:
		resId = args[0]
		level = args[1]
		org = args[2]
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify team")
		color.Unset()
		team = getTeam()
	case 4:
		resId = args[0]
		level = args[1]
		org = args[2]
		team = args[3]
	default:
		defer color.Set(color.FgRed, color.Bold)
		return "", "", "", "", errors.New("too many arguments - check again")
	}
	return
}
