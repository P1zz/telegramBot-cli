package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestWriteTextEscapesDelimiters(t *testing.T) {
	var buffer bytes.Buffer
	err := Write(&buffer, FormatText, Message{ChatID: 42, Text: "hello|world\nnext"})
	if err != nil {
		t.Fatal(err)
	}
	want := "|CHAT_ID:42|TEXT:hello\\|world\\nnext|\n"
	if buffer.String() != want {
		t.Fatalf("text output = %q, want %q", buffer.String(), want)
	}
}

func TestWriteJSONL(t *testing.T) {
	var buffer bytes.Buffer
	message := Message{ChatID: -1001, MessageID: 7, Text: "hello"}
	if err := Write(&buffer, FormatJSONL, message); err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(buffer.String(), "\n") {
		t.Fatal("JSONL output does not end with a newline")
	}
	var decoded Message
	if err := json.Unmarshal(buffer.Bytes(), &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded != message {
		t.Fatalf("decoded message = %#v, want %#v", decoded, message)
	}
}

func TestValidateFormat(t *testing.T) {
	if err := ValidateFormat(FormatJSONL); err != nil {
		t.Fatal(err)
	}
	if err := ValidateFormat("yaml"); err == nil {
		t.Fatal("ValidateFormat accepted yaml")
	}
}
