package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"

	telegramapi "github.com/P1zz/telegramBot-cli/internal/telegram"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configFile        string
	configErr         error
	version                               = "dev"
	commit                                = "unknown"
	buildDate                             = "unknown"
	newTelegramClient telegramapi.Factory = telegramapi.New
)

var rootCmd = &cobra.Command{
	Use:           "telegram-cli",
	Short:         "Use the Telegram Bot API from the command line",
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Flags().Changed("chatId") {
			value, _ := cmd.Flags().GetInt64("chatId")
			viper.Set("chatId", value)
		}
		return configErr
	},
}

func Execute(ctx context.Context) error {
	rootCmd.SetContext(ctx)
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "C", "", "path to a TOML config file")
	rootCmd.PersistentFlags().StringP("token", "t", "", "Telegram bot token")
	rootCmd.PersistentFlags().String("token-file", "", "read the Telegram bot token from a file, or - for stdin")
	rootCmd.PersistentFlags().Int64P("chat-id", "c", 0, "Telegram chat ID")
	rootCmd.PersistentFlags().Int64("chatId", 0, "deprecated alias for --chat-id")
	_ = rootCmd.PersistentFlags().MarkDeprecated("chatId", "use --chat-id")

	cobra.CheckErr(viper.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token")))
	cobra.CheckErr(viper.BindPFlag("chatId", rootCmd.PersistentFlags().Lookup("chat-id")))
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if info, ok := debug.ReadBuildInfo(); ok && version == "dev" && info.Main.Version != "(devel)" {
				version = info.Main.Version
			}
			fmt.Fprintf(cmd.OutOrStdout(), "telegram-cli %s\ncommit: %s\nbuilt: %s\ngo: %s\n", version, commit, buildDate, runtime.Version())
		},
	})
}

func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
		if dir, err := os.UserConfigDir(); err == nil {
			viper.AddConfigPath(filepath.Join(dir, "telegram-cli"))
		}
	}
	viper.SetConfigType("toml")
	viper.SetEnvPrefix("TELEGRAM_CLI")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, missing := err.(viper.ConfigFileNotFoundError); missing && configFile == "" {
			return
		}
		if errors.Is(err, os.ErrNotExist) && configFile == "" {
			return
		}
		configErr = fmt.Errorf("read config %q: %w", configFile, err)
	}
}

func resolvedToken() (string, error) {
	if path, _ := rootCmd.PersistentFlags().GetString("token-file"); path != "" {
		if path == "-" {
			data, err := io.ReadAll(os.Stdin)
			if err != nil {
				return "", fmt.Errorf("read token from stdin: %w", err)
			}
			return strings.TrimSpace(string(data)), nil
		}
		info, err := os.Stat(path)
		if err != nil {
			return "", fmt.Errorf("read token file: %w", err)
		}
		if info.Mode().Perm()&0o077 != 0 {
			return "", fmt.Errorf("token file %q must not be accessible by group or others", path)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return "", fmt.Errorf("read token file: %w", err)
		}
		return strings.TrimSpace(string(data)), nil
	}
	if used := viper.ConfigFileUsed(); used != "" && viper.GetString("token") != "" {
		if info, err := os.Stat(used); err == nil && info.Mode().Perm()&0o077 != 0 {
			return "", fmt.Errorf("config %q contains a token and must not be accessible by group or others", used)
		}
	}
	return viper.GetString("token"), nil
}

func safeError(action, token string, err error) error {
	message := err.Error()
	if token != "" {
		message = strings.ReplaceAll(message, token, "[REDACTED]")
	}
	return fmt.Errorf("%s: %s", action, message)
}
