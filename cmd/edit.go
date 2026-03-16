package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-telegram/bot"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type EditConfig struct {
	token        string
	chatId       int
	newMessage   string
	oldMessageId int
}

var editTextCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit message",
	Long:  "Edit a text message",
	//Link the validation function to the validateArgsEdit
	Args: validateArgsEdit,
	//Link the function with the capabilities of returning an error
	RunE: editMessage,
}

func init() {
	rootCmd.AddCommand(editTextCmd)

	editTextCmd.Flags().StringP("token", "t", "", "Token from bot fathers")
	editTextCmd.Flags().IntP("chatId", "c", 0, "ID of the chat, leave blank or set 0 if you want to listen all chats")
	editTextCmd.Flags().IntP("oldMessageId", "i", 0, "ID of the message you wan't to edit")
	editTextCmd.Flags().StringP("newMessage", "m", "", "Text of the new message")

	viper.BindPFlag("token", editTextCmd.Flags().Lookup("token"))
	viper.BindPFlag("chatId", editTextCmd.Flags().Lookup("chatId"))
	viper.BindPFlag("edit.newMessage", editTextCmd.Flags().Lookup("newMessage"))
	viper.BindPFlag("edit.oldMessageId", editTextCmd.Flags().Lookup("oldMessageId"))
}

func validateArgsEdit(cmd *cobra.Command, args []string) error {

	cfg := EditConfig{
		token:        viper.GetString("token"),
		chatId:       viper.GetInt("chatId"),
		newMessage:   viper.GetString("edit.newMessage"),
		oldMessageId: viper.GetInt("edit.oldMessageId"),
	}

	//Validate the token
	if cfg.token == "" {
		return fmt.Errorf("No token provided")
	}

	//Validate the chat ID
	if cfg.chatId != 0 && len(strconv.Itoa(cfg.chatId)) != 9 {
		return fmt.Errorf("Wrong chat ID provided")
	}

	//Validate the text
	if cfg.newMessage == "" {
		return fmt.Errorf("New message not provided")
	}

	//No need to validate the message ID

	//Send config into the context
	cmd.SetContext(context.WithValue(cmd.Context(), EditConfig{}, cfg))

	return nil
}

func editMessage(cmd *cobra.Command, args []string) error {

	cfg := cmd.Context().Value(EditConfig{}).(EditConfig)

	//Create a context
	bgCtx := context.Background()

	//Create the bot
	tgBot, err := bot.New(cfg.token)
	if err != nil {
		return err
	}

	//Populate parameters
	parameters := &bot.EditMessageTextParams{
		ChatID:    cfg.chatId,
		MessageID: cfg.oldMessageId,
		Text:      cfg.newMessage,
	}

	//Edit message
	_, err = tgBot.EditMessageText(bgCtx, parameters)
	if err != nil {
		return err
	}

	//Close context
	bgCtx.Done()

	return nil
}
