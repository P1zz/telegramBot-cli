package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/P1zz/telegramBot-cli/internal/validation"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type SendConfig struct {
	Token          string
	ChatID         int64
	Message        string
	MarkdownV2     bool
	FileHasSpoiler bool
	PrintMessageID bool
	FilePath       string
	FileTimeout    int
	FileIsImage    bool
	FileIsVideo    bool
	ReplyChatID    int64
	ReplyMessageID int
}

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send a text message, image, or video",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := loadSendConfig()
		if err := cfg.Validate(); err != nil {
			return err
		}
		message, err := sendMessage(cmd.Context(), cfg)
		if err != nil {
			return err
		}
		if cfg.PrintMessageID {
			fmt.Fprintf(cmd.OutOrStdout(), "MESSAGE_ID:%d\n", message.ID)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)

	flags := sendCmd.Flags()
	flags.StringP("message", "m", "", "message text or media caption")
	flags.StringP("file", "p", "", "path of the image or video to send")
	flags.IntP("file-timeout", "T", 60, "file upload timeout in seconds")
	flags.BoolP("image", "i", false, "send the file as an image")
	flags.BoolP("video", "v", false, "send the file as a video")
	flags.Bool("spoiler", false, "hide the media preview behind a spoiler")
	flags.Int64("reply-chat-id", 0, "chat ID containing the message to reply to")
	flags.Int("reply-message-id", 0, "message ID to reply to")
	flags.Bool("markdown-v2", false, "parse the message as Telegram MarkdownV2")
	flags.Bool("print-message-id", false, "print the sent message ID")

	bindFlag("send.message", sendCmd, "message")
	bindFlag("send.filePath", sendCmd, "file")
	bindFlag("send.fileTimeout", sendCmd, "file-timeout")
	bindFlag("send.fileIsImage", sendCmd, "image")
	bindFlag("send.fileIsVideo", sendCmd, "video")
	bindFlag("send.fileHasSpoiler", sendCmd, "spoiler")
	bindFlag("send.replyChatID", sendCmd, "reply-chat-id")
	bindFlag("send.replyMessageID", sendCmd, "reply-message-id")
	bindFlag("send.markdownV2", sendCmd, "markdown-v2")
	bindFlag("send.printMessageID", sendCmd, "print-message-id")
}

func bindFlag(key string, command *cobra.Command, name string) {
	cobra.CheckErr(viper.BindPFlag(key, command.Flags().Lookup(name)))
}

func loadSendConfig() SendConfig {
	return SendConfig{
		Token:          viper.GetString("token"),
		ChatID:         viper.GetInt64("chatId"),
		Message:        viper.GetString("send.message"),
		MarkdownV2:     viper.GetBool("send.markdownV2"),
		FileHasSpoiler: viper.GetBool("send.fileHasSpoiler"),
		PrintMessageID: viper.GetBool("send.printMessageID"),
		FilePath:       viper.GetString("send.filePath"),
		FileTimeout:    viper.GetInt("send.fileTimeout"),
		FileIsImage:    viper.GetBool("send.fileIsImage"),
		FileIsVideo:    viper.GetBool("send.fileIsVideo"),
		ReplyChatID:    viper.GetInt64("send.replyChatID"),
		ReplyMessageID: viper.GetInt("send.replyMessageID"),
	}
}

func (cfg SendConfig) Validate() error {
	if err := validation.Token(cfg.Token); err != nil {
		return err
	}
	if err := validation.RequiredChatID(cfg.ChatID); err != nil {
		return err
	}
	if cfg.Message == "" && !cfg.FileIsImage && !cfg.FileIsVideo {
		return errors.New("a message, image, or video is required")
	}
	if cfg.FileIsImage && cfg.FileIsVideo {
		return errors.New("image and video modes are mutually exclusive")
	}
	if cfg.FileIsImage || cfg.FileIsVideo {
		if err := validation.Positive("file timeout", cfg.FileTimeout); err != nil {
			return err
		}
		if err := validation.ReadableRegularFile(cfg.FilePath); err != nil {
			return err
		}
	}
	if cfg.ReplyChatID != 0 && cfg.ReplyMessageID == 0 {
		return errors.New("reply message ID is required when reply chat ID is set")
	}
	if cfg.ReplyMessageID < 0 {
		return errors.New("reply message ID must not be negative")
	}
	return nil
}

func sendMessage(ctx context.Context, cfg SendConfig) (*models.Message, error) {
	tgBot, err := bot.New(cfg.Token)
	if err != nil {
		return nil, fmt.Errorf("create bot: %w", err)
	}

	var parseMode models.ParseMode
	if cfg.MarkdownV2 {
		parseMode = models.ParseModeMarkdown
	}

	var reply *models.ReplyParameters
	if cfg.ReplyMessageID != 0 {
		replyChatID := cfg.ReplyChatID
		if replyChatID == 0 {
			replyChatID = cfg.ChatID
		}
		reply = &models.ReplyParameters{ChatID: replyChatID, MessageID: cfg.ReplyMessageID}
	}

	if !cfg.FileIsImage && !cfg.FileIsVideo {
		message, err := tgBot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: cfg.ChatID, Text: cfg.Message, ReplyParameters: reply, ParseMode: parseMode,
		})
		if err != nil {
			return nil, fmt.Errorf("send message: %w", err)
		}
		return message, nil
	}

	file, err := os.Open(cfg.FilePath)
	if err != nil {
		return nil, fmt.Errorf("open file %q: %w", cfg.FilePath, err)
	}
	defer file.Close()

	uploadCtx, cancel := context.WithTimeout(ctx, time.Duration(cfg.FileTimeout)*time.Second)
	defer cancel()
	input := &models.InputFileUpload{Filename: cfg.FilePath, Data: file}

	var message *models.Message
	if cfg.FileIsImage {
		message, err = tgBot.SendPhoto(uploadCtx, &bot.SendPhotoParams{
			ChatID: cfg.ChatID, Photo: input, Caption: cfg.Message, ReplyParameters: reply,
			ParseMode: parseMode, HasSpoiler: cfg.FileHasSpoiler,
		})
	} else {
		message, err = tgBot.SendVideo(uploadCtx, &bot.SendVideoParams{
			ChatID: cfg.ChatID, Video: input, Caption: cfg.Message, ReplyParameters: reply,
			ParseMode: parseMode, HasSpoiler: cfg.FileHasSpoiler,
		})
	}
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("upload exceeded the %d-second timeout: %w", cfg.FileTimeout, err)
		}
		return nil, fmt.Errorf("send media: %w", err)
	}
	return message, nil
}
