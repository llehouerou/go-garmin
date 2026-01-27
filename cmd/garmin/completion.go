package main

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion script",
	Long: `Generate completion script for your shell.

To load completions:

Bash:
  $ source <(garmin completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ garmin completion bash > /etc/bash_completion.d/garmin
  # macOS:
  $ garmin completion bash > /usr/local/etc/bash_completion.d/garmin

Zsh:
  # If shell completion is not already enabled in your environment,
  # you will need to enable it. You can execute the following once:
  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ garmin completion zsh > "${fpath[1]}/_garmin"

  # You will need to start a new shell for this setup to take effect.

Fish:
  $ garmin completion fish | source

  # To load completions for each session, execute once:
  $ garmin completion fish > ~/.config/fish/completions/garmin.fish

PowerShell:
  PS> garmin completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> garmin completion powershell > garmin.ps1
  # and source this file from your PowerShell profile.
`,
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
	RunE: func(_ *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			return rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			return rootCmd.GenFishCompletion(os.Stdout, true)
		case "powershell":
			return rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
		}
		return nil
	},
}
