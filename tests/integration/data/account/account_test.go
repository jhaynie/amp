package account

//
//import (
//	"github.com/appcelerator/amp/config"
//	"github.com/appcelerator/amp/data/account"
//	"github.com/appcelerator/amp/data/schema"
//	"github.com/appcelerator/amp/data/storage"
//	"github.com/appcelerator/amp/data/storage/etcd"
//	"golang.org/x/net/context"
//	"log"
//	"os"
//	"strings"
//	"testing"
//	"time"
//)
//
//const (
//	defTimeout = 5 * time.Second
//)
//
//var (
//	acct         account.Interface
//	testAcct     schema.Account
//	testTeam     schema.Team
//	testResource schema.Resource
//)
//
//func newContext() context.Context {
//	ctx, _ := context.WithTimeout(context.Background(), defTimeout)
//	return ctx
//}
//
//func initData() {
//	testAcct = schema.Account{
//		Id:         "",
//		Name:       "axway",
//		Type:       schema.AccountType_ORGANIZATION,
//		Email:      "testowner@axway.com",
//		IsVerified: false,
//	}
//	testTeam = schema.Team{
//		Id:   "",
//		Name: "Falcons",
//		Desc: "The Falcons",
//	}
//	testResource = schema.Resource{
//		Id:           "",
//		Name:         "test-stack",
//		TeamId:       "",
//		OrgAccountId: "",
//		Type:         schema.ResourceType_STACK,
//	}
//	// delete any data from previous unit tests OR from failed executions
//	// prefix is amp-test
//	deleteAll()
//}
//
//var store storage.Interface
//
//func TestMain(m *testing.M) {
//	log.SetOutput(os.Stdout)
//	log.SetFlags(log.Lshortfile)
//	log.SetPrefix("test: ")
//
//	etcdEndpoints := []string{amp.EtcdDefaultEndpoint}
//	log.Printf("connecting to etcd at %s", strings.Join(etcdEndpoints, ","))
//	store = etcd.New(etcdEndpoints, "amp-test")
//	if err := store.Connect(defTimeout); err != nil {
//		log.Panicf("Unable to connect to etcd on: %s\n%v", etcdEndpoints, err)
//	}
//	initData()
//	log.Printf("connected to etcd at %v", strings.Join(store.Endpoints(), ","))
//	acct = account.NewStore(store)
//	os.Exit(m.Run())
//}
//
//func deleteAll() {
//	// Recursively delete all the test data, so each test starts clean
//	store.Delete(context.Background(), "/accounts", true, nil)
//}
//func deleteAccounts() {
//	store.Delete(context.Background(), "/accounts/team", true, nil)
//
//}
//func addAccounts() {
//	acct.AddAccount(newContext(), &testAcct)
//	testAcct.Name = "axway2"
//	// Add Org Account
//	acct.AddAccount(newContext(), &testAcct)
//
//	// Add a second Org Account
//	testAcct.Name = "atomiq"
//	testAcct.Id = ""
//	acct.AddAccount(newContext(), &testAcct)
//
//	//Add first User account
//	testAcct.Id = ""
//	testAcct.Name = "theuser"
//	testAcct.Type = schema.AccountType_USER
//	acct.AddAccount(newContext(), &testAcct)
//
//}
//
//func addTeam() error {
//	//Team Depends on Presence of an Account
//	addAccounts()
//	a, err := acct.GetAccount(newContext(), "axway")
//	if err == nil {
//		testTeam.OrgAccountId = a.Id
//		_, err = acct.AddTeam(newContext(), &testTeam)
//	}
//	return err
//}
//
//func TestAddAccount(t *testing.T) {
//	//Store the account
//	_, err := acct.AddAccount(newContext(), &testAcct)
//	if err != nil {
//		t.Error(err)
//	}
//
//	//Test that the account was stored
//	a, _ := acct.GetAccount(newContext(), testAcct.Name)
//	if a.Id == "" {
//		t.Errorf("Failed to retrieve Account from path suffix %s", testAcct.Name)
//	}
//}
//func TestAddDuplicateAccount(t *testing.T) {
//	acct.AddAccount(newContext(), &testAcct)
//	_, err := acct.AddAccount(newContext(), &testAcct)
//
//	// This should result in a duplicate error
//	if err == nil || !strings.Contains(err.Error(), "already exists") {
//		t.Errorf("Expected \"already exists\" Errorv\n")
//	}
//}
//
//func TestAddTeam(t *testing.T) {
//	// Requires Account->Team Relationship
//	addAccounts()
//	err := addTeam()
//	if err != nil {
//		t.Error(err)
//	}
//	tm, _ := acct.GetTeam(newContext(), "Falcons")
//	if tm.Id == "" {
//		t.Errorf("Failed to retrieve Team from path suffix %s", "Falcons")
//	}
//}
//
//func TestAddTeamMember(t *testing.T) {
//	// Add Team Member depends on presence of a Team
//	addTeam()
//
//	// Retrieve Team Record
//	tm, err := acct.GetTeam(newContext(), "Falcons")
//	if err != nil {
//		t.Error(err)
//	}
//	u, err := acct.GetAccount(newContext(), "theuser")
//
//	if err != nil {
//		t.Error(err)
//	}
//
//	//Store TeamMember Record
//	mem := &schema.TeamMember{}
//	mem.UserAccountId = u.Id
//	mem.TeamId = tm.Id
//	_, err = acct.AddTeamMember(newContext(), mem)
//	if err != nil {
//		t.Error(err)
//	}
//
//	//Test the TeamMember was stored
//	m, _ := acct.GetTeamMember(newContext(), tm.Id, mem.Id)
//	if m.Id == "" {
//		t.Errorf("Failed to retrieve TeamMember from path suffix %s", mem.TeamId+"/"+mem.Id)
//	}
//}
//
//func TestAddDuplicateTeam(t *testing.T) {
//	acct.AddTeam(newContext(), &testTeam)
//	_, err := acct.AddTeam(newContext(), &testTeam)
//	if err == nil || !strings.Contains(err.Error(), "already exists") {
//		t.Errorf("Expected \"already exists\" Errorv\n")
//	}
//}
//
//func TestListAccount(t *testing.T) {
//	deleteAccounts()
//	addAccounts()
//	accList, err := acct.GetAccounts(newContext(), schema.AccountType_USER)
//	if err != nil {
//		log.Panicf("Unable to Fetch Account List: %v", err)
//	}
//	if len(accList) != 1 {
//		t.Errorf("expected %v, got %v", 1, len(accList))
//	}
//
//}
//
//func TestAddResource(t *testing.T) {
//
//	//Make sure dependencies exist
//	addTeam()
//
//	//Get Team Record
//	team, _ := acct.GetTeam(newContext(), "Falcons")
//
//	//Store the resource
//	testResource.TeamId = team.Id
//	_, err := acct.AddResource(newContext(), &testResource)
//	if err != nil {
//		log.Panicf("Unable to add Resource: %v", err)
//	}
//
//	//Test that the Resource was stored
//	res, _ := acct.GetResource(newContext(), testResource.Name)
//	if res.Id == "" {
//		t.Errorf("Failed to retrieve Resource from path suffix %s", testResource.Name)
//	}
//
//}
//
//func addResource() {
//
//	addTeam()
//	team, _ := acct.GetTeam(newContext(), "Falcons")
//	testResource.TeamId = team.Id
//	acct.AddResource(newContext(), &testResource)
//}
//
//func TestAddResourceSettings(t *testing.T) {
//
//	//Make Sure the dependencies exist
//	addResource()
//	resource, _ := acct.GetResource(newContext(), "test-stack")
//
//	// Store the Resource Settings Record
//	setting := &schema.ResourceSettings{}
//	setting.ResourceId = resource.Id
//	setting.Key = "foo"
//	setting.Value = "bar"
//	_, err := acct.AddResourceSettings(newContext(), setting)
//	if err != nil {
//		t.Errorf("Unable to add Resource: %v", err)
//	}
//
//	// Test that the ResourceSetting record was stored
//	resList, _ := acct.GetResourceSettings(newContext(), resource.Id)
//	if len(resList) != 1 {
//		t.Errorf("Expected 1 entry to be fetched from ResourceSettings got %d", len(resList))
//	}
//}
//func TestAddPermission(t *testing.T) {
//
//	addResource()
//	resource, _ := acct.GetResource(newContext(), "test-stack")
//	team, _ := acct.GetTeam(newContext(), "Falcons")
//	perm := schema.Permission{
//		ResourceId: resource.Id,
//		GrantType:  schema.GrantType_ALL,
//		TeamId:     team.Id,
//	}
//	acct.AddPermission(newContext(), &perm)
//
//}
//
//func TestDeleteResource(t *testing.T) {
//
//	addResource()
//	err := acct.DeleteResource(newContext(), "test-stack")
//	if err != nil {
//		t.Error(err)
//	}
//}
