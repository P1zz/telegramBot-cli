package cmd

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/spf13/cobra"
)

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
	receiveTextCmd.Flags().BoolP("sync", "s", false, "Sync old messages sended while the bot was not running")
	receiveTextCmd.Flags().BoolP("printChatId", "r", false, "Print the chat ID")
	receiveTextCmd.Flags().BoolP("printMessageId", "M", false, "Print the message ID of each message")
	receiveTextCmd.Flags().BoolP("printTimestampUnix", "U", false, "Print the datetime UNIX")
	receiveTextCmd.Flags().BoolP("printTimestampHuman", "H", false, "Print the datetime human readable")
	receiveTextCmd.Flags().BoolP("printPhotoUrl", "P", false, "Print the photo url")
	receiveTextCmd.Flags().BoolP("printFileUrl", "F", false, "Print the file url")
	receiveTextCmd.Flags().BoolP("printAudioUrl", "A", false, "Print the audio url")
}

func receiveMessage(cmd *cobra.Command, args []string) error {
	token, _ := cmd.Flags().GetString("token")
	chatId, _ := cmd.Flags().GetInt("chatId")
	maxMessages, _ := cmd.Flags().GetInt("messageCounter")
	counter := 0
	sync, _ := cmd.Flags().GetBool("sync")
	wantChatId, _ := cmd.Flags().GetBool("printChatId")
	wantMessageId, _ := cmd.Flags().GetBool("printMessageId")
	wantTimestampUnix, _ := cmd.Flags().GetBool("printTimestampUnix")
	wantTimestampHuman, _ := cmd.Flags().GetBool("printTimestampHuman")
	wantPhoto, _ := cmd.Flags().GetBool("printPhotoUrl")
	wantFile, _ := cmd.Flags().GetBool("printFileUrl")
	wantAudio, _ := cmd.Flags().GetBool("printAudioUrl")

	//Create a context
	bgCtx, cancel := context.WithCancel(context.Background())

	//Create the handler
	defaultHandler := func(ctx context.Context, tgBot *bot.Bot, update *models.Update, cancelFunc context.CancelFunc) {

		//No message exist
		if update.Message == nil {
			return
		}

		//If don't want to sync discard old message
		if int64(update.Message.Date) < time.Now().Unix() && !sync {
			return
		}

		//Listen only for the specified chat ID
		if update.Message.Chat.ID != int64(chatId) && chatId != 0 {
			return
		}
		//Create an empty message that will be filled by the functions
		outputMessage := "|"

		//Append the Date and Time
		if wantTimestampHuman {
			outputMessage += fmt.Sprintf("DATE:%s|", time.Unix(int64(update.Message.Date), 0))
		} else if wantTimestampUnix {
			outputMessage += fmt.Sprintf("DATE:%d|", update.Message.Date)
		}

		//Append Chat ID
		if wantChatId {
			outputMessage += fmt.Sprintf("CHAT_ID:%d|", update.Message.Chat.ID)
		}

		//Append Message ID
		if wantMessageId {
			outputMessage += fmt.Sprintf("MESSAGE_ID:%d|", update.Message.ID)
		}

		if wantPhoto || wantFile || wantAudio {
			var fileID string
			var fileName string

			fileIsPresent := false

			if wantPhoto {
				if update.Message.Photo != nil && wantPhoto {
					fileIsPresent = true
					//Get highest resolution photo
					fileID = update.Message.Photo[len(update.Message.Photo)-1].FileID
					fileName = ""
				}
			}

			if wantFile {
				if update.Message.Document != nil && wantFile {
					fileIsPresent = true
					fileID = update.Message.Document.FileID
					fileName = update.Message.Document.FileName
				}
			}

			if wantAudio {
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
				outputMessage += fmt.Sprintf("URL:https://api.telegram.org/file/bot%s/%s|", token, file.FilePath)

				if fileName != "" {
					outputMessage += "FILE_NAME:" + fileName + "|"
				}

				outputMessage += "FILE_CAPTION:" + update.Message.Caption + "|"

			}
		} else if update.Message.Text != "" { //Handle text message
			//Append message
			outputMessage += "TEXT:" + update.Message.Text + "|"
		}

		//Print out complete message
		fmt.Println(outputMessage)

		//Increase the counter only if user want a cuntdown
		if maxMessages != 0 {
			counter++
		}

		//Check if counter has reach the user value
		if counter >= maxMessages && maxMessages != 0 {
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
	}

	//Create the bot
	tgBot, err := bot.New(token, opts...)
	if nil != err {
		return (err)
	}

	//Start the bot
	tgBot.Start(bgCtx)

	//Close context
	bgCtx.Done()

	return nil
}

func validateArgsReceive(cmd *cobra.Command, args []string) error {
	//Validate the token
	token, _ := cmd.Flags().GetString("token")
	if token == "" {
		return fmt.Errorf("no token provided")
	}

	//Validate the chat ID
	chatId, _ := cmd.Flags().GetInt("chatId")
	if chatId != 0 && len(strconv.Itoa(chatId)) != 9 {
		return fmt.Errorf("wrong chat ID provided")
	}

	//No need to validate the messageCounter, sync, printChatId, printMessageId, printTimestampUnix, printTimestampHuman

	return nil
}
