package account

import (
	"context"

	"github.com/appcelerator/amp/data/storage"

	"fmt"
	"path"

	"github.com/appcelerator/amp/data/schema"
	"github.com/docker/docker/pkg/stringid"
	"github.com/golang/protobuf/proto"
)

const accountSchemaRootKey = "accounts"
const accountUserByNameKey = accountSchemaRootKey + "/account/name"
const accountTeamKey = accountSchemaRootKey + "/team"
const accountTeamByNameKey = accountSchemaRootKey + "/team/name"
const accountUserKey = accountSchemaRootKey + "/account/id"
const accountTeamMemberKey = accountSchemaRootKey + "/team/member"
const accountResourceKey = accountSchemaRootKey + "/resource"
const accountResourceByNameKey = accountSchemaRootKey + "/resource/name"
const accountResourceSettingsKey = accountSchemaRootKey + "/resource/settings"
const accountPermissionKey = accountSchemaRootKey + "/permission"

// Store impliments account data.Interface
type Store struct {
	Store storage.Interface
}

// NewStore returns a Storage wrapper with functions to operate against the backing database
func NewStore(store storage.Interface) *Store {
	return &Store{
		Store: store,
	}
}

//AddResource adds a Resource to the resource table
func (s *Store) AddResource(ctx context.Context, resource *schema.Resource) (id string, err error) {

	// Integrity Check
	if err = s.checkResource(resource); err == nil {
		// Store the Account struct and the alternate key
		if err = s.Store.Create(ctx, path.Join(accountResourceKey, resource.Id), resource, nil, 0); err == nil {
			err = s.Store.Create(ctx, path.Join(accountResourceByNameKey, resource.Name), &schema.Key{KeyId: resource.Id}, nil, 0)
		}
	}
	return resource.Id, err
}

func (s *Store) checkResource(resource *schema.Resource) (err error) {

	// Assign a new UUID if the resource does not exist
	if resource.Id == "" {
		resource.Id = generateUUID()

	} else {
		err = fmt.Errorf("Resource.Id must be blank/nil ")
	}
	return err
}

// GetResource returns a Resource from the Resource table
func (s *Store) GetResource(ctx context.Context, name string) (resource *schema.Resource, err error) {
	key := &schema.Key{}
	//Grab the ID
	err = s.Store.Get(ctx, path.Join(accountResourceByNameKey, name), key, true)
	if err == nil && key.KeyId != "" {
		resource, err = s.GetResourceById(ctx, key.KeyId)
	}
	return resource, err
}

// GetResource returns a Resource from the Resource table
func (s *Store) GetResourceById(ctx context.Context, id string) (resource *schema.Resource, err error) {
	resource = &schema.Resource{}
	err = s.Store.Get(ctx, path.Join(accountResourceKey+"/"+id), resource, true)
	return resource, err
}

// GetResourceSettings returns a List of ResourceSettings from the ResourceSettings table
func (s *Store) GetResourceSettings(ctx context.Context, resourceId string) (rs []*schema.ResourceSettings, err error) {
	var out []proto.Message
	settings := &schema.ResourceSettings{}
	err = s.Store.List(ctx, accountResourceSettingsKey+"/"+resourceId, storage.Everything, settings, &out)
	if err == nil {
		// Unfortunately we have to iterate and filter
		for i := 0; i < len(out); i++ {
			m, ok := out[i].(*schema.ResourceSettings)
			if ok {
				rs = append(rs, m)
			} else if !ok {
				err = fmt.Errorf("Unexpected Type Encountered")
				break
			}
		}
	}
	return
}

//AddTeamMember adds a user account to the team table
func (s *Store) AddResourceSettings(ctx context.Context, rs *schema.ResourceSettings) (id string, err error) {

	if err = s.checkResourceSettings(rs); err == nil {
		// Store the Account struct and the alternate key
		err = s.Store.Create(ctx, path.Join(accountResourceSettingsKey, rs.ResourceId+"/"+rs.Id), rs, nil, 0)
	}
	return rs.Id, err
}

func (s *Store) checkResourceSettings(rs *schema.ResourceSettings) (err error) {

	if rs.Id == "" {
		rs.Id = generateUUID()

	} else {
		err = fmt.Errorf("Resource Id must not be populated")
	}
	return err
}

//AddTeamMember adds a user account to the team table
func (s *Store) AddTeamMember(ctx context.Context, member *schema.TeamMember) (id string, err error) {

	if err = s.checkTeamMember(ctx, member); err == nil {
		// Store the Account struct and the alternate key
		err = s.Store.Create(ctx, path.Join(accountTeamMemberKey, member.TeamId+"/"+member.Id), member, nil, 0)
	}
	return member.Id, err
}

func (s *Store) checkTeamMember(ctx context.Context, member *schema.TeamMember) error {

	mem, err := s.GetTeamMember(ctx, member.TeamId, member.Id)
	if err == nil && member.Id == "" {
		member.Id = generateUUID()

	} else {
		err = fmt.Errorf("TeamMember %s already exists", mem.Id)
	}
	return err
}

// GetTeamMember returns a TeamMember from the TeamMember table
func (s *Store) GetTeamMember(ctx context.Context, teamId string, memberId string) (member *schema.TeamMember, err error) {
	member = &schema.TeamMember{}
	key := &schema.Key{KeyId: memberId}
	err = s.Store.Get(ctx, path.Join(accountTeamMemberKey+"/"+teamId, key.KeyId), member, true)
	return member, err
}

// generateUUID place holder until we standardize the approach we want to use
func generateUUID() (id string) {
	return stringid.GenerateNonCryptoID()
}

// AddAccount adds a new account to the account table
func (s *Store) AddAccount(ctx context.Context, account *schema.Account) (id string, err error) {

	if err = s.checkAccount(ctx, account); err == nil {
		// Store the account struct and the alternate key
		if err = s.Store.Create(ctx, path.Join(accountUserKey, account.Id), account, nil, 0); err == nil {
			fk := &schema.Key{KeyId: account.Id}
			err = s.Store.Create(ctx, path.Join(accountUserByNameKey, account.Name), fk, nil, 0)
		}
	}
	return account.Id, err
}
func (s *Store) checkAccount(ctx context.Context, account *schema.Account) error {
	acct, err := s.GetAccount(ctx, account.Name)
	if err == nil && acct.Id == "" {
		account.Id = generateUUID()
	} else {
		err = fmt.Errorf("Account %s already exists", acct.Name)
	}
	return err
}
func (s *Store) checkTeam(ctx context.Context, team *schema.Team) error {

	t, err := s.GetTeam(ctx, team.Name)
	if err == nil && t.Id == "" {
		team.Id = generateUUID()

	} else {
		err = fmt.Errorf("Team %s already exists", t.Name)
	}
	return err
}

// Verify sets an account verification to true
func (s *Store) Verify(ctx context.Context, name string) error {
	acct, err := s.GetAccount(ctx, name)
	if err == nil && acct.Name != "" && !acct.IsVerified {
		acct.IsVerified = true
		err = s.Store.Put(ctx, path.Join(accountUserKey, acct.Id), acct, 0)
	}
	return err
}

// AddTeam adds a new team to the team table
func (s *Store) AddTeam(ctx context.Context, team *schema.Team) (id string, err error) {
	//TODO Add data integrity checks
	if team.Id == "" {
		team.Id = generateUUID()
	}
	// Store Team struct and alternate Key
	if err = s.Store.Create(ctx, path.Join(accountTeamKey, team.Id), team, nil, 0); err == nil {
		err = s.Store.Create(ctx, path.Join(accountTeamByNameKey, team.Name), &schema.Key{KeyId: team.Id}, nil, 0)
	}
	return team.Id, err
}

// GetTeam returns a Team from the Team table
func (s *Store) GetTeam(ctx context.Context, name string) (team *schema.Team, err error) {
	team = &schema.Team{}
	key := &schema.Key{}
	//Grab the ID
	err = s.Store.Get(ctx, path.Join(accountTeamByNameKey, name), key, true)
	if err == nil && key.KeyId != "" {
		err = s.Store.Get(ctx, path.Join(accountTeamKey, key.KeyId), team, true)
	}
	return team, err
}

// GetAccount returns an account from the accounts table
func (s *Store) GetAccount(ctx context.Context, name string) (account *schema.Account, err error) {
	acct := &schema.Account{}
	key := &schema.Key{}
	//Grab the ID
	err = s.Store.Get(ctx, path.Join(accountUserByNameKey, name), key, true)
	if err == nil && key.KeyId != "" {
		err = s.Store.Get(ctx, path.Join(accountUserKey, key.KeyId), acct, true)
	}
	return acct, err
}

// GetAccounts implements Inrface.GetAccounts
func (s *Store) GetAccounts(ctx context.Context, accountType schema.AccountType) (accounts []*schema.Account, err error) {

	var out []proto.Message
	account := &schema.Account{}
	err = s.Store.List(ctx, accountUserKey, storage.Everything, account, &out)
	if err == nil {
		// Unfortunately we have to iterate and filter
		for i := 0; i < len(out); i++ {
			m, ok := out[i].(*schema.Account)
			if ok && m.Type == accountType {
				accounts = append(accounts, m)
			} else if !ok {
				err = fmt.Errorf("Unexpected Type Encountered")
				break
			}
		}
	}
	return
}

//AddPermission adds a permission record to the Permission Table
func (s *Store) AddPermission(ctx context.Context, perm *schema.Permission) (id string, err error) {

	if err = s.checkPermission(perm); err == nil {
		// Store the Permission struct and the alternate key
		err = s.Store.Create(ctx, path.Join(accountPermissionKey, perm.ResourceId+"/"+perm.Id), perm, nil, 0)
	}
	return perm.Id, err
}

func (s *Store) checkPermission(perm *schema.Permission) (err error) {

	if perm.Id == "" {
		perm.Id = generateUUID()

	} else {
		err = fmt.Errorf("Perm ID must be blank/nil")
	}
	return err
}

//GetPermission retrieves a collection of permission records from the permission table
func (s *Store) GetPermission(ctx context.Context, resourceId string) (perms []*schema.Permission, err error) {
	var out []proto.Message
	perm := &schema.Permission{}
	err = s.Store.List(ctx, accountPermissionKey+"/"+resourceId, storage.Everything, perm, &out)
	if err == nil {
		// Unfortunately we have to iterate and filter
		for i := 0; i < len(out); i++ {
			m, ok := out[i].(*schema.Permission)
			if ok {
				perms = append(perms, m)
			} else if !ok {
				err = fmt.Errorf("Unexpected Type Encountered")
				break
			}
		}
	}
	return
}

//DeleteResourceSettings removes the Resource entry for a given Id
func (s *Store) DeleteResourceSettings(ctx context.Context, resourceId string) (err error) {

	err = s.Store.Delete(ctx, accountResourceSettingsKey+"/"+resourceId, true, nil)
	return err
}

//DeleteResource removes the Resource entry for a given Id
func (s *Store) DeleteResource(ctx context.Context, name string) (err error) {

	if resource, err := s.GetResource(ctx, name); err == nil && resource.Id != "" {
		//Cascade delete
		s.DeleteResourceSettings(ctx, resource.Id)
		err = s.Store.Delete(ctx, accountResourceKey+"/"+resource.Id, false, nil)
	}
	return
}

//DeleteTeamMember
func (s *Store) DeleteTeamMember(ctx context.Context, teamId string, memberId string) (err error) {
	//TODO
	return
}
