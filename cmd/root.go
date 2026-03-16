package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "telegram-cli",
	Short: "CLI tool for interact with telegram",
	Long:  `With this tool you can interact with telegram api as a bot`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("config", "C", "config.toml", "config file")
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	if viper.GetString("config") == "" {
		viper.SetConfigFile("config.toml")
	} else {
		viper.SetConfigFile(viper.GetString("config"))
	}

	viper.SetConfigType("toml")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(fmt.Errorf("fatal: %w", err))
		}
	}

	viper.AutomaticEnv()
}
