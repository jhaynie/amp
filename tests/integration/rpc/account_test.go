package tests

import (
	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/data/account/schema"
	"github.com/docker/distribution/context"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	signUpRequestUser = account.SignUpRequest{
		UserName:    "User",
		Password:    "UserPassword",
		Email:       "user@amp.io",
		AccountType: schema.AccountType_USER,
	}
)

func TestUserSignUpInvalidUserNameShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	invalidSignUp := signUpRequestUser
	invalidSignUp.UserName = ""
	_, signUpErr := accountClient.SignUp(ctx, &invalidSignUp)
	assert.Error(t, signUpErr)
}

func TestUserSignUpInvalidEmailShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	invalidSignUp := signUpRequestUser
	invalidSignUp.Email = "this is not an email"
	_, signUpErr := accountClient.SignUp(ctx, &invalidSignUp)
	assert.Error(t, signUpErr)
}

func TestUserSignUpInvalidPasswordShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	invalidSignUp := signUpRequestUser
	invalidSignUp.Password = ""
	_, signUpErr := accountClient.SignUp(ctx, &invalidSignUp)
	assert.Error(t, signUpErr)
}

func TestUserShouldSignUpAndVerify(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	signUpReply, signUpErr := accountClient.SignUp(ctx, &signUpRequestUser)
	assert.NoError(t, signUpErr)

	// Verify
	_, verifyErr := accountClient.Verify(ctx, &account.VerificationRequest{Token: signUpReply.Token})
	assert.NoError(t, verifyErr)
}

func TestUserSignUpAlreadyExistsShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	_, err1 := accountClient.SignUp(ctx, &signUpRequestUser)
	assert.NoError(t, err1)

	// SignUp
	_, err2 := accountClient.SignUp(ctx, &signUpRequestUser)
	assert.Error(t, err2)
}

func TestUserVerifyNotATokenShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	_, signUpErr := accountClient.SignUp(ctx, &signUpRequestUser)
	assert.NoError(t, signUpErr)

	// Verify
	_, verifyErr := accountClient.Verify(ctx, &account.VerificationRequest{Token: "this is not a token"})
	assert.Error(t, verifyErr)
}

// TODO: Check token with invalid signature
// TODO: Check token with non existing account id
// TODO: Check expired token

func TestUserLogin(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	signUpReply, signUpErr := accountClient.SignUp(ctx, &signUpRequestUser)
	assert.NoError(t, signUpErr)

	// Verify
	_, verifyErr := accountClient.Verify(ctx, &account.VerificationRequest{Token: signUpReply.Token})
	assert.NoError(t, verifyErr)

	// Login
	_, loginErr := accountClient.Login(ctx, &account.LogInRequest{
		UserName: signUpRequestUser.UserName,
		Password: signUpRequestUser.Password,
	})
	assert.NoError(t, loginErr)
}

func TestUserLoginNonExistingAccountShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Login
	_, loginErr := accountClient.Login(ctx, &account.LogInRequest{
		UserName: signUpRequestUser.UserName,
		Password: signUpRequestUser.Password,
	})
	assert.Error(t, loginErr)
}

func TestUserLoginNonVerifiedAccountShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	_, signUpErr := accountClient.SignUp(ctx, &signUpRequestUser)
	assert.NoError(t, signUpErr)

	// Login
	_, loginErr := accountClient.Login(ctx, &account.LogInRequest{
		UserName: signUpRequestUser.UserName,
		Password: signUpRequestUser.Password,
	})
	assert.Error(t, loginErr)
}

func TestUserLoginInvalidUserNameShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	signUpReply, signUpErr := accountClient.SignUp(ctx, &signUpRequestUser)
	assert.NoError(t, signUpErr)

	// Verify
	_, verifyErr := accountClient.Verify(ctx, &account.VerificationRequest{Token: signUpReply.Token})
	assert.NoError(t, verifyErr)

	// Login
	_, loginErr := accountClient.Login(ctx, &account.LogInRequest{
		UserName: "not the right user name",
		Password: signUpRequestUser.Password,
	})
	assert.Error(t, loginErr)
}

func TestUserLoginInvalidPasswordShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	signUpReply, signUpErr := accountClient.SignUp(ctx, &signUpRequestUser)
	assert.NoError(t, signUpErr)

	// Verify
	_, verifyErr := accountClient.Verify(ctx, &account.VerificationRequest{Token: signUpReply.Token})
	assert.NoError(t, verifyErr)

	// Login
	_, loginErr := accountClient.Login(ctx, &account.LogInRequest{
		UserName: signUpRequestUser.UserName,
		Password: "not the right password",
	})
	assert.Error(t, loginErr)
}

func TestUserPasswordReset(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	signUpReply, signUpErr := accountClient.SignUp(ctx, &signUpRequestUser)
	assert.NoError(t, signUpErr)

	// Verify
	_, verifyErr := accountClient.Verify(ctx, &account.VerificationRequest{Token: signUpReply.Token})
	assert.NoError(t, verifyErr)

	// Password Reset
	_, resetPasswordErr := accountClient.PasswordReset(ctx, &account.PasswordResetRequest{UserName: signUpRequestUser.UserName})
	assert.NoError(t, resetPasswordErr)
}

func TestUserPasswordResetNonExistingAccountShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	signUpReply, signUpErr := accountClient.SignUp(ctx, &signUpRequestUser)
	assert.NoError(t, signUpErr)

	// Verify
	_, verifyErr := accountClient.Verify(ctx, &account.VerificationRequest{Token: signUpReply.Token})
	assert.NoError(t, verifyErr)

	// Password Reset
	_, resetPasswordErr := accountClient.PasswordReset(ctx, &account.PasswordResetRequest{UserName: "This is not an existing account"})
	assert.Error(t, resetPasswordErr)
}

func TestUserPasswordSet(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	signUpReply, signUpErr := accountClient.SignUp(ctx, &signUpRequestUser)
	assert.NoError(t, signUpErr)

	// Verify
	_, verifyErr := accountClient.Verify(ctx, &account.VerificationRequest{Token: signUpReply.Token})
	assert.NoError(t, verifyErr)

	// Password Reset
	passwordResetReply, resetPasswordErr := accountClient.PasswordReset(ctx, &account.PasswordResetRequest{UserName: signUpRequestUser.UserName})
	assert.NoError(t, resetPasswordErr)

	// Password Set
	_, setPasswordErr := accountClient.PasswordSet(ctx, &account.PasswordSetRequest{
		Token:    passwordResetReply.Token,
		Password: "newPassword",
	})
	assert.NoError(t, setPasswordErr)

	// Login
	_, loginErr := accountClient.Login(ctx, &account.LogInRequest{
		UserName: signUpRequestUser.UserName,
		Password: "newPassword",
	})
	assert.NoError(t, loginErr)
}

func TestUserPasswordSetInvalidTokenShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	signUpReply, signUpErr := accountClient.SignUp(ctx, &signUpRequestUser)
	assert.NoError(t, signUpErr)

	// Verify
	_, verifyErr := accountClient.Verify(ctx, &account.VerificationRequest{Token: signUpReply.Token})
	assert.NoError(t, verifyErr)

	// Password Reset
	_, resetPasswordErr := accountClient.PasswordReset(ctx, &account.PasswordResetRequest{UserName: signUpRequestUser.UserName})
	assert.NoError(t, resetPasswordErr)

	// Password Set
	_, setPasswordErr := accountClient.PasswordSet(ctx, &account.PasswordSetRequest{
		Token:    "this is an invalid token",
		Password: "newPassword",
	})
	assert.Error(t, setPasswordErr)

	// Login
	_, loginErr := accountClient.Login(ctx, &account.LogInRequest{
		UserName: signUpRequestUser.UserName,
		Password: "newPassword",
	})
	assert.Error(t, loginErr)
}

func TestUserPasswordSetInvalidPasswordShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	signUpReply, signUpErr := accountClient.SignUp(ctx, &signUpRequestUser)
	assert.NoError(t, signUpErr)

	// Verify
	_, verifyErr := accountClient.Verify(ctx, &account.VerificationRequest{Token: signUpReply.Token})
	assert.NoError(t, verifyErr)

	// Password Reset
	passwordResetReply, resetPasswordErr := accountClient.PasswordReset(ctx, &account.PasswordResetRequest{UserName: signUpRequestUser.UserName})
	assert.NoError(t, resetPasswordErr)

	// Password Set
	_, setPasswordErr := accountClient.PasswordSet(ctx, &account.PasswordSetRequest{
		Token:    passwordResetReply.Token,
		Password: "",
	})
	assert.Error(t, setPasswordErr)

	// Login
	_, loginErr := accountClient.Login(ctx, &account.LogInRequest{
		UserName: signUpRequestUser.UserName,
		Password: "",
	})
	assert.Error(t, loginErr)
}
