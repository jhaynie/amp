package tests

import (
	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/data/account/schema"
	"github.com/docker/distribution/context"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	signUpRequestOrganization = account.SignUpRequest{
		UserName:    "Organization",
		Password:    "OrganizationPassword",
		Email:       "organization@amp.io",
		AccountType: schema.AccountType_ORGANIZATION,
	}

	signUpRequestUser = account.SignUpRequest{
		UserName:    "User",
		Password:    "UserPassword",
		Email:       "user@amp.io",
		AccountType: schema.AccountType_USER,
	}
)

func TestOrganizationSignUpInvalidUserNameShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	invalidSignUp := signUpRequestOrganization
	invalidSignUp.UserName = ""
	_, signUpErr := accountClient.SignUp(ctx, &invalidSignUp)
	assert.Error(t, signUpErr)
}

func TestOrganizationSignUpInvalidEmailShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	invalidSignUp := signUpRequestOrganization
	invalidSignUp.Email = "this is not an email"
	_, signUpErr := accountClient.SignUp(ctx, &invalidSignUp)
	assert.Error(t, signUpErr)
}

func TestOrganizationSignUpInvalidPasswordShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	invalidSignUp := signUpRequestOrganization
	invalidSignUp.Password = ""
	_, signUpErr := accountClient.SignUp(ctx, &invalidSignUp)
	assert.Error(t, signUpErr)
}

func TestOrganizationShouldSignUpAndVerify(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	signUpReply, signUpErr := accountClient.SignUp(ctx, &signUpRequestOrganization)
	assert.NoError(t, signUpErr)

	// Verify
	_, verifyErr := accountClient.Verify(ctx, &account.VerificationRequest{Token: signUpReply.Token})
	assert.NoError(t, verifyErr)
}

func TestOrganizationSignUpAlreadyExistsShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	_, err1 := accountClient.SignUp(ctx, &signUpRequestOrganization)
	assert.NoError(t, err1)

	// SignUp
	_, err2 := accountClient.SignUp(ctx, &signUpRequestOrganization)
	assert.Error(t, err2)
}

func TestOrganizationVerifyNotATokenShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	_, signUpErr := accountClient.SignUp(ctx, &signUpRequestOrganization)
	assert.NoError(t, signUpErr)

	// Verify
	_, verifyErr := accountClient.Verify(ctx, &account.VerificationRequest{Token: "this is not a token"})
	assert.Error(t, verifyErr)
}

// TODO: Check token with invalid signature
// TODO: Check token with invalid non existing account id
// TODO: Check expired token

func TestOrganizationLogin(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	signUpReply, signUpErr := accountClient.SignUp(ctx, &signUpRequestOrganization)
	assert.NoError(t, signUpErr)

	// Verify
	_, verifyErr := accountClient.Verify(ctx, &account.VerificationRequest{Token: signUpReply.Token})
	assert.NoError(t, verifyErr)

	// Login
	_, loginErr := accountClient.Login(ctx, &account.LogInRequest{
		UserName: signUpRequestOrganization.UserName,
		Password: signUpRequestOrganization.Password,
	})
	assert.NoError(t, loginErr)
}

// TODO: Check login with non existing account
// TODO: Check login with non verified account

func TestOrganizationLoginInvalidUserNameShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	signUpReply, signUpErr := accountClient.SignUp(ctx, &signUpRequestOrganization)
	assert.NoError(t, signUpErr)

	// Verify
	_, verifyErr := accountClient.Verify(ctx, &account.VerificationRequest{Token: signUpReply.Token})
	assert.NoError(t, verifyErr)

	// Login
	_, loginErr := accountClient.Login(ctx, &account.LogInRequest{
		UserName: "not the right user name",
		Password: signUpRequestOrganization.Password,
	})
	assert.Error(t, loginErr)
}

func TestOrganizationLoginInvalidPasswordShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	signUpReply, signUpErr := accountClient.SignUp(ctx, &signUpRequestOrganization)
	assert.NoError(t, signUpErr)

	// Verify
	_, verifyErr := accountClient.Verify(ctx, &account.VerificationRequest{Token: signUpReply.Token})
	assert.NoError(t, verifyErr)

	// Login
	_, loginErr := accountClient.Login(ctx, &account.LogInRequest{
		UserName: signUpRequestOrganization.UserName,
		Password: "not the right password",
	})
	assert.Error(t, loginErr)
}
