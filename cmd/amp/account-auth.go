package main

import (
	"fmt"

	"github.com/appcelerator/amp/api/client"
	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/fatih/color"
	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"strings"
)

var (
	LoginCmd = &cobra.Command{
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

	AccountCmd = &cobra.Command{
		Use:   "account",
		Short: "Account operations",
		Long:  `The account command manages all account-related operations.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return AMP.Connect()
		},
	}

	signUpCmd = &cobra.Command{
		Use:   "signup",
		Short: "Signup for a new account",
		Long:  `The signup command creates a new account and sends a verification link to their registered email address.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return signUp(AMP)
		},
	}

	verifyCmd = &cobra.Command{
		Use:   "verify",
		Short: "Verify email using code",
		Long: `The verify command is used to verify the users account using the code sent to them via email.
This is used if the user cannot access the verification link sent.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return verify(AMP)
		},
	}

	pwdCmd = &cobra.Command{
		Use:   "password",
		Short: "Password operations",
		Long:  "The password command allows users to reset or update password.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return pwd(AMP, cmd, args)
		},
	}

	pwdResetCmd = &cobra.Command{
		Use:   "reset ACCOUNT-NAME EMAIL",
		Short: "Reset password",
		Long:  "The reset command allows users to reset password. A link to reset password will be sent to their registered email address.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return pwdReset(AMP, cmd, args)
		},
	}

	pwdChangeCmd = &cobra.Command{
		Use:   "change ACCOUNT-NAME EXISTING-PASSWORD NEW-PASSWORD CONFIRM-NEW-PASSWORD",
		Short: "Change password",
		Long:  "The change command allows users to update their existing password.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return pwdChange(AMP, cmd, args)
		},
	}
)

func init() {
	RootCmd.AddCommand(LoginCmd)
	RootCmd.AddCommand(AccountCmd)
	AccountCmd.AddCommand(signUpCmd)
	AccountCmd.AddCommand(verifyCmd)
	AccountCmd.AddCommand(pwdCmd)
	pwdCmd.AddCommand(pwdResetCmd)
	pwdCmd.AddCommand(pwdChangeCmd)
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
	fmt.Println("Welcome back, ", username)
	color.Unset()
	return nil
}

// signup signs up visitor for a new personal account.
// Sends a verification link to their email address.
func signUp(amp *client.AMP) error {
	fmt.Println("This will sign you up for a new personal AMP account.")
	username := getUserName()
	email, err := getEmailAddress()
	if err != nil {
		return fmt.Errorf("user error: %v", err)
	}
	request := &account.SignUpRequest{
		Name:  username,
		Email: email,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	_, err = client.SignUp(context.Background(), request)
	if err != nil {
		return fmt.Errorf("server error: %v", err)
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("Hi", username, "!, Please check your email to complete the signup process.")
	color.Unset()
	return nil
}

// verify gets the unique code sent to the visitor in the email verification, registered username and new password,
// validates the command line inputs and activates their account.
func verify(amp *client.AMP) (err error) {
	fmt.Println("This will verify your account and confirm your password.")
	code := getCode()
	username := getUserName()
	password, err := getPwd()
	if err != nil {
		return fmt.Errorf("user error: %v", err)
	}
	request := &account.VerificationRequest{
		Name:     username,
		Password: password,
		Code:     code,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	_, err = client.Verify(context.Background(), request)
	if err != nil {
		return fmt.Errorf("server error: %v", err)
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("Hi", username, "! Your account has now be activated.")
	color.Unset()
	return nil
}

func pwd(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
	color.Set(color.FgYellow, color.Bold)
	fmt.Println("Choose a command for password operation")
	fmt.Println("Use amp account password -h for help")
	color.Unset()
	return nil
}

// pwdReset validates the input command line arguments and resets the current password
// by invoking the corresponding rpc/storage method
func pwdReset(amp *client.AMP, cmd *cobra.Command, args []string) error {
	fmt.Println("This will send a password reset email to your email address.")
	username := getUserName()
	email, err := getEmailAddress()
	if err != nil {
		return fmt.Errorf("user error: %v", err)
	}
	request := &account.PasswordResetRequest{
		Name:  username,
		Email: email,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	_, err = client.PasswordReset(context.Background(), request)
	if err != nil {
		return fmt.Errorf("server error: %v", err)
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("Hi", username, "! Please check your email to complete the password reset process.")
	color.Unset()
	return nil
}

// pwdChange validates the input command line arguments and changes the current password
// by invoking the corresponding rpc/storage method
func pwdChange(amp *client.AMP, cmd *cobra.Command, args []string) error {
	fmt.Println("This will allow you to update your existing password.")
	username := getUserName()
	color.Set(color.FgYellow, color.Bold)
	fmt.Println("Enter your current password.")
	color.Unset()
	existingPwd, err := getPwd()
	if err != nil {
		return fmt.Errorf("user error: %v", err)
	}
	color.Set(color.FgYellow, color.Bold)
	fmt.Println("Enter new password.")
	color.Unset()
	newPwd, err := getPwd()
	if err != nil {
		return fmt.Errorf("user error: %v", err)
	}
	getConfirmNewPwd(newPwd)
	request := &account.PasswordChangeRequest{
		Name:             username,
		ExistingPassword: existingPwd,
		NewPassword:      newPwd,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	_, err = client.PasswordChange(context.Background(), request)
	if err != nil {
		return fmt.Errorf("server error : %v", err)
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("Hi ", username, "! Your recent password change has been successful.")
	color.Unset()
	return nil
}

func getUserName() (username string) {
	fmt.Print("username: ")
	color.Set(color.FgGreen, color.Bold)
	fmt.Scanln(&username)
	color.Unset()
	err := account.CheckUserName(username)
	if err != nil {
		color.Set(color.FgRed, color.Bold)
		fmt.Println("Username is mandatory. Try again!")
		color.Unset()
		fmt.Println("")
		return getUserName()
	}
	return
}

func getEmailAddress() (email string, err error) {
	fmt.Print("email: ")
	color.Set(color.FgGreen, color.Bold)
	fmt.Scanln(&email)
	color.Unset()
	email, err = account.CheckEmailAddress(email)
	if err != nil {
		color.Set(color.FgRed, color.Bold)
		fmt.Println("Email in incorrect format. Try again!")
		color.Unset()
		fmt.Println("")
		return getEmailAddress()
	}
	return
}

func getCode() (code string) {
	fmt.Print("code: ")
	color.Set(color.FgGreen, color.Bold)
	fmt.Scanln(&code)
	color.Unset()
	err := account.CheckVerificationCodeFormat(code)
	if err != nil {
		color.Set(color.FgRed, color.Bold)
		fmt.Println("Code is invalid. Code must be 8 characters long. Try again!")
		color.Unset()
		fmt.Println("")
		return getCode()
	}
	return
}

func getPwd() (password string, err error) {
	fmt.Print("password: ")
	pw, err := gopass.GetPasswd()
	if err != nil {
		color.Set(color.FgRed, color.Bold)
		fmt.Println("Password is mandatory. Try again!")
		color.Unset()
		fmt.Println("")
		return getPwd()
	}
	password = string(pw)
	err = account.CheckPassword(password)
	if err != nil {
		fmt.Println(err)
	}

	err = account.CheckPasswordStrength(password)
	if err != nil {
		if strings.Contains(err.Error(), "password too weak") {
			color.Set(color.FgRed, color.Bold)
			fmt.Println("Password entered is too weak. Password must be at least 8 characters long. Try again!")
			color.Unset()
			fmt.Println("")
			return getPwd()
		} else {
			return
		}
	}
	return
}

func getConfirmNewPwd(newPwd string) {
	color.Set(color.FgYellow, color.Bold)
	fmt.Println("Enter Password again for confirmation.")
	color.Unset()
	confirmNewPwd, _ := getPwd()
	if confirmNewPwd != newPwd {
		color.Set(color.FgRed, color.Bold)
		fmt.Println("Password mismatch. Try again!")
		color.Unset()
		fmt.Println("")
		getConfirmNewPwd(newPwd)
	} else {
		return
	}
	return
}
