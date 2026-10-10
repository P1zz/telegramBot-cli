package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configFile string
	configErr  error
	version    = "dev"
)

var rootCmd = &cobra.Command{
	Use:           "telegram-cli",
	Short:         "Use the Telegram Bot API from the command line",
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return configErr
	},
}

func Execute(ctx context.Context) error {
	rootCmd.SetContext(ctx)
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "C", "config.toml", "path to a TOML config file")
	rootCmd.PersistentFlags().StringP("token", "t", "", "Telegram bot token")
	rootCmd.PersistentFlags().Int64P("chat-id", "c", 0, "Telegram chat ID")

	cobra.CheckErr(viper.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token")))
	cobra.CheckErr(viper.BindPFlag("chatId", rootCmd.PersistentFlags().Lookup("chat-id")))
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(cmd.OutOrStdout(), version)
		},
	})
}

func initConfig() {
	viper.SetConfigFile(configFile)
	viper.SetConfigType("toml")
	viper.SetEnvPrefix("TELEGRAM_CLI")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if errors.Is(err, os.ErrNotExist) && !rootCmd.PersistentFlags().Changed("config") {
			return
		}
		configErr = fmt.Errorf("read config %q: %w", configFile, err)
	}
}
