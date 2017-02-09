
package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/appcelerator/amp/api/rpc/account"
)

var (
	accountName string
	accountEmail string
	accountPwd string
	accountType string
)

//Test two stacks life cycle in the same time
func TestOrgShouldSignUpAndVerify(t *testing.T) {
	accountName = fmt.Sprintf("userorg%d", time.Now().Unix())
	accountPwd = "pwd"
	accountEmail = "amp@axway.com"
	accountType="organization"

	//SignUp      Billing
	signupAnswer, signUpErr :=accountClient.SignUp(ctx, &account.SignUpRequest {
		Name: accountName,
		Password: accountPwd,
		Email: accountEmail,
		AccountType: accountType,
	})
	if signUpErr!=nil {
		t.Fatal(signUpErr)
	}
	if signupAnswer.SessionKey=="" {
		t.Fatal("invalide sessionKey")
	}

	code :="bad"
	//Verify with a bad code
	_, verifyErr1 := accountClient.Verify(ctx, &account.VerificationRequest {
		Name: accountName,
		Password: accountPwd,
		Code: code,
	})
	if verifyErr1==nil {
		t.Fatal("Should have error executing Verify with a bad code")
	}

	//a way to get the code...
	code= "ok"
	_, verifyErr2 := accountClient.Verify(ctx,  &account.VerificationRequest {
		Name: accountName,
		Password: accountPwd,
		Code: code,
	})
	if verifyErr2!=nil {
		t.Fatal("invalide verify code")
	}
}

func TestOrgShouldLogInResetPwdChangeItLogInAgain(t *testing.T) {
	//Login with bad pwd
	_, badLoginErr :=accountClient.Login(ctx, &account.LogInRequest {
		Name: accountName,
		Password: "badPwd",
	})
	if badLoginErr==nil {
		t.Fatal("Shoud have error, login with a bad pwd")
	}

	//Login with the right pwd
	loginAnswer1, loginErr :=accountClient.Login(ctx, &account.LogInRequest {
		Name: accountName,
		Password: accountPwd,
	})
	if loginErr!=nil {
		t.Fatal(loginErr)
	}
	if loginAnswer1.SessionKey=="" {
		t.Fatal("invalide sessionKey")
	}
	//Reset password
	_, pwdResetErr :=accountClient.PasswordReset(ctx, &account.PasswordResetRequest {
		Name: accountName,
		Email: accountEmail,
	})
	if pwdResetErr!=nil {
		t.Fatal(pwdResetErr)
	}

	//Change password
	req:= &account.PasswordChangeRequest {
		Name: accountName,
		ExistingPassword: accountPwd,
	}
	accountPwd="newPwd"
	req.NewPassword=accountPwd
	_, pwdChangeErr :=accountClient.PasswordChange(ctx, req)
	if pwdChangeErr!=nil {
		t.Fatal(pwdChangeErr)
	}

	//Login again using new password
	loginAnswer2, loginErr :=accountClient.Login(ctx, &account.LogInRequest {
		Name: accountName,
		Password: accountPwd,
	})
	if loginErr!=nil {
		t.Fatal(loginErr)
	}
	if loginAnswer2.SessionKey==nil {
		t.Fatal("invalide sessionKey")
	}
}

func TestOrgShouldEditAccount(t *testing.T) {
	//Edit email account
	accountEmail="ampnew@axway.com"
	_, editErr1:= accountClient.Edit(ctx, &account.EditAccountRequest{
		Name: accountName,
		Email: accountEmail,
	})
	if editErr1!=nil {
		t.Fatal(editErr1)
	}

	//Shouldn't be able to log without eamil verified
	_, loginErr :=accountClient.Login(ctx, &account.LogInRequest {
		Name: accountName,
		Password: accountPwd,
	})
	if loginErr==nil {
		t.Fatal("Shouldn't be able to log without email verified")
	}
	//Verify email
	//a way to get the code...
	code:= "ok"
	_, verifyErr := accountClient.Verify(ctx,  &account.VerificationRequest {
		Name: accountName,
		Password: accountPwd,
		Code: code,
	})
	if verifyErr!=nil {
		t.Fatal("invalide verify code")
	}
	settingMap:=make(map[string]string)
	settingMap["testkey"]="testval"
	//edit biiling and settings
	_, editErr2 := accountClient.Edit(ctx, &account.EditAccountRequest  {
		Name: accountName,
		Billing: &account.Billing{
			Contact: "Axway",
			AddressLine_1: "Tour W",
			AddressLine_2 : "esplanade boieldieu",
			City: "Paris - La Defense",
			State: "Paris",
			ZipCode: "92030",
			PhoneNumber: "33147172222",
			FaxNumber: "Does fax still exist??",
			Email: "payerMan@axway.com",
			CardType: "visa",
			CardNumber: "1234567890",
		},
		Settings: &account.Settings{
			param: settingMap,
		},
	})
	if editErr2!=nil {
		t.Fatal(editErr2)
	}
}

func TestOrgShouldListAndGetAccountDetails(t *testing.T) {
	//List existing accounts
	list, listErr := accountClient.List(ctx, &account.ListAccountRequest{})
	if listErr!=nil {
		t.Fatal(listErr)
	}
	found:=false
	for _, name := range list.Accounts {
		if name==accountName {
			found = true
			break
		}
		if !found {
			t.Fatal("List error, juste created account not found")
		}
		ret, getDetailsErr:= accountClient.GetDetails(ctx, account.GetAccountDetailsRequest{
			Name: accountName,
		})
		if getDetailsErr != nil {
			t.Fatal(getDetailsErr)
		}
		if ret.Account.Name != accountName {
			t.Fatalf("Account.Name should be: %s\n", accountName)
		}
		if ret.Account.Email != accountEmail {
			t.Fatalf("Account.Email should be: %s\n", accountEmail)
		}
		if !ret.Account.EmailVerified {
			t.Fatalf("Account.EmailVerified should be: true\n")
		}
		if ret.Account.AccountType != accountType {
			t.Fatalf("Account.AccountType should be: %s\n", accountType)
		}
		if ret.Billing.Contact != "Axway" {
			t.Fatalf("Billing.Contact should be: Axway\n")
		}
		if ret.Billing.AddressLine_1 != "Tour W" {
			t.Fatalf("Billing.AddressLine1 should be: Tour W\n")
		}
		if ret.Billing.AddressLine_2 != "esplanade boieldieu" {
			t.Fatalf("Billing.AddressLine2 should be: esplanade boieldieu\n")
		}
		if ret.Billing.City != "Paris - La Defense" {
			t.Fatalf("Billing.City should be: Paris - La Defense\n")
		}
		if ret.Billing.State != "Paris" {
			t.Fatalf("Billing.State should be: Paris\n")
		}
		if ret.Billing.ZipCode != "92030" {
			t.Fatalf("Billing.ZipCode should be: 92030\n")
		}
		if ret.Billing.PhoneNumber != "33147172222" {
			t.Fatalf("Billing.PhoneNumber should be: 33147172222\n")
		}
		if ret.Billing.FaxNumber != "Does fax still exist??" {
			t.Fatalf("Billing.FaxNumber should be: Does fax still exist??\n")
		}
		if ret.Billing.CardType != "visa" {
			t.Fatalf("Billing.CardType should be: visa\n")
		}
		if ret.Billing.CardNumber != "1234567890" {
			t.Fatalf("Billing.CardNumber should be: 1234567890\n")
		}
		if val, ok := ret.Settings.Param["testkey"]; !ok {
			t.Fatalf("Setting['testKey'] doesn't exist")
		} else if val!="testval" {
			t.Fatalf("Settings['testkey] should be: testval\n")
		}
	}
}

func TestOrgShouldRemoveAccount(t *testing.T) {
	//Remove account
	_, removeErr:=accountClient.Delete(ctx, &account.DeleteAccountRequest{
		Name: accountName,
	})
	if removeErr != nil {
		t.Fatalf(removeErr)
	}
	//try to login again
	//Shouldn't be able to log without eamil verified
	_, loginErr :=accountClient.Login(ctx, &account.LogInRequest {
		Name: accountName,
		Password: accountPwd,
	})
	if loginErr==nil {
		t.Fatal("Shouldn't be able to login after accound deleted")
	}
}

