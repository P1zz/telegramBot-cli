package output

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

const (
	FormatText  = "text"
	FormatJSONL = "jsonl"
)

type Message struct {
	Date        string `json:"date,omitempty"`
	ChatID      int64  `json:"chat_id,omitempty"`
	MessageID   int    `json:"message_id,omitempty"`
	FilePath    string `json:"file_path,omitempty"`
	FileName    string `json:"file_name,omitempty"`
	FileCaption string `json:"file_caption,omitempty"`
	Text        string `json:"text,omitempty"`
}

func ValidateFormat(format string) error {
	if format != FormatText && format != FormatJSONL {
		return fmt.Errorf("output must be %q or %q", FormatText, FormatJSONL)
	}
	return nil
}

func Write(w io.Writer, format string, message Message) error {
	if format == FormatJSONL {
		return json.NewEncoder(w).Encode(message)
	}
	_, err := fmt.Fprintln(w, formatText(message))
	return err
}

func formatText(message Message) string {
	fields := make([]string, 0, 7)
	appendField := func(name, value string) {
		if value != "" {
			fields = append(fields, name+":"+escape(value))
		}
	}

	appendField("DATE", message.Date)
	if message.ChatID != 0 {
		appendField("CHAT_ID", fmt.Sprint(message.ChatID))
	}
	if message.MessageID != 0 {
		appendField("MESSAGE_ID", fmt.Sprint(message.MessageID))
	}
	appendField("FILE_PATH", message.FilePath)
	appendField("FILE_NAME", message.FileName)
	appendField("FILE_CAPTION", message.FileCaption)
	appendField("TEXT", message.Text)
	return "|" + strings.Join(fields, "|") + "|"
}

func escape(value string) string {
	replacer := strings.NewReplacer("\\", "\\\\", "|", "\\|", "\r", "\\r", "\n", "\\n")
	return replacer.Replace(value)
}
