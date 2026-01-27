package main

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var devicesCmd = &cobra.Command{
	Use:   "devices",
	Short: "Device data (list, settings, messages)",
}

var devicesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List registered devices",
	Args:  cobra.NoArgs,
	RunE:  runDevicesList,
}

var devicesSettingsCmd = &cobra.Command{
	Use:   "settings <device-id>",
	Short: "Get settings for a specific device",
	Args:  cobra.ExactArgs(1),
	RunE:  runDevicesSettings,
}

var devicesMessagesCmd = &cobra.Command{
	Use:   "messages",
	Short: "Get device messages",
	Args:  cobra.NoArgs,
	RunE:  runDevicesMessages,
}

var devicesPrimaryCmd = &cobra.Command{
	Use:   "primary",
	Short: "Get primary training device info",
	Args:  cobra.NoArgs,
	RunE:  runDevicesPrimary,
}

func init() {
	devicesCmd.AddCommand(devicesListCmd)
	devicesCmd.AddCommand(devicesSettingsCmd)
	devicesCmd.AddCommand(devicesMessagesCmd)
	devicesCmd.AddCommand(devicesPrimaryCmd)
}

func runDevicesList(cmd *cobra.Command, _ []string) error {
	client, err := loadClient()
	if err != nil {
		return err
	}

	devices, err := client.Devices.GetDevices(cmd.Context())
	if err != nil {
		return err
	}

	return printJSON(devices)
}

func runDevicesSettings(cmd *cobra.Command, args []string) error {
	deviceID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid device ID: %s", args[0])
	}

	client, err := loadClient()
	if err != nil {
		return err
	}

	settings, err := client.Devices.GetSettings(cmd.Context(), deviceID)
	if err != nil {
		return err
	}

	return printJSON(settings)
}

func runDevicesMessages(cmd *cobra.Command, _ []string) error {
	client, err := loadClient()
	if err != nil {
		return err
	}

	messages, err := client.Devices.GetMessages(cmd.Context())
	if err != nil {
		return err
	}

	return printJSON(messages)
}

func runDevicesPrimary(cmd *cobra.Command, _ []string) error {
	client, err := loadClient()
	if err != nil {
		return err
	}

	info, err := client.Devices.GetPrimaryTrainingDevice(cmd.Context())
	if err != nil {
		return err
	}

	return printJSON(info)
}
