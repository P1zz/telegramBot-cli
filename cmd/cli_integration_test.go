package cmd

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	telegramapi "github.com/P1zz/telegramBot-cli/internal/telegram"
	"github.com/go-telegram/bot"
)

func TestCLIFlagEnvironmentAndConfigPrecedence(t *testing.T) {
	directory := t.TempDir()
	path := filepath.Join(directory, "config.toml")
	if err := os.WriteFile(path, []byte("chatId = 111\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("TELEGRAM_CLI_TOKEN", "environment-token")

	fake := &fakeClient{}
	previousFactory := newTelegramClient
	newTelegramClient = func(token string, _ ...bot.Option) (telegramapi.Client, error) {
		if token != "environment-token" {
			t.Fatalf("token = %q", token)
		}
		return fake, nil
	}
	t.Cleanup(func() { newTelegramClient = previousFactory })

	var stdout bytes.Buffer
	rootCmd.SetOut(&stdout)
	rootCmd.SetErr(&stdout)
	rootCmd.SetArgs([]string{"--config", path, "send", "--chat-id", "222", "--message", "hello", "--output", "json"})
	if err := rootCmd.ExecuteContext(context.Background()); err != nil {
		t.Fatal(err)
	}
	if fake.sendMessageParams.ChatID != int64(222) {
		t.Fatalf("chat ID = %#v", fake.sendMessageParams.ChatID)
	}
	if !strings.Contains(stdout.String(), `"message_id":42`) {
		t.Fatalf("stdout = %q", stdout.String())
	}
}
