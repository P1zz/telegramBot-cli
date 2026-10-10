package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
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
	FileIsDocument bool
	FileIsAudio    bool
	FileIsVoice    bool
	Stdin          bool
	Silent         bool
	ProtectContent bool
	Output         string
	ReplyChatID    int64
	ReplyMessageID int
}

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send text or media",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		applyLegacyFlags(cmd, map[string]string{"messageText": "send.message", "filePath": "send.filePath", "pathIsImage": "send.fileIsImage", "pathIsVideo": "send.fileIsVideo", "printMessageId": "send.printMessageID"})
		cfg := loadSendConfig()
		if tokenFile, _ := cmd.Flags().GetString("token-file"); cfg.Stdin && tokenFile == "-" {
			return errors.New("--stdin cannot be combined with --token-file -")
		}
		var err error
		cfg.Token, err = resolvedToken()
		if err != nil {
			return err
		}
		if cfg.Stdin {
			data, readErr := io.ReadAll(cmd.InOrStdin())
			if readErr != nil {
				return fmt.Errorf("read stdin: %w", readErr)
			}
			cfg.Message = strings.TrimSuffix(string(data), "\n")
		}
		if err := cfg.Validate(); err != nil {
			return err
		}
		message, err := sendMessage(cmd.Context(), cfg)
		if err != nil {
			return err
		}
		if cfg.Output == "json" {
			return json.NewEncoder(cmd.OutOrStdout()).Encode(map[string]any{"message_id": message.ID, "chat_id": cfg.ChatID})
		} else if cfg.PrintMessageID {
			fmt.Fprintf(cmd.OutOrStdout(), "MESSAGE_ID:%d\n", message.ID)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)

	flags := sendCmd.Flags()
	flags.StringP("message", "m", "", "message text or media caption")
	flags.String("messageText", "", "deprecated alias for --message")
	flags.StringP("file", "p", "", "path of the image or video to send")
	flags.IntP("file-timeout", "T", 60, "file upload timeout in seconds")
	flags.BoolP("image", "i", false, "send the file as an image")
	flags.BoolP("video", "v", false, "send the file as a video")
	flags.Bool("document", false, "send the file as a document")
	flags.Bool("audio", false, "send the file as audio")
	flags.Bool("voice", false, "send the file as a voice message")
	flags.Bool("stdin", false, "read message text from stdin")
	flags.Bool("silent", false, "send without a notification")
	flags.Bool("protect-content", false, "prevent forwarding and saving")
	flags.StringP("output", "o", "text", "output format: text or json")
	flags.Bool("spoiler", false, "hide the media preview behind a spoiler")
	flags.Int64("reply-chat-id", 0, "chat ID containing the message to reply to")
	flags.Int("reply-message-id", 0, "message ID to reply to")
	flags.Bool("markdown-v2", false, "parse the message as Telegram MarkdownV2")
	flags.Bool("print-message-id", false, "print the sent message ID")
	flags.String("filePath", "", "deprecated alias for --file")
	flags.Bool("pathIsImage", false, "deprecated alias for --image")
	flags.Bool("pathIsVideo", false, "deprecated alias for --video")
	flags.Bool("printMessageId", false, "deprecated alias for --print-message-id")
	for old, replacement := range map[string]string{"messageText": "message", "filePath": "file", "pathIsImage": "image", "pathIsVideo": "video", "printMessageId": "print-message-id"} {
		_ = flags.MarkDeprecated(old, "use --"+replacement)
	}

	bindFlag("send.message", sendCmd, "message")
	bindFlag("send.filePath", sendCmd, "file")
	bindFlag("send.fileTimeout", sendCmd, "file-timeout")
	bindFlag("send.fileIsImage", sendCmd, "image")
	bindFlag("send.fileIsVideo", sendCmd, "video")
	bindFlag("send.fileIsDocument", sendCmd, "document")
	bindFlag("send.fileIsAudio", sendCmd, "audio")
	bindFlag("send.fileIsVoice", sendCmd, "voice")
	bindFlag("send.stdin", sendCmd, "stdin")
	bindFlag("send.silent", sendCmd, "silent")
	bindFlag("send.protectContent", sendCmd, "protect-content")
	bindFlag("send.output", sendCmd, "output")
	bindFlag("send.fileHasSpoiler", sendCmd, "spoiler")
	bindFlag("send.replyChatID", sendCmd, "reply-chat-id")
	bindFlag("send.replyMessageID", sendCmd, "reply-message-id")
	bindFlag("send.markdownV2", sendCmd, "markdown-v2")
	bindFlag("send.printMessageID", sendCmd, "print-message-id")
}

func bindFlag(key string, command *cobra.Command, name string) {
	cobra.CheckErr(viper.BindPFlag(key, command.Flags().Lookup(name)))
}

func applyLegacyFlags(command *cobra.Command, aliases map[string]string) {
	for flag, key := range aliases {
		if !command.Flags().Changed(flag) {
			continue
		}
		value, err := command.Flags().GetString(flag)
		if err == nil {
			viper.Set(key, value)
			continue
		}
		if boolean, boolErr := command.Flags().GetBool(flag); boolErr == nil {
			viper.Set(key, boolean)
		}
	}
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
		FileIsDocument: viper.GetBool("send.fileIsDocument"),
		FileIsAudio:    viper.GetBool("send.fileIsAudio"),
		FileIsVoice:    viper.GetBool("send.fileIsVoice"),
		Stdin:          viper.GetBool("send.stdin"), Silent: viper.GetBool("send.silent"),
		ProtectContent: viper.GetBool("send.protectContent"), Output: viper.GetString("send.output"),
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
	mediaCount := 0
	for _, enabled := range []bool{cfg.FileIsImage, cfg.FileIsVideo, cfg.FileIsDocument, cfg.FileIsAudio, cfg.FileIsVoice} {
		if enabled {
			mediaCount++
		}
	}
	if cfg.Message == "" && mediaCount == 0 {
		return errors.New("a message, image, or video is required")
	}
	if mediaCount > 1 {
		return errors.New("media modes are mutually exclusive")
	}
	if mediaCount == 1 {
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
	if cfg.Output != "" && cfg.Output != "text" && cfg.Output != "json" {
		return errors.New("output must be \"text\" or \"json\"")
	}
	return nil
}

func sendMessage(ctx context.Context, cfg SendConfig) (*models.Message, error) {
	tgBot, err := newTelegramClient(cfg.Token)
	if err != nil {
		return nil, safeError("create bot", cfg.Token, err)
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

	if !cfg.FileIsImage && !cfg.FileIsVideo && !cfg.FileIsDocument && !cfg.FileIsAudio && !cfg.FileIsVoice {
		message, err := tgBot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: cfg.ChatID, Text: cfg.Message, ReplyParameters: reply, ParseMode: parseMode,
			DisableNotification: cfg.Silent, ProtectContent: cfg.ProtectContent,
		})
		if err != nil {
			return nil, safeError("send message", cfg.Token, err)
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
			ParseMode: parseMode, HasSpoiler: cfg.FileHasSpoiler, DisableNotification: cfg.Silent, ProtectContent: cfg.ProtectContent,
		})
	} else if cfg.FileIsVideo {
		message, err = tgBot.SendVideo(uploadCtx, &bot.SendVideoParams{
			ChatID: cfg.ChatID, Video: input, Caption: cfg.Message, ReplyParameters: reply,
			ParseMode: parseMode, HasSpoiler: cfg.FileHasSpoiler, DisableNotification: cfg.Silent, ProtectContent: cfg.ProtectContent,
		})
	} else if cfg.FileIsDocument {
		message, err = tgBot.SendDocument(uploadCtx, &bot.SendDocumentParams{ChatID: cfg.ChatID, Document: input, Caption: cfg.Message, ParseMode: parseMode, ReplyParameters: reply, DisableNotification: cfg.Silent, ProtectContent: cfg.ProtectContent})
	} else if cfg.FileIsAudio {
		message, err = tgBot.SendAudio(uploadCtx, &bot.SendAudioParams{ChatID: cfg.ChatID, Audio: input, Caption: cfg.Message, ParseMode: parseMode, ReplyParameters: reply, DisableNotification: cfg.Silent, ProtectContent: cfg.ProtectContent})
	} else {
		message, err = tgBot.SendVoice(uploadCtx, &bot.SendVoiceParams{ChatID: cfg.ChatID, Voice: input, Caption: cfg.Message, ParseMode: parseMode, ReplyParameters: reply, DisableNotification: cfg.Silent, ProtectContent: cfg.ProtectContent})
	}
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("upload exceeded the %d-second timeout: %w", cfg.FileTimeout, err)
		}
		return nil, safeError("send media", cfg.Token, err)
	}
	return message, nil
}
