package cmd

import (
	"context"
	"encoding/json"
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
		var err error
		cfg.Token, err = resolvedToken()
		if err != nil {
			return err
		}
		if err := cfg.Validate(); err != nil {
			return err
		}
		output, _ := cmd.Flags().GetString("output")
		if output != "text" && output != "json" {
			return fmt.Errorf("output must be \"text\" or \"json\"")
		}
		if err := deleteMessage(cmd.Context(), cfg); err != nil {
			return err
		}
		if output == "json" {
			return json.NewEncoder(cmd.OutOrStdout()).Encode(map[string]any{"deleted": true, "chat_id": cfg.ChatID, "message_id": cfg.MessageID})
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().IntP("message-id", "i", 0, "ID of the message to delete")
	deleteCmd.Flags().StringP("output", "o", "text", "output format: text or json")
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
	tgBot, err := newTelegramClient(cfg.Token)
	if err != nil {
		return safeError("create bot", cfg.Token, err)
	}
	deleted, err := tgBot.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID: cfg.ChatID, MessageID: cfg.MessageID,
	})
	if err != nil {
		return safeError("delete message", cfg.Token, err)
	}
	if !deleted {
		return fmt.Errorf("telegram did not delete message %d", cfg.MessageID)
	}
	return nil
}
