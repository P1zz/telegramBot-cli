package cmd

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/go-telegram/bot"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type DeleteConfig struct {
	token     string
	chatId    int
	messageId int
}

var deleteTextCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete message",
	Long:  "Delete a message",
	//Link the validation function to the validateArgsDelete
	Args: validateArgsDelete,
	//Link the function with the capabilities of returning an error
	Run: deleteMessage,
}

func init() {
	rootCmd.AddCommand(deleteTextCmd)

	deleteTextCmd.Flags().StringP("token", "t", "", "Token from bot fathers")
	deleteTextCmd.Flags().IntP("chatId", "c", 0, "ID of the chat, leave blank or set 0 if you want to listen all chats")
	deleteTextCmd.Flags().IntP("messageId", "i", 0, "ID of the message you wan't to delete")

	viper.BindPFlag("token", deleteTextCmd.Flags().Lookup("token"))
	viper.BindPFlag("chatId", deleteTextCmd.Flags().Lookup("chatId"))
	viper.BindPFlag("delete.messageId", deleteTextCmd.Flags().Lookup("messageId"))
}

func validateArgsDelete(cmd *cobra.Command, args []string) error {

	cfg := DeleteConfig{
		token:     viper.GetString("token"),
		chatId:    viper.GetInt("chatId"),
		messageId: viper.GetInt("delete.messageId"),
	}

	//Validate the token
	if cfg.token == "" {
		return fmt.Errorf("No token provided")
	}

	//Validate the chat ID
	if cfg.chatId != 0 && len(strconv.Itoa(cfg.chatId)) != 9 {
		return fmt.Errorf("Wrong chat ID provided")
	}

	//No need to validate the message ID

	//Send config into the context
	cmd.SetContext(context.WithValue(cmd.Context(), DeleteConfig{}, cfg))

	return nil
}

func deleteMessage(cmd *cobra.Command, args []string) {

	cfg := cmd.Context().Value(DeleteConfig{}).(DeleteConfig)

	//Create a context
	bgCtx := context.Background()

	//Create the bot
	tgBot, err := bot.New(cfg.token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while creating the bot instance\n")
		os.Exit(1)
	}

	//Populate parameters
	parameters := &bot.DeleteMessageParams{
		ChatID:    cfg.chatId,
		MessageID: cfg.messageId,
	}

	//Delete message
	res, err := tgBot.DeleteMessage(bgCtx, parameters)
	if !res && err != nil {
		fmt.Fprintf(os.Stderr, "Error while deleting text message\n")
		os.Exit(2)
	}

	//Close context
	bgCtx.Done()
}
