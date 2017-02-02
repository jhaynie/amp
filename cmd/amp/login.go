package main

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/appcelerator/amp/api/client"
	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var	LoginCmd = &cobra.Command{
	Use:   "login",
	Short: "log in to amp",
	Long:  `The login command logs the user into existing account`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := AMP.Connect()
		if err != nil {
			return err
		}
		return login(AMP)
	},
}

func init() {
	RootCmd.AddCommand(LoginCmd)
}

// login gets the username and password, validates the command line inputs
// and logs the user into their account
func login(amp *client.AMP) (err error) {
	fmt.Println("This will login an existing personal AMP account.")
	username := getUserName()
	password, err := getPwd()
	if err != nil {
		return fmt.Errorf("user error: %v", err)
	}
	request := &account.LogInRequest{
		Name:     username,
		Password: password,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	_, err = client.Login(context.Background(), request)
	if err != nil {
		return fmt.Errorf("server error: %v", err.Error())
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("Welcome back,", username)
	color.Unset()
	return nil
}
