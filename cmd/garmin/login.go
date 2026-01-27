package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/llehouerou/go-garmin"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with Garmin Connect",
	Long:  "Interactively authenticate with Garmin Connect using email and password.",
	Args:  cobra.NoArgs,
	RunE:  runLogin,
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Remove saved session",
	Long:  "Remove the saved Garmin Connect session from disk.",
	Args:  cobra.NoArgs,
	RunE:  runLogout,
}

func runLogin(_ *cobra.Command, _ []string) error {
	// Check if already logged in
	if client, err := loadClient(); err == nil {
		_ = client
		return errors.New("already logged in, use 'garmin logout' first")
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Fprint(os.Stderr, "Email: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	fmt.Fprint(os.Stderr, "Password: ")
	passwordBytes, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		return fmt.Errorf("failed to read password: %w", err)
	}
	fmt.Fprintln(os.Stderr) // newline after password
	password := string(passwordBytes)

	client := garmin.New(garmin.Options{
		MFAHandler: func() (string, error) {
			fmt.Fprint(os.Stderr, "MFA Code: ")
			code, _ := reader.ReadString('\n')
			return strings.TrimSpace(code), nil
		},
	})

	ctx := context.Background()
	if err := client.Login(ctx, email, password); err != nil {
		return err
	}

	if err := saveClient(client); err != nil {
		return fmt.Errorf("login succeeded but failed to save session: %w", err)
	}

	fmt.Fprintln(os.Stderr, "Login successful.")
	return nil
}

func runLogout(_ *cobra.Command, _ []string) error {
	if err := removeSession(); err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintln(os.Stderr, "Not logged in.")
			return nil
		}
		return err
	}
	fmt.Fprintln(os.Stderr, "Logged out.")
	return nil
}
