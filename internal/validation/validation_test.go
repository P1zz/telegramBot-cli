package validation

import "testing"

func TestRequiredChatID(t *testing.T) {
	for _, test := range []struct {
		name    string
		id      int64
		wantErr bool
	}{
		{name: "private chat", id: 123456789},
		{name: "group chat", id: -1001234567890},
		{name: "missing", id: 0, wantErr: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			if got := RequiredChatID(test.id); (got != nil) != test.wantErr {
				t.Fatalf("RequiredChatID(%d) error = %v, wantErr %v", test.id, got, test.wantErr)
			}
		})
	}
}

func TestMessageID(t *testing.T) {
	if err := MessageID(1); err != nil {
		t.Fatalf("MessageID(1): %v", err)
	}
	if err := MessageID(0); err == nil {
		t.Fatal("MessageID(0) returned no error")
	}
}

func TestNumbers(t *testing.T) {
	if err := NonNegative("counter", -1); err == nil {
		t.Fatal("NonNegative accepted -1")
	}
	if err := Positive("timeout", 0); err == nil {
		t.Fatal("Positive accepted zero")
	}
}
