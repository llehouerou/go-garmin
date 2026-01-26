package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"

	"github.com/llehouerou/go-garmin"
)

func loginCmd(_ []string) {
	// Check if already logged in
	if client, err := loadClient(); err == nil {
		_ = client
		fmt.Fprintln(os.Stderr, "Already logged in. Use 'garmin logout' first.")
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Fprint(os.Stderr, "Email: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	fmt.Fprint(os.Stderr, "Password: ")
	passwordBytes, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		printError(fmt.Errorf("failed to read password: %w", err))
		os.Exit(1)
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
		printError(err)
		os.Exit(1)
	}

	if err := saveClient(client); err != nil {
		printError(fmt.Errorf("login succeeded but failed to save session: %w", err))
		os.Exit(1)
	}

	fmt.Fprintln(os.Stderr, "Login successful.")
}

func logoutCmd(_ []string) {
	if err := removeSession(); err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintln(os.Stderr, "Not logged in.")
			return
		}
		printError(err)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stderr, "Logged out.")
}

func printError(err error) {
	_ = json.NewEncoder(os.Stderr).Encode(map[string]string{"error": err.Error()})
}
