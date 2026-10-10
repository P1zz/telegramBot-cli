package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

const configTemplate = `# Prefer TELEGRAM_CLI_TOKEN or --token-file over storing a token here.
chatId = 0

[send]
fileTimeout = 60
output = "text"

[receive]
counter = 0
timeout = 0
sync = false
output = "jsonl"
`

var configCmd = &cobra.Command{Use: "config", Short: "Manage CLI configuration"}

var configInitCmd = &cobra.Command{
	Use: "init", Short: "Create a documented configuration file", Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		path, _ := cmd.Flags().GetString("output")
		force, _ := cmd.Flags().GetBool("force")
		if path == "" {
			dir, err := os.UserConfigDir()
			if err != nil {
				return err
			}
			path = filepath.Join(dir, "telegram-cli", "config.toml")
		}
		if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
			return fmt.Errorf("create config directory: %w", err)
		}
		flags := os.O_WRONLY | os.O_CREATE | os.O_EXCL
		if force {
			flags = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
		}
		file, err := os.OpenFile(path, flags, 0o600)
		if err != nil {
			return fmt.Errorf("create config: %w", err)
		}
		if _, err = file.WriteString(configTemplate); err != nil {
			file.Close()
			return fmt.Errorf("write config: %w", err)
		}
		if err = file.Close(); err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), path)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configInitCmd)
	configInitCmd.Flags().StringP("output", "o", "", "configuration destination")
	configInitCmd.Flags().Bool("force", false, "overwrite an existing file")
}
