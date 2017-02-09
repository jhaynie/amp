package tests

//import (
//	"fmt"
//	"testing"
//	"time"
//
//	"github.com/appcelerator/amp/api/rpc/account"
//	"github.com/appcelerator/amp/data/account/schema"
//)
//
//var (
//	accountName     string
//	accountEmail    string
//	accountPwd      string
//	accountUserType = schema.AccountType_USER
//)
//
////Test two stacks life cycle in the same time
//func TestAccountShouldSignUpAndVerify(t *testing.T) {
//	accountName = fmt.Sprintf("user%d", time.Now().Unix())
//	accountPwd = "pwd"
//	accountEmail = "amp@axway.com"
//
//	//SignUp      Billing
//	signUpReply, signUpErr := accountClient.SignUp(ctx, &account.SignUpRequest{
//		UserName:    accountName,
//		Password:    accountPwd,
//		Email:       accountEmail,
//		AccountType: accountUserType,
//	})
//	if signUpErr != nil {
//		t.Fatal(signUpErr)
//	}
//	if signUpReply.Token == "" {
//		t.Fatal("invalid token")
//	}
//
//	token := "bad"
//	//Verify with a bad code
//	_, verifyErr1 := accountClient.Verify(ctx, &account.VerificationRequest{
//		Token: token,
//	})
//	if verifyErr1 == nil {
//		t.Fatal("Should have error executing Verify with a bad code")
//	}
//
//	//a way to get the code...
//	_, verifyErr2 := accountClient.Verify(ctx, &account.VerificationRequest{
//		Token: signUpReply.Token,
//	})
//	if verifyErr2 != nil {
//		t.Fatal(verifyErr2)
//	}
//}
//
//func TestAccountShouldLogInResetPwdChangeItLogInAgain(t *testing.T) {
//	//Login with bad pwd
//	_, badLoginErr := accountClient.Login(ctx, &account.LogInRequest{
//		UserName: accountName,
//		Password: "badPwd",
//	})
//	if badLoginErr == nil {
//		t.Fatal("Shoud have error, login with a bad pwd")
//	}
//
//	//Login with the right pwd
//	loginAnswer1, loginErr := accountClient.Login(ctx, &account.LogInRequest{
//		UserName: accountName,
//		Password: accountPwd,
//	})
//	if loginErr != nil {
//		t.Fatal(loginErr)
//	}
//	if loginAnswer1.Token == "" {
//		t.Fatal("invalide sessionKey")
//	}
//	//Reset password
//	_, pwdResetErr := accountClient.PasswordReset(ctx, &account.PasswordResetRequest{
//		UserName: accountName,
//		Email:    accountEmail,
//	})
//	if pwdResetErr != nil {
//		t.Fatal(pwdResetErr)
//	}
//
//	//Change password
//	req := &account.PasswordChangeRequest{
//		UserName:         accountName,
//		ExistingPassword: accountPwd,
//	}
//	accountPwd = "newPwd"
//	req.NewPassword = accountPwd
//	_, pwdChangeErr := accountClient.PasswordChange(ctx, req)
//	if pwdChangeErr != nil {
//		t.Fatal(pwdChangeErr)
//	}
//
//	//Login again using new password
//	loginAnswer2, loginErr := accountClient.Login(ctx, &account.LogInRequest{
//		UserName: accountName,
//		Password: accountPwd,
//	})
//	if loginErr != nil {
//		t.Fatal(loginErr)
//	}
//	if loginAnswer2.Token == "" {
//		t.Fatal("invalide sessionKey")
//	}
//}
//
//func TestAccountShouldEditAccount(t *testing.T) {
//	//Edit email account
//	accountEmail = "ampnew@axway.com"
//	_, editErr1 := accountClient.Edit(ctx, &account.EditAccountRequest{
//		UserName: accountName,
//		Email:    accountEmail,
//	})
//	if editErr1 != nil {
//		t.Fatal(editErr1)
//	}
//
//	//Shouldn't be able to log without email verified
//	_, loginErr := accountClient.Login(ctx, &account.LogInRequest{
//		UserName: accountName,
//		Password: accountPwd,
//	})
//	if loginErr == nil {
//		t.Fatal("Shouldn't be able to log without email verified")
//	}
//	////Verify email
//	////a way to get the code...
//	//code := "ok"
//	//_, verifyErr := accountClient.Verify(ctx, &account.VerificationRequest{
//	//	Name:     accountName,
//	//	Password: accountPwd,
//	//	Code:     code,
//	//})
//	//if verifyErr != nil {
//	//	t.Fatal("invalide verify code")
//	//}
//	settingMap := make(map[string]string)
//	settingMap["testkey"] = "testval"
//	//edit billing and settings
//	_, editErr2 := accountClient.Edit(ctx, &account.EditAccountRequest{
//		UserName: accountName,
//		Settings: &account.Settings{
//			Param: settingMap,
//		},
//	})
//	if editErr2 != nil {
//		t.Fatal(editErr2)
//	}
//}
//
//func TestAccountShouldListAndGetAccountDetails(t *testing.T) {
//	//List existing accounts
//	list, listErr := accountClient.List(ctx, &account.ListAccountRequest{})
//	if listErr != nil {
//		t.Fatal(listErr)
//	}
//	found := false
//	for _, acc := range list.Accounts {
//		if acc.UserName == accountName {
//			found = true
//			break
//		}
//	}
//	if !found {
//		t.Fatal("List error, juste created account not found")
//	}
//	ret, getDetailsErr := accountClient.GetDetails(ctx, &account.GetAccountDetailsRequest{
//		UserName: accountName,
//	})
//	if getDetailsErr != nil {
//		t.Fatal(getDetailsErr)
//	}
//	if ret.Account.UserName != accountName {
//		t.Fatalf("Account.Name should be: %s\n", accountName)
//	}
//	if ret.Account.Email != accountEmail {
//		t.Fatalf("Account.Email should be: %s\n", accountEmail)
//	}
//	if !ret.Account.EmailVerified {
//		t.Fatalf("Account.EmailVerified should be: true\n")
//	}
//	if ret.Account.AccountType != accountUserType {
//		t.Fatalf("Account.AccountType should be: %s\n", accountUserType)
//	}
//	if val, ok := ret.Settings.Param["testkey"]; !ok {
//		t.Fatalf("Setting['testKey'] doesn't exist")
//	} else if val != "testval" {
//		t.Fatalf("Settings['testkey] should be: testval\n")
//	}
//}
//
//func TestAccountShouldRemoveAccount(t *testing.T) {
//	//Remove account
//	_, removeErr := accountClient.Delete(ctx, &account.DeleteAccountRequest{
//		UserName: accountName,
//	})
//	if removeErr != nil {
//		t.Fatal(removeErr)
//	}
//	//try to login again
//	//Shouldn't be able to log without eamil verified
//	_, loginErr := accountClient.Login(ctx, &account.LogInRequest{
//		UserName: accountName,
//		Password: accountPwd,
//	})
//	if loginErr == nil {
//		t.Fatal("Shouldn't be able to login after the account has been deleted")
//	}
//}
