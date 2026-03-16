package cmd

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ReceiveConfig struct {
	token               string
	chatId              int
	messageCounter      int
	sync                bool
	printChatId         bool
	printMessageId      bool
	printTimestampUnix  bool
	printTimestampHuman bool
	printPhotoUrl       bool
	printFileUrl        bool
	printAudioUrl       bool
}

var receiveTextCmd = &cobra.Command{
	Use:   "receive",
	Short: "Receive message with text",
	Long:  "Receive a message as bot with the pattern below\n|DATA|CHAT_ID|MESSAGE_ID|URL|FILE_NAME|FILE_CAPTION|TEXT|",
	//Link the validation function to the receiveTextCmd
	Args: validateArgsReceive,
	//Link the function with the capabilities of returning an error
	RunE: receiveMessage,
}

func init() {
	rootCmd.AddCommand(receiveTextCmd)

	receiveTextCmd.Flags().StringP("token", "t", "", "Token from bot fathers")
	receiveTextCmd.Flags().IntP("chatId", "c", 0, "ID of the chat, leave blank or set 0 if you want to listen all chats")
	receiveTextCmd.Flags().IntP("messageCounter", "n", 0, "Numer of messages to receive, leave blank or set 0 for continuous receiving")
	//TODO document cron logic for sync
	receiveTextCmd.Flags().BoolP("sync", "s", false, "Sync old messages sended while the bot was not running")
	receiveTextCmd.Flags().BoolP("printChatId", "r", false, "Print the chat ID")
	receiveTextCmd.Flags().BoolP("printMessageId", "M", false, "Print the message ID of each message")
	receiveTextCmd.Flags().BoolP("printTimestampUnix", "U", false, "Print the datetime UNIX")
	receiveTextCmd.Flags().BoolP("printTimestampHuman", "H", false, "Print the datetime human readable")
	receiveTextCmd.Flags().BoolP("printPhotoUrl", "P", false, "Print the photo url")
	receiveTextCmd.Flags().BoolP("printFileUrl", "F", false, "Print the file url")
	receiveTextCmd.Flags().BoolP("printAudioUrl", "A", false, "Print the audio url")

	viper.BindPFlag("token", receiveTextCmd.Flags().Lookup("token"))
	viper.BindPFlag("chatId", receiveTextCmd.Flags().Lookup("chatId"))
	viper.BindPFlag("receive.messageCounter", receiveTextCmd.Flags().Lookup("messageCounter"))
	viper.BindPFlag("receive.sync", receiveTextCmd.Flags().Lookup("sync"))
	viper.BindPFlag("receive.printChatId", receiveTextCmd.Flags().Lookup("printChatId"))
	viper.BindPFlag("receive.printMessageId", receiveTextCmd.Flags().Lookup("printMessageId"))
	viper.BindPFlag("receive.printTimestampUnix", receiveTextCmd.Flags().Lookup("printTimestampUnix"))
	viper.BindPFlag("receive.printTimestampHuman", receiveTextCmd.Flags().Lookup("printTimestampHuman"))
	viper.BindPFlag("receive.printPhotoUrl", receiveTextCmd.Flags().Lookup("printPhotoUrl"))
	viper.BindPFlag("receive.printFileUrl", receiveTextCmd.Flags().Lookup("printFileUrl"))
	viper.BindPFlag("receive.printAudioUrl", receiveTextCmd.Flags().Lookup("printAudioUrl"))
}

func validateArgsReceive(cmd *cobra.Command, args []string) error {

	cfg := ReceiveConfig{
		token:               viper.GetString("token"),
		chatId:              viper.GetInt("chatId"),
		messageCounter:      viper.GetInt("receive.messageCounter"),
		sync:                viper.GetBool("receive.sync"),
		printChatId:         viper.GetBool("receive.printChatId"),
		printMessageId:      viper.GetBool("receive.printMessageId"),
		printTimestampUnix:  viper.GetBool("receive.printTimestampUnix"),
		printTimestampHuman: viper.GetBool("receive.printTimestampHuman"),
		printPhotoUrl:       viper.GetBool("receive.printPhotoUrl"),
		printFileUrl:        viper.GetBool("receive.printFileUrl"),
		printAudioUrl:       viper.GetBool("receive.printAudioUrl"),
	}

	//Validate the token
	if cfg.token == "" {
		return fmt.Errorf("No token provided")
	}

	//Validate the chat ID
	if cfg.chatId != 0 && len(strconv.Itoa(cfg.chatId)) != 9 {
		return fmt.Errorf("Wrong chat ID provided")
	}

	//No need to validate the messageCounter, sync, printChatId, printMessageId, printTimestampUnix, printTimestampHuman

	//Send config into the context
	cmd.SetContext(context.WithValue(cmd.Context(), ReceiveConfig{}, cfg))

	return nil
}

func receiveMessage(cmd *cobra.Command, args []string) error {

	cfg := cmd.Context().Value(ReceiveConfig{}).(ReceiveConfig)

	//Create a context
	bgCtx, cancel := context.WithCancel(context.Background())

	//Create the handler
	defaultHandler := func(ctx context.Context, tgBot *bot.Bot, update *models.Update, cancelFunc context.CancelFunc) {

		counter := 0

		//No message exist
		if update.Message == nil {
			return
		}

		//If don't want to sync discard old message
		if int64(update.Message.Date) < time.Now().Unix() && !cfg.sync {
			return
		}

		//Listen only for the specified chat ID
		if update.Message.Chat.ID != int64(cfg.chatId) && cfg.chatId != 0 {
			return
		}
		//Create an empty message that will be filled by the functions
		outputMessage := "|"

		//Append the Date and Time
		if cfg.printTimestampHuman {
			outputMessage += fmt.Sprintf("DATE:%s|", time.Unix(int64(update.Message.Date), 0))
		} else if cfg.printTimestampUnix {
			outputMessage += fmt.Sprintf("DATE:%d|", update.Message.Date)
		}

		//Append Chat ID
		if cfg.printChatId {
			outputMessage += fmt.Sprintf("CHAT_ID:%d|", update.Message.Chat.ID)
		}

		//Append Message ID
		if cfg.printChatId {
			outputMessage += fmt.Sprintf("MESSAGE_ID:%d|", update.Message.ID)
		}

		if cfg.printPhotoUrl || cfg.printFileUrl || cfg.printAudioUrl {
			var fileID string
			var fileName string

			fileIsPresent := false

			if cfg.printPhotoUrl {
				if update.Message.Photo != nil && cfg.printPhotoUrl {
					fileIsPresent = true
					//Get highest resolution photo
					fileID = update.Message.Photo[len(update.Message.Photo)-1].FileID
					fileName = ""
				}
			}

			if cfg.printFileUrl {
				if update.Message.Document != nil && cfg.printFileUrl {
					fileIsPresent = true
					fileID = update.Message.Document.FileID
					fileName = update.Message.Document.FileName
				}
			}

			if cfg.printAudioUrl {
				if update.Message.Audio != nil {
					fileIsPresent = true
					fileID = update.Message.Audio.FileID
					fileName = update.Message.Audio.FileName
				} else if update.Message.Voice != nil {
					fileIsPresent = true
					fileID = update.Message.Voice.FileID
					fileName = ""
				}
			}

			if fileIsPresent {
				// Get file info from Telegram API
				file, err := tgBot.GetFile(ctx, &bot.GetFileParams{FileID: fileID})
				if err != nil {
					fmt.Println("Error file not valid")
					return
				}

				//Append image path
				outputMessage += fmt.Sprintf("URL:https://api.telegram.org/file/bot%s/%s|", cfg.token, file.FilePath)

				if fileName != "" {
					outputMessage += "FILE_NAME:" + fileName + "|"
				}

				outputMessage += "FILE_CAPTION:" + update.Message.Caption + "|"

			}
		}

		//Receive message
		if update.Message.Text != "" { //Handle text message
			//Append message
			outputMessage += "TEXT:" + update.Message.Text + "|"
		}

		//Print out complete message
		fmt.Println(outputMessage)

		//Increase the counter only if user want a cuntdown
		if cfg.messageCounter != 0 {
			counter++
		}

		//Check if counter has reach the user value
		if counter >= cfg.messageCounter && cfg.messageCounter != 0 {
			//Close the bot
			tgBot.Close(ctx)

			//Cancel the Context
			cancelFunc()
		}
	}

	opts := []bot.Option{
		//Link the handler to the bot
		bot.WithDefaultHandler(func(ctx context.Context, b *bot.Bot, update *models.Update) {
			//Pass the param from the default handler + the context cancellation function
			defaultHandler(ctx, b, update, cancel)
		}),
		//Redirect bot library/API errors
		bot.WithErrorsHandler(func(err error) {}),
	}

	//Create the bot
	tgBot, err := bot.New(cfg.token, opts...)
	if nil != err {
		return (err)
	}

	//Start the bot
	tgBot.Start(bgCtx)

	//Close context
	bgCtx.Done()

	return nil
}
