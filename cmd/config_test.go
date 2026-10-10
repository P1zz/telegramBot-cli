package cmd

import (
	"testing"

	"github.com/go-telegram/bot/models"
)

func TestSendConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		cfg     SendConfig
		wantErr bool
	}{
		{name: "text", cfg: SendConfig{Token: "token", ChatID: -100123, Message: "hello"}},
		{name: "missing content", cfg: SendConfig{Token: "token", ChatID: 123}, wantErr: true},
		{name: "two media modes", cfg: SendConfig{Token: "token", ChatID: 123, FileIsImage: true, FileIsVideo: true}, wantErr: true},
		{name: "orphan reply chat", cfg: SendConfig{Token: "token", ChatID: 123, Message: "hello", ReplyChatID: 456}, wantErr: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := test.cfg.Validate(); (err != nil) != test.wantErr {
				t.Fatalf("Validate() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}

func TestReceiveConfigValidation(t *testing.T) {
	if err := (ReceiveConfig{Token: "token", Output: "jsonl"}).Validate(); err != nil {
		t.Fatal(err)
	}
	if err := (ReceiveConfig{Token: "token", Output: "text", Counter: -1}).Validate(); err == nil {
		t.Fatal("negative counter was accepted")
	}
	if err := (ReceiveConfig{Token: "token", Output: "text", PrintTimestampHuman: true, PrintTimestampUnix: true}).Validate(); err == nil {
		t.Fatal("two timestamp formats were accepted")
	}
}

func TestUpdateAccepted(t *testing.T) {
	update := &models.Update{Message: &models.Message{Date: 100, Chat: models.Chat{ID: -1001}}}
	if !updateAccepted(update, ReceiveConfig{Sync: true}, 200) {
		t.Fatal("sync mode rejected an old update")
	}
	if updateAccepted(update, ReceiveConfig{}, 200) {
		t.Fatal("non-sync mode accepted an old update")
	}
	if updateAccepted(update, ReceiveConfig{Sync: true, ChatID: 99}, 200) {
		t.Fatal("accepted an update from another chat")
	}
}

func TestMessageFromUpdate(t *testing.T) {
	message := &models.Message{ID: 9}
	got, event := messageFromUpdate(&models.Update{EditedChannelPost: message})
	if got != message || event != "edited_channel_post" {
		t.Fatalf("got %#v, %q", got, event)
	}
}

func TestSelectedFilePrefersPhoto(t *testing.T) {
	message := &models.Message{
		Photo:    []models.PhotoSize{{FileID: "small"}, {FileID: "large"}},
		Document: &models.Document{FileID: "document", FileName: "report.pdf"},
	}
	fileID, fileName := selectedFile(message, ReceiveConfig{PrintPhotoPath: true, PrintFilePath: true})
	if fileID != "large" || fileName != "" {
		t.Fatalf("selectedFile() = %q, %q", fileID, fileName)
	}
}
