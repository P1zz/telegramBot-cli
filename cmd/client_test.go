package cmd

import (
	"context"
	"testing"

	telegramapi "github.com/P1zz/telegramBot-cli/internal/telegram"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type fakeClient struct {
	telegramapi.Client
	sendMessageParams *bot.SendMessageParams
}

func (client *fakeClient) SendMessage(_ context.Context, params *bot.SendMessageParams) (*models.Message, error) {
	client.sendMessageParams = params
	return &models.Message{ID: 42}, nil
}

func TestSendMessageBuildsReplyAndOptions(t *testing.T) {
	fake := &fakeClient{}
	previous := newTelegramClient
	newTelegramClient = func(string, ...bot.Option) (telegramapi.Client, error) { return fake, nil }
	t.Cleanup(func() { newTelegramClient = previous })

	message, err := sendMessage(context.Background(), SendConfig{
		Token: "secret", ChatID: -1001, Message: "hello", ReplyMessageID: 7,
		Silent: true, ProtectContent: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if message.ID != 42 {
		t.Fatalf("message ID = %d", message.ID)
	}
	params := fake.sendMessageParams
	if params == nil || params.ReplyParameters == nil {
		t.Fatal("reply parameters were not supplied")
	}
	if params.ReplyParameters.ChatID != int64(-1001) || params.ReplyParameters.MessageID != 7 {
		t.Fatalf("unexpected reply: %#v", params.ReplyParameters)
	}
	if !params.DisableNotification || !params.ProtectContent {
		t.Fatal("send options were not propagated")
	}
}
