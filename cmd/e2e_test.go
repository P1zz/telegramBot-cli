//go:build e2e

package cmd

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestTelegramLifecycleE2E(t *testing.T) {
	if os.Getenv("TELEGRAM_CLI_E2E") != "1" {
		t.Skip("set TELEGRAM_CLI_E2E=1 to run")
	}
	token := os.Getenv("TELEGRAM_CLI_TOKEN")
	chatID, err := strconv.ParseInt(os.Getenv("TELEGRAM_CLI_CHAT_ID"), 10, 64)
	if err != nil || token == "" {
		t.Fatal("TELEGRAM_CLI_TOKEN and TELEGRAM_CLI_CHAT_ID are required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	message, err := sendMessage(ctx, SendConfig{Token: token, ChatID: chatID, Message: "telegram-cli e2e"})
	if err != nil {
		t.Fatal(err)
	}
	if err = editMessage(ctx, EditConfig{Token: token, ChatID: chatID, OldMessageID: message.ID, NewMessage: "telegram-cli e2e edited"}); err != nil {
		t.Fatal(err)
	}
	if err = deleteMessage(ctx, DeleteConfig{Token: token, ChatID: chatID, MessageID: message.ID}); err != nil {
		t.Fatal(err)
	}
}
