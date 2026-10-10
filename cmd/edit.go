package cmd

import (
	"context"
	"fmt"

	"github.com/P1zz/telegramBot-cli/internal/validation"

	"github.com/go-telegram/bot"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type EditConfig struct {
	Token        string
	ChatID       int64
	NewMessage   string
	OldMessageID int
}

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit a text message",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := loadEditConfig()
		if err := cfg.Validate(); err != nil {
			return err
		}
		return editMessage(cmd.Context(), cfg)
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
	editCmd.Flags().IntP("message-id", "i", 0, "ID of the message to edit")
	editCmd.Flags().StringP("message", "m", "", "replacement message text")
	bindFlag("edit.oldMessageId", editCmd, "message-id")
	bindFlag("edit.newMessage", editCmd, "message")
}

func loadEditConfig() EditConfig {
	return EditConfig{
		Token:        viper.GetString("token"),
		ChatID:       viper.GetInt64("chatId"),
		NewMessage:   viper.GetString("edit.newMessage"),
		OldMessageID: viper.GetInt("edit.oldMessageId"),
	}
}

func (cfg EditConfig) Validate() error {
	if err := validation.Token(cfg.Token); err != nil {
		return err
	}
	if err := validation.RequiredChatID(cfg.ChatID); err != nil {
		return err
	}
	if err := validation.MessageID(cfg.OldMessageID); err != nil {
		return err
	}
	if cfg.NewMessage == "" {
		return fmt.Errorf("new message is required")
	}
	return nil
}

func editMessage(ctx context.Context, cfg EditConfig) error {
	tgBot, err := bot.New(cfg.Token)
	if err != nil {
		return fmt.Errorf("create bot: %w", err)
	}
	_, err = tgBot.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID: cfg.ChatID, MessageID: cfg.OldMessageID, Text: cfg.NewMessage,
	})
	if err != nil {
		return fmt.Errorf("edit message: %w", err)
	}
	return nil
}
