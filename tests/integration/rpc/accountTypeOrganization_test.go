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
//	account1Name   string
//	account1Email  string
//	account1Pwd    string
//	account2Name   string
//	account2Email  string
//	account2Pwd    string
//	accountOrgType = schema.AccountType_ORGANIZATION
//)
//
////Create two organization accounts
//func TestAccountOrgShouldSignUpAndVerify(t *testing.T) {
//	accountName = fmt.Sprintf("user%d", time.Now().Unix())
//	accountPwd = "pwd"
//	accountEmail = "amp@axway.com"
//
//	account1Name = fmt.Sprintf("accountOrg%d", time.Now().Unix())
//	account1Pwd = "pwd1"
//	account1Email = "amp@axway.com"
//
//	account2Name = fmt.Sprintf("accountOrg%d", time.Now().Unix())
//	account2Pwd = "pwd2"
//	account2Email = "amp@axway.com"
//
//	//SignUp account 1
//	signUpReply1, signUpErr1 := accountClient.SignUp(ctx, &account.SignUpRequest{
//		UserName:    account1Name,
//		Password:    account1Pwd,
//		Email:       account1Email,
//		AccountType: accountOrgType,
//	})
//	if signUpErr1 != nil {
//		t.Fatal(signUpErr1)
//	}
//	//a way to get the code...
//	_, verifyErr1 := accountClient.Verify(ctx, &account.VerificationRequest{
//		Token: signUpReply1.Token,
//	})
//	if verifyErr1 != nil {
//		t.Fatal(verifyErr1)
//	}
//	//SignUp Account 2
//	signUpReply2, signUpErr2 := accountClient.SignUp(ctx, &account.SignUpRequest{
//		UserName:    account2Name,
//		Password:    account1Pwd,
//		Email:       account2Email,
//		AccountType: accountOrgType,
//	})
//	if signUpErr2 != nil {
//		t.Fatal(signUpErr2)
//	}
//	//a way to get the code...
//	_, verifyErr2 := accountClient.Verify(ctx, &account.VerificationRequest{
//		Token: signUpReply2.Token,
//	})
//	if verifyErr2 != nil {
//		t.Fatal(verifyErr2)
//	}
//	//SignUp user
//	signUpReply, signUpErr := accountClient.SignUp(ctx, &account.SignUpRequest{
//		UserName:    accountName,
//		Password:    accountPwd,
//		Email:       accountEmail,
//		AccountType: accountUserType,
//	})
//	if signUpErr != nil {
//		t.Fatal(signUpErr)
//	}
//	//a way to get the code...
//	_, verifyErr := accountClient.Verify(ctx, &account.VerificationRequest{
//		Token: signUpReply.Token,
//	})
//	if verifyErr != nil {
//		t.Fatal(verifyErr)
//	}
//}
//
//func TestAccountOrgShouldLogInResetPwdChangeItLogInAgain(t *testing.T) {
//	//Login with bad pwd
//	_, badLoginErr := accountClient.Login(ctx, &account.LogInRequest{
//		UserName: accountName,
//		Password: "badPwd",
//	})
//	if badLoginErr == nil {
//		t.Fatal("Shoud have error, login with a bad pwd")
//	}
//}
//
//func TestAccountOrgShouldEditAccount(t *testing.T) {
//	//create settings object
//	settingMap := make(map[string]string)
//	settingMap["testkey"] = "testval"
//	//create billing object
//	billing := &account.Billing{
//		Contact:       "Axway",
//		AddressLine_1: "Tour W",
//		AddressLine_2: "esplanade boieldieu",
//		City:          "Paris - La Defense",
//		State:         "Paris",
//		ZipCode:       "92030",
//		PhoneNumber:   "33147172222",
//		FaxNumber:     "Does fax still exist??",
//		Email:         "payerMan@axway.com",
//		CardType:      "visa",
//		CardNumber:    "1234567890",
//	}
//	//Update account1
//	_, editErr1 := accountClient.Edit(ctx, &account.EditAccountRequest{
//		UserName: account1Name,
//		Billing:  billing,
//		Settings: &account.Settings{
//			Param: settingMap,
//		},
//	})
//	if editErr1 != nil {
//		t.Fatal(editErr1)
//	}
//	//update account 2
//	_, editErr2 := accountClient.Edit(ctx, &account.EditAccountRequest{
//		UserName: account2Name,
//		Billing:  billing,
//		Settings: &account.Settings{
//			Param: settingMap,
//		},
//	})
//	if editErr2 != nil {
//		t.Fatal(editErr2)
//	}
//}
//
//func TestAccountOrgShouldListAndGetAccountDetails(t *testing.T) {
//	//List existing accounts
//	list, listErr := accountClient.List(ctx, &account.ListAccountRequest{})
//	if listErr != nil {
//		t.Fatal(listErr)
//	}
//	found := 0
//	for _, acc := range list.Accounts {
//		if acc.UserName == account1Name || acc.UserName == account2Name {
//			found++
//		}
//	}
//	if found != 2 {
//		t.Fatal("List error, juste created accounts are not found")
//	}
//	ret, getDetailsErr := accountClient.GetDetails(ctx, &account.GetAccountDetailsRequest{
//		UserName: account1Name,
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
//	if ret.Account.AccountType != accountOrgType {
//		t.Fatalf("Account.AccountType should be: %s\n", accountOrgType)
//	}
//	if ret.Billing.Contact != "Axway" {
//		t.Fatalf("Billing.Contact should be: Axway\n")
//	}
//	if ret.Billing.AddressLine_1 != "Tour W" {
//		t.Fatalf("Billing.AddressLine1 should be: Tour W\n")
//	}
//	if ret.Billing.AddressLine_2 != "esplanade boieldieu" {
//		t.Fatalf("Billing.AddressLine2 should be: esplanade boieldieu\n")
//	}
//	if ret.Billing.City != "Paris - La Defense" {
//		t.Fatalf("Billing.City should be: Paris - La Defense\n")
//	}
//	if ret.Billing.State != "Paris" {
//		t.Fatalf("Billing.State should be: Paris\n")
//	}
//	if ret.Billing.ZipCode != "92030" {
//		t.Fatalf("Billing.ZipCode should be: 92030\n")
//	}
//	if ret.Billing.PhoneNumber != "33147172222" {
//		t.Fatalf("Billing.PhoneNumber should be: 33147172222\n")
//	}
//	if ret.Billing.FaxNumber != "Does fax still exist??" {
//		t.Fatalf("Billing.FaxNumber should be: Does fax still exist??\n")
//	}
//	if ret.Billing.CardType != "visa" {
//		t.Fatalf("Billing.CardType should be: visa\n")
//	}
//	if ret.Billing.CardNumber != "1234567890" {
//		t.Fatalf("Billing.CardNumber should be: 1234567890\n")
//	}
//	if val, ok := ret.Settings.Param["testkey"]; !ok {
//		t.Fatalf("Setting['testKey'] doesn't exist")
//	} else if val != "testval" {
//		t.Fatalf("Settings['testkey] should be: testval\n")
//	}
//}
//
//func TestAccountOrgAddTeam(t *testing.T) {
//	//create a team
//	_, createTeamErr := accountClient.CreateTeam(ctx, &account.CreateTeamRequest{
//		Organization: account1Name,
//		Name:         "myTeam",
//		Description:  "My team for test",
//	})
//	if createTeamErr != nil {
//		t.Fatal(createTeamErr)
//	}
//	//List team and verify my team is in it
//	teamList, teamListErr := accountClient.ListTeam(ctx, &account.ListTeamRequest{
//		Organization: account1Name,
//	})
//	if teamListErr != nil {
//		t.Fatal(teamListErr)
//	}
//	found := false
//	for _, teamName := range teamList.TeamNames {
//		if teamName == "myTeam" {
//			found = true
//			break
//		}
//	}
//	if !found {
//		t.Fatal("the just created team 'myTeam' is not found in team list")
//	}
//}
//
//func TestAccountOrgEditTeamAndList(t *testing.T) {
//	//create a team
//	_, editTeamErr := accountClient.EditTeam(ctx, &account.EditTeamRequest{
//		Organization:   account1Name,
//		Name:           "myTeam",
//		NewName:        "myNewTeam",
//		NewDescription: "My new team for test",
//	})
//	if editTeamErr != nil {
//		t.Fatal(editTeamErr)
//	}
//	//List team and verify my team is in it
//	teamList, teamListErr := accountClient.ListTeam(ctx, &account.ListTeamRequest{
//		Organization: account1Name,
//	})
//	if teamListErr != nil {
//		t.Fatal(teamListErr)
//	}
//	found := false
//	for _, teamName := range teamList.TeamNames {
//		if teamName == "myNewTeam" {
//			found = true
//			break
//		}
//	}
//	if !found {
//		t.Fatal("the just created team 'myTeam' is not found in team list")
//	}
//}
//
//func TestAccountOrgMoveTeamToAnotherOrganization(t *testing.T) {
//	//create a team
//	_, editTeamErr := accountClient.EditTeam(ctx, &account.EditTeamRequest{
//		Organization:    account1Name,
//		Name:            "myTeam",
//		NewOrganization: account2Name,
//	})
//	if editTeamErr != nil {
//		t.Fatal(editTeamErr)
//	}
//	//List team of account2 and verify my team is in it
//	teamList1, teamListErr1 := accountClient.ListTeam(ctx, &account.ListTeamRequest{
//		Organization: account2Name,
//	})
//	if teamListErr1 != nil {
//		t.Fatal(teamListErr1)
//	}
//	found := false
//	for _, teamName := range teamList1.TeamNames {
//		if teamName == "myNewTeam" {
//			found = true
//			break
//		}
//	}
//	if !found {
//		t.Fatal("the team 'myTeam' has not been move to another organization")
//	}
//	//List team of account1 and verify my team is not there anymore
//	teamList2, teamListErr2 := accountClient.ListTeam(ctx, &account.ListTeamRequest{
//		Organization: account1Name,
//	})
//	if teamListErr2 != nil {
//		t.Fatal(teamListErr2)
//	}
//	found = false
//	for _, teamName := range teamList2.TeamNames {
//		if teamName == "myNewTeam" {
//			found = true
//			break
//		}
//	}
//	if found {
//		t.Fatal("the team 'myTeam' should have been moved, but it stays in its original organization")
//	}
//}
//
//func TestAccountOrgShouldAddUserToTeam(t *testing.T) {
//	//Try to add a bad user in the team
//	_, addUserErr1 := accountClient.AddTeamMemberships(ctx, &account.AddTeamMembershipsRequest{
//		Organization: account2Name,
//		Name:         "myNewTeam",
//		Members:      []string{"badUser"},
//	})
//	if addUserErr1 == nil {
//		t.Fatal("Should have error when added the user who doesn't exist in a team")
//	}
//	//add a reguler user in a team
//	_, addUserErr2 := accountClient.AddTeamMemberships(ctx, &account.AddTeamMembershipsRequest{
//		Organization: account2Name,
//		Name:         "myNewTeam",
//		Members:      []string{accountName},
//	})
//	if addUserErr2 == nil {
//		t.Fatal(addUserErr2)
//	}
//}
//
//func TestAccountOrgShouldRemoveAccount(t *testing.T) {
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
