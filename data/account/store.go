package account

import (
	"context"
	"fmt"
	"github.com/appcelerator/amp/data/schema/account"
	"github.com/appcelerator/amp/data/storage"
	"github.com/docker/docker/pkg/stringid"
	"path"
	"strings"
)

const accountsRootKey = "accounts"

// Store implements account data.Interface
type Store struct {
	Store storage.Interface
}

// NewStore returns an etcd implementation of account.Interface
func NewStore(store storage.Interface) *Store {
	return &Store{
		Store: store,
	}
}

// CreateAccount creates a new account
func (s *Store) CreateAccount(ctx context.Context, in *schema.Account) (string, error) {
	alreadyExists, err := s.GetAccountByEmail(ctx, in.Email)
	if err != nil {
		return nil, err
	}
	if alreadyExists != nil {
		return nil, fmt.Errorf("account already exists: %s", in.Email)
	}
	in.Id = stringid.GenerateNonCryptoID()
	if err := s.Store.Create(ctx, path.Join(accountsRootKey, in.Id), in, nil, 0); err != nil {
		return nil, err
	}
	return in.Id, nil
}

// GetAccount fetches an account by id
func (s *Store) GetAccount(ctx context.Context, id string) (*schema.Account, error) {
	account := &schema.Account{}
	if err := s.Store.Get(ctx, path.Join(accountsRootKey, id), account, false); err != nil {
		return nil, err
	}
	return account, nil
}

// GetAccount fetches an account by email
func (s *Store) GetAccountByEmail(ctx context.Context, email string) (*schema.Account, error) {
	accounts, err := s.ListAccounts(ctx)
	if err != nil {
		return nil, err
	}
	for _, acc := range accounts {
		if strings.EqualFold(acc.Email, email) {
			return acc, nil
		}
	}
	return nil, nil
}

// ListAccounts lists accounts
func (s *Store) ListAccounts(ctx context.Context) ([]*schema.Account, error) {
	accounts := []*schema.Account{}
	if err := s.Store.List(ctx, accountsRootKey, storage.Everything, &schema.Account{}, &accounts); err != nil {
		return nil, err
	}
	return accounts, nil
}

// UpdateAccount updates an account
func (s *Store) UpdateAccount(ctx context.Context, in *schema.Account) (*schema.Account, error) {
	if err := s.Store.Put(ctx, path.Join(accountsRootKey, in.Id), in, nil, 0); err != nil {
		return nil, err
	}
	return in, nil
}

// DeleteAccount deletes an account by id
func (s *Store) DeleteAccount(ctx context.Context, id string) error {
	if err := s.Store.Delete(ctx, path.Join(accountsRootKey, id), false, nil); err != nil {
		return nil, err
	}
	return nil, nil
}
