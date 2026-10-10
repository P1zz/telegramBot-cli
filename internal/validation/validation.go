package validation

import (
	"errors"
	"fmt"
	"os"
)

func Token(token string) error {
	if token == "" {
		return errors.New("bot token is required")
	}
	return nil
}

func RequiredChatID(id int64) error {
	if id == 0 {
		return errors.New("chat ID is required")
	}
	return nil
}

func MessageID(id int) error {
	if id <= 0 {
		return errors.New("message ID must be greater than zero")
	}
	return nil
}

func NonNegative(name string, value int) error {
	if value < 0 {
		return fmt.Errorf("%s must not be negative", name)
	}
	return nil
}

func Positive(name string, value int) error {
	if value <= 0 {
		return fmt.Errorf("%s must be greater than zero", name)
	}
	return nil
}

func ReadableRegularFile(path string) error {
	if path == "" {
		return errors.New("file path is required")
	}
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("access file %q: %w", path, err)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("file %q is not a regular file", path)
	}
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("read file %q: %w", path, err)
	}
	return file.Close()
}
