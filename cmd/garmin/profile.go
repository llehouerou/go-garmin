package main

import (
	"github.com/spf13/cobra"
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "User profile data (social, settings)",
}

var profileSocialCmd = &cobra.Command{
	Use:   "social",
	Short: "Get social profile",
	Args:  cobra.NoArgs,
	RunE:  runProfileSocial,
}

var profileSettingsCmd = &cobra.Command{
	Use:   "settings",
	Short: "Get profile settings",
	Args:  cobra.NoArgs,
	RunE:  runProfileSettings,
}

var profileUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Get user settings (detailed)",
	Args:  cobra.NoArgs,
	RunE:  runProfileUser,
}

func init() {
	profileCmd.AddCommand(profileSocialCmd)
	profileCmd.AddCommand(profileSettingsCmd)
	profileCmd.AddCommand(profileUserCmd)
}

func runProfileSocial(cmd *cobra.Command, _ []string) error {
	client, err := loadClient()
	if err != nil {
		return err
	}

	data, err := client.UserProfile.GetSocialProfile(cmd.Context())
	if err != nil {
		return err
	}

	return printJSON(data)
}

func runProfileSettings(cmd *cobra.Command, _ []string) error {
	client, err := loadClient()
	if err != nil {
		return err
	}

	data, err := client.UserProfile.GetProfileSettings(cmd.Context())
	if err != nil {
		return err
	}

	return printJSON(data)
}

func runProfileUser(cmd *cobra.Command, _ []string) error {
	client, err := loadClient()
	if err != nil {
		return err
	}

	data, err := client.UserProfile.GetUserSettings(cmd.Context())
	if err != nil {
		return err
	}

	return printJSON(data)
}
