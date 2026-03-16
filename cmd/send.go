package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type SendConfig struct {
	token           string
	chatId          int
	message         string
	isMarkdownV2    bool
	fileHasSpoiler  bool
	getTheMessageId bool
	filePath        string
	fileTimeout     int
	fileIsImage     bool
	fileIsVideo     bool
	replyChatID     int
	replyMessageID  int
}

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send message with text or image",
	Long:  "Send a message in a chat as bot with text or an image",
	Args:  validateArgsSend,
	//Link the validation function to the sendTextCmd
	RunE: sendMessage,
	//Link the function with the capabilities of returning an error
}

func init() {
	rootCmd.AddCommand(sendCmd)

	sendCmd.Flags().StringP("token", "t", "", "Token from bot fathers")
	sendCmd.Flags().IntP("chatId", "c", 0, "Your chat ID")
	sendCmd.Flags().StringP("messageText", "m", "", "Message text to send")
	sendCmd.Flags().StringP("filePath", "p", "", "Path of the image/video to send")
	sendCmd.Flags().IntP("fileTimeout", "T", 60, "Timeout in seconds for sending a file")
	sendCmd.Flags().BoolP("pathIsImage", "i", false, "The path is an image to send")
	sendCmd.Flags().BoolP("pathIsVideo", "v", false, "The path is a video to send")
	sendCmd.Flags().BoolP("fileHasSpoiler", "H", false, "The file is send with hidden preview")
	sendCmd.Flags().IntP("replyChatId", "x", 0, "Chat id you want to reply")
	sendCmd.Flags().IntP("replyMessageId", "y", 0, "Message id you want to reply")
	sendCmd.Flags().BoolP("markDownV2", "2", false, "Message text is parsed in markdown v2")
	sendCmd.Flags().BoolP("printMessageId", "M", false, "Print message id of your message")

	viper.BindPFlag("token", sendCmd.Flags().Lookup("token"))
	viper.BindPFlag("chatId", sendCmd.Flags().Lookup("chatId"))
	viper.BindPFlag("send.message", sendCmd.Flags().Lookup("messageText"))
	viper.BindPFlag("send.markdownV2", sendCmd.Flags().Lookup("markDownV2"))
	viper.BindPFlag("send.fileHasSpoiler", sendCmd.Flags().Lookup("fileHasSpoiler"))
	viper.BindPFlag("send.getTheMessageId", sendCmd.Flags().Lookup("printMessageId"))
	viper.BindPFlag("send.filePath", sendCmd.Flags().Lookup("filePath"))
	viper.BindPFlag("send.fileTimeout", sendCmd.Flags().Lookup("fileTimeout"))
	viper.BindPFlag("send.fileIsImage", sendCmd.Flags().Lookup("pathIsImage"))
	viper.BindPFlag("send.fileIsVideo", sendCmd.Flags().Lookup("pathIsVideo"))
	viper.BindPFlag("send.replyChatID", sendCmd.Flags().Lookup("replyChatId"))
	viper.BindPFlag("send.replyMessageID", sendCmd.Flags().Lookup("replyMessageId"))
}

func validateArgsSend(cmd *cobra.Command, args []string) error {

	cfg := SendConfig{
		token:           viper.GetString("token"),
		chatId:          viper.GetInt("chatId"),
		message:         viper.GetString("send.message"),
		isMarkdownV2:    viper.GetBool("send.markdownV2"),
		fileHasSpoiler:  viper.GetBool("send.fileHasSpoiler"),
		getTheMessageId: viper.GetBool("send.getTheMessageId"),
		filePath:        viper.GetString("send.filePath"),
		fileTimeout:     viper.GetInt("send.fileTimeout"),
		fileIsImage:     viper.GetBool("send.fileIsImage"),
		fileIsVideo:     viper.GetBool("send.fileIsVideo"),
		replyChatID:     viper.GetInt("send.replyChatID"),
		replyMessageID:  viper.GetInt("send.replyMessageID"),
	}

	//Validate the token
	if cfg.token == "" {
		return fmt.Errorf("No token provided")
	}

	//Validate the chat ID
	if cfg.chatId == 0 || len(strconv.Itoa(cfg.chatId)) != 9 {
		return fmt.Errorf("Wrong chat ID provided")
	}

	//There should be something to send
	if cfg.message == "" && (!cfg.fileIsImage && !cfg.fileIsVideo) {
		return fmt.Errorf("You should provide at least a message or a file")
	}

	//If there is a file to send
	if cfg.fileIsImage || cfg.fileIsVideo {
		//Check that are not both a video and an image
		if cfg.fileIsImage && cfg.fileIsVideo {
			return fmt.Errorf("Images and videos are mutually exclusive")
			//Check the file exist
		} else if _, err := os.Stat(cfg.filePath); errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("Wrong path provided")
		}
	}

	//Send config into the context
	cmd.SetContext(context.WithValue(cmd.Context(), SendConfig{}, cfg))
	return nil
}

func sendMessage(cmd *cobra.Command, args []string) error {

	cfg := cmd.Context().Value(SendConfig{}).(SendConfig)

	//Create a context
	bgCtx := context.Background()

	//Create the bot
	tgBot, err := bot.New(cfg.token)
	if err != nil {
		return err
	}

	//Create the return message structure
	var rtrn *models.Message

	//Create and fill parsing parameters
	var parsing models.ParseMode
	if cfg.isMarkdownV2 {
		parsing = models.ParseModeMarkdown
	}

	//If user does no has provided the chat ID use the current one
	if cfg.replyChatID == 0 {
		cfg.replyChatID = cfg.chatId
	}

	//Create and fill reply parameters
	replyParameters := &models.ReplyParameters{}
	if cfg.replyChatID != 0 {
		replyParameters.ChatID = cfg.replyChatID
		replyParameters.MessageID = cfg.replyMessageID
	}

	//Send a file
	if cfg.filePath != "" {
		//Open image
		file, err := os.Open(cfg.filePath)
		if err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.fileTimeout)*time.Second)
		defer cancel()

		if cfg.fileIsImage {
			//Create image parameters
			parameters := &bot.SendPhotoParams{
				ChatID:          cfg.chatId,
				Photo:           &models.InputFileUpload{Filename: cfg.filePath, Data: file},
				Caption:         cfg.message,
				ReplyParameters: replyParameters,
				ParseMode:       parsing,
				HasSpoiler:      cfg.fileHasSpoiler,
			}

			//Send image
			rtrn, err = tgBot.SendPhoto(ctx, parameters)

			//Check for errors
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					return fmt.Errorf("Send file request exceeded timeout of %d seconds, try a smaller file or increase -T", cfg.fileTimeout)
				}
				return err
			}
		} else if cfg.fileIsVideo {
			parameters := &bot.SendVideoParams{
				ChatID:          cfg.chatId,
				Video:           &models.InputFileUpload{Filename: cfg.filePath, Data: file},
				Caption:         cfg.message,
				ReplyParameters: replyParameters,
				ParseMode:       parsing,
				HasSpoiler:      cfg.fileHasSpoiler,
			}

			//Send video
			rtrn, err = tgBot.SendVideo(ctx, parameters)

			//Check for errors
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					return fmt.Errorf("Send file request exceeded timeout of %d seconds, try a smaller file or increase -T", cfg.fileTimeout)
				}
				return err
			}
		}

	} else { //Send a message

		parameters := &bot.SendMessageParams{
			ChatID:          cfg.chatId,
			Text:            cfg.message,
			ReplyParameters: replyParameters,
			ParseMode:       parsing,
		}

		//Send the message
		rtrn, err = tgBot.SendMessage(bgCtx, parameters)
	}

	//Check for errors
	if err != nil {
		return err
	}

	//If requested print messsage ID
	if cfg.getTheMessageId {
		fmt.Printf("MESSAGE_ID:%d\n", rtrn.ID)
	}

	//Close context
	bgCtx.Done()

	return nil
}
