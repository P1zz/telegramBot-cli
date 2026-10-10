package telegram

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// Client is the Telegram API surface used by the CLI. Keeping it small makes
// command behavior testable without network access.
type Client interface {
	SendMessage(context.Context, *bot.SendMessageParams) (*models.Message, error)
	SendPhoto(context.Context, *bot.SendPhotoParams) (*models.Message, error)
	SendVideo(context.Context, *bot.SendVideoParams) (*models.Message, error)
	SendDocument(context.Context, *bot.SendDocumentParams) (*models.Message, error)
	SendAudio(context.Context, *bot.SendAudioParams) (*models.Message, error)
	SendVoice(context.Context, *bot.SendVoiceParams) (*models.Message, error)
	EditMessageText(context.Context, *bot.EditMessageTextParams) (*models.Message, error)
	DeleteMessage(context.Context, *bot.DeleteMessageParams) (bool, error)
	GetFile(context.Context, *bot.GetFileParams) (*models.File, error)
	Start(context.Context)
}

type Factory func(string, ...bot.Option) (Client, error)

func New(token string, options ...bot.Option) (Client, error) {
	return bot.New(token, options...)
}
