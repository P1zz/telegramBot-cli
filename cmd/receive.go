package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"sync"
	"time"

	clioutput "github.com/P1zz/telegramBot-cli/internal/output"
	telegramapi "github.com/P1zz/telegramBot-cli/internal/telegram"
	"github.com/P1zz/telegramBot-cli/internal/validation"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ReceiveConfig struct {
	Token               string
	ChatID              int64
	Counter             int
	Timeout             int
	Sync                bool
	PrintChatID         bool
	PrintMessageID      bool
	PrintTimestampUnix  bool
	PrintTimestampHuman bool
	PrintPhotoPath      bool
	PrintFilePath       bool
	PrintAudioPath      bool
	Output              string
}

var receiveCmd = &cobra.Command{
	Use:   "receive",
	Short: "Receive Telegram messages",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		applyLegacyFlags(cmd, map[string]string{"printChatId": "receive.printChatId", "printMessageId": "receive.printMessageId", "printPhotoUrl": "receive.printPhotoPath", "printFileUrl": "receive.printFilePath", "printAudioUrl": "receive.printAudioPath"})
		cfg := loadReceiveConfig()
		var err error
		cfg.Token, err = resolvedToken()
		if err != nil {
			return err
		}
		if err := cfg.Validate(); err != nil {
			return err
		}
		return receiveMessages(cmd.Context(), cfg, cmd.OutOrStdout(), cmd.ErrOrStderr())
	},
}

func init() {
	rootCmd.AddCommand(receiveCmd)
	flags := receiveCmd.Flags()
	flags.IntP("counter", "n", 0, "stop after receiving this many messages (0 means unlimited)")
	flags.IntP("timeout", "w", 0, "stop after this many seconds (0 means no timeout)")
	flags.BoolP("sync", "s", false, "include queued messages sent before startup")
	flags.Bool("print-chat-id", false, "include the chat ID in text output")
	flags.Bool("print-message-id", false, "include the message ID in text output")
	flags.Bool("print-timestamp-unix", false, "include the message time as a Unix timestamp")
	flags.Bool("print-timestamp-human", false, "include the message time in RFC3339 format")
	flags.Bool("print-photo-path", false, "include Telegram's photo file path")
	flags.Bool("print-file-path", false, "include Telegram's document file path")
	flags.Bool("print-audio-path", false, "include Telegram's audio or voice file path")
	flags.StringP("output", "o", clioutput.FormatText, "output format: text or jsonl")
	for _, old := range []string{"printChatId", "printMessageId", "printPhotoUrl", "printFileUrl", "printAudioUrl"} {
		flags.Bool(old, false, "deprecated compatibility flag")
		_ = flags.MarkDeprecated(old, "use the kebab-case path flag shown in --help")
	}

	bindFlag("receive.counter", receiveCmd, "counter")
	bindFlag("receive.timeout", receiveCmd, "timeout")
	bindFlag("receive.sync", receiveCmd, "sync")
	bindFlag("receive.printChatId", receiveCmd, "print-chat-id")
	bindFlag("receive.printMessageId", receiveCmd, "print-message-id")
	bindFlag("receive.printTimestampUnix", receiveCmd, "print-timestamp-unix")
	bindFlag("receive.printTimestampHuman", receiveCmd, "print-timestamp-human")
	bindFlag("receive.printPhotoPath", receiveCmd, "print-photo-path")
	bindFlag("receive.printFilePath", receiveCmd, "print-file-path")
	bindFlag("receive.printAudioPath", receiveCmd, "print-audio-path")
	bindFlag("receive.output", receiveCmd, "output")
}

func loadReceiveConfig() ReceiveConfig {
	return ReceiveConfig{
		Token: viper.GetString("token"), ChatID: viper.GetInt64("chatId"),
		Counter: viper.GetInt("receive.counter"), Timeout: viper.GetInt("receive.timeout"),
		Sync: viper.GetBool("receive.sync"), PrintChatID: viper.GetBool("receive.printChatId"),
		PrintMessageID:      viper.GetBool("receive.printMessageId"),
		PrintTimestampUnix:  viper.GetBool("receive.printTimestampUnix"),
		PrintTimestampHuman: viper.GetBool("receive.printTimestampHuman"),
		PrintPhotoPath:      viper.GetBool("receive.printPhotoPath"),
		PrintFilePath:       viper.GetBool("receive.printFilePath"),
		PrintAudioPath:      viper.GetBool("receive.printAudioPath"),
		Output:              viper.GetString("receive.output"),
	}
}

func (cfg ReceiveConfig) Validate() error {
	if err := validation.Token(cfg.Token); err != nil {
		return err
	}
	if err := validation.NonNegative("counter", cfg.Counter); err != nil {
		return err
	}
	if err := validation.NonNegative("timeout", cfg.Timeout); err != nil {
		return err
	}
	if cfg.PrintTimestampUnix && cfg.PrintTimestampHuman {
		return errors.New("unix and human-readable timestamps are mutually exclusive")
	}
	return clioutput.ValidateFormat(cfg.Output)
}

func receiveMessages(parent context.Context, cfg ReceiveConfig, stdout, stderr io.Writer) error {
	ctx, cancel := context.WithCancel(parent)
	defer cancel()
	if cfg.Timeout > 0 {
		var timeoutCancel context.CancelFunc
		ctx, timeoutCancel = context.WithTimeout(ctx, time.Duration(cfg.Timeout)*time.Second)
		defer timeoutCancel()
	}

	startedAt := time.Now().Unix()
	var received int
	var handlerErr error
	var handlerMu sync.Mutex

	options := []bot.Option{
		bot.WithErrorsHandler(func(err error) {
			_, _ = fmt.Fprintf(stderr, "Telegram polling error: %v\n", err)
		}),
		bot.WithDefaultHandler(func(handlerCtx context.Context, tgBot *bot.Bot, update *models.Update) {
			handlerMu.Lock()
			defer handlerMu.Unlock()
			if err := handleUpdate(handlerCtx, tgBot, update, cfg, startedAt, stdout); err != nil {
				handlerErr = err
				cancel()
				return
			}
			if updateAccepted(update, cfg, startedAt) && cfg.Counter > 0 {
				received++
				if received >= cfg.Counter {
					cancel()
				}
			}
		}),
	}
	if !cfg.Sync {
		options = append(options, bot.WithInitialOffset(-1))
	}

	tgBot, err := newTelegramClient(cfg.Token, options...)
	if err != nil {
		return safeError("create bot", cfg.Token, err)
	}
	tgBot.Start(ctx)
	handlerMu.Lock()
	defer handlerMu.Unlock()
	if handlerErr != nil {
		return handlerErr
	}
	if err := parent.Err(); err != nil {
		return err
	}
	return nil
}

func updateAccepted(update *models.Update, cfg ReceiveConfig, startedAt int64) bool {
	message, _ := messageFromUpdate(update)
	if message == nil {
		return false
	}
	if !cfg.Sync && int64(message.Date) < startedAt {
		return false
	}
	return cfg.ChatID == 0 || message.Chat.ID == cfg.ChatID
}

func handleUpdate(ctx context.Context, tgBot telegramapi.Client, update *models.Update, cfg ReceiveConfig, startedAt int64, writer io.Writer) error {
	if !updateAccepted(update, cfg, startedAt) {
		return nil
	}
	message, event := messageFromUpdate(update)
	result := clioutput.Message{Event: event, Text: message.Text, FileCaption: message.Caption}
	if cfg.PrintChatID || cfg.Output == clioutput.FormatJSONL {
		result.ChatID = message.Chat.ID
	}
	if cfg.PrintMessageID || cfg.Output == clioutput.FormatJSONL {
		result.MessageID = message.ID
	}
	if cfg.PrintTimestampHuman {
		result.Date = time.Unix(int64(message.Date), 0).Format(time.RFC3339)
	} else if cfg.PrintTimestampUnix {
		result.Date = strconv.Itoa(message.Date)
	}

	fileID, fileName := selectedFile(message, cfg)
	if fileID != "" {
		file, err := tgBot.GetFile(ctx, &bot.GetFileParams{FileID: fileID})
		if err != nil {
			return safeError("get Telegram file metadata", cfg.Token, err)
		}
		// FilePath is safe to log. A downloadable URL would contain the secret bot token.
		result.FilePath = file.FilePath
		result.FileName = fileName
	}
	return clioutput.Write(writer, cfg.Output, result)
}

func messageFromUpdate(update *models.Update) (*models.Message, string) {
	if update == nil {
		return nil, ""
	}
	switch {
	case update.Message != nil:
		return update.Message, "message"
	case update.EditedMessage != nil:
		return update.EditedMessage, "edited_message"
	case update.ChannelPost != nil:
		return update.ChannelPost, "channel_post"
	case update.EditedChannelPost != nil:
		return update.EditedChannelPost, "edited_channel_post"
	default:
		return nil, ""
	}
}

func selectedFile(message *models.Message, cfg ReceiveConfig) (fileID, fileName string) {
	if cfg.PrintPhotoPath && len(message.Photo) > 0 {
		return message.Photo[len(message.Photo)-1].FileID, ""
	}
	if cfg.PrintFilePath && message.Document != nil {
		return message.Document.FileID, message.Document.FileName
	}
	if cfg.PrintAudioPath && message.Audio != nil {
		return message.Audio.FileID, message.Audio.FileName
	}
	if cfg.PrintAudioPath && message.Voice != nil {
		return message.Voice.FileID, ""
	}
	return "", ""
}
