package cmd

import (
	"context"
	"fmt"

	"github.com/P1zz/telegramBot-cli/internal/validation"

	"github.com/go-telegram/bot"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type DeleteConfig struct {
	Token     string
	ChatID    int64
	MessageID int
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a message",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := loadDeleteConfig()
		if err := cfg.Validate(); err != nil {
			return err
		}
		return deleteMessage(cmd.Context(), cfg)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().IntP("message-id", "i", 0, "ID of the message to delete")
	bindFlag("delete.messageId", deleteCmd, "message-id")
}

func loadDeleteConfig() DeleteConfig {
	return DeleteConfig{
		Token: viper.GetString("token"), ChatID: viper.GetInt64("chatId"),
		MessageID: viper.GetInt("delete.messageId"),
	}
}

func (cfg DeleteConfig) Validate() error {
	if err := validation.Token(cfg.Token); err != nil {
		return err
	}
	if err := validation.RequiredChatID(cfg.ChatID); err != nil {
		return err
	}
	return validation.MessageID(cfg.MessageID)
}

func deleteMessage(ctx context.Context, cfg DeleteConfig) error {
	tgBot, err := bot.New(cfg.Token)
	if err != nil {
		return fmt.Errorf("create bot: %w", err)
	}
	deleted, err := tgBot.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID: cfg.ChatID, MessageID: cfg.MessageID,
	})
	if err != nil {
		return fmt.Errorf("delete message: %w", err)
	}
	if !deleted {
		return fmt.Errorf("Telegram did not delete message %d", cfg.MessageID)
	}
	return nil
}
