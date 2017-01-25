package account

import (
	"context"

	"github.com/appcelerator/amp/data/schema"
)

// Interface must be implemented an account database
type Interface interface {
	// AddAccount adds a new account to the account table
	AddAccount(ctx context.Context, account *schema.Account) (id string, err error)

	// Verify sets an account verification to true
	Verify(ctx context.Context, name string) error

	// AddTeam adds a new team to the team table
	AddTeam(ctx context.Context, team *schema.Team) (id string, err error)

	//GetTeam returns a team from the team table
	GetTeam(ctx context.Context, name string) (team *schema.Team, err error)

	// AddTeamMember adds a new team to the team table
	AddTeamMember(ctx context.Context, member *schema.TeamMember) (id string, err error)

	// GetTeamMember returns a new team from the team table
	GetTeamMember(ctx context.Context, teamId string, memberId string) (member *schema.TeamMember, err error)

	//DeleteTeamMember
	DeleteTeamMember(ctx context.Context, teamId string, memberId string) (err error)

	// GetAccount returns an account from the accounts table
	GetAccount(ctx context.Context, name string) (*schema.Account, error)

	// GetAccounts returns accounts matching a query
	GetAccounts(ctx context.Context, accountType schema.AccountType) ([]*schema.Account, error)

	//AddResource Adds Resource to resource table
	AddResource(ctx context.Context, resource *schema.Resource) (id string, err error)

	//GetResource returns a team from the team table
	GetResource(ctx context.Context, name string) (team *schema.Resource, err error)

	//DeleteResource removes the Resource entry for a given Id
	DeleteResource(ctx context.Context, name string) (err error)

	//AddResource Adds Resource to resource table
	AddResourceSettings(ctx context.Context, resource *schema.ResourceSettings) (id string, err error)

	//GetResourceSettings returns a list of ResourceSettings for a given resource
	GetResourceSettings(ctx context.Context, resourceId string) (rs []*schema.ResourceSettings, err error)

	//DeleteResourceSettings removes the Resource entry for a given Id
	DeleteResourceSettings(ctx context.Context, resourceId string) (err error)

	//AddPermission Adds Permission to the Permission table
	AddPermission(ctx context.Context, resource *schema.Permission) (id string, err error)

	//GetPermission returns the permission record
	GetPermission(ctx context.Context, resourceId string) (perm []*schema.Permission, err error)
}
