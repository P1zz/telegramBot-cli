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
		if err := editMessage(cmd.Context(), cfg); err != nil {
			return err
		}
		if output == "json" {
			return json.NewEncoder(cmd.OutOrStdout()).Encode(map[string]any{"edited": true, "chat_id": cfg.ChatID, "message_id": cfg.OldMessageID})
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
	editCmd.Flags().IntP("message-id", "i", 0, "ID of the message to edit")
	editCmd.Flags().StringP("message", "m", "", "replacement message text")
	editCmd.Flags().StringP("output", "o", "text", "output format: text or json")
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
	tgBot, err := newTelegramClient(cfg.Token)
	if err != nil {
		return safeError("create bot", cfg.Token, err)
	}
	_, err = tgBot.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID: cfg.ChatID, MessageID: cfg.OldMessageID, Text: cfg.NewMessage,
	})
	if err != nil {
		return safeError("edit message", cfg.Token, err)
	}
	return nil
}
