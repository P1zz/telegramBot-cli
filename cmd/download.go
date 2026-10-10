package cmd

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/go-telegram/bot"
	"github.com/spf13/cobra"
)

var downloadHTTPClient = http.DefaultClient

var downloadCmd = &cobra.Command{
	Use: "download", Short: "Download a Telegram file without exposing the bot token", Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		token, err := resolvedToken()
		if err != nil {
			return err
		}
		fileID, _ := cmd.Flags().GetString("file-id")
		if fileID == "" {
			return fmt.Errorf("file ID is required")
		}
		outputPath, _ := cmd.Flags().GetString("output")
		if outputPath == "" {
			return fmt.Errorf("output path is required")
		}
		client, err := newTelegramClient(token)
		if err != nil {
			return fmt.Errorf("create bot: %w", err)
		}
		file, err := client.GetFile(cmd.Context(), &bot.GetFileParams{FileID: fileID})
		if err != nil {
			return fmt.Errorf("get file metadata: %w", err)
		}
		return downloadFile(cmd.Context(), token, file.FilePath, outputPath, cmd.OutOrStdout())
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().String("file-id", "", "Telegram file ID")
	downloadCmd.Flags().StringP("output", "o", "-", "destination path, or - for stdout")
}

func downloadFile(ctx context.Context, token, filePath, outputPath string, stdout io.Writer) (returnErr error) {
	downloadURL := "https://api.telegram.org/file/bot" + url.PathEscape(token) + "/" + filePath
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadURL, nil)
	if err != nil {
		return fmt.Errorf("create download request: %w", err)
	}
	response, err := downloadHTTPClient.Do(req)
	if err != nil {
		return safeError("download file", token, err)
	}
	defer func() { _ = response.Body.Close() }()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("download file: Telegram returned %s", response.Status)
	}
	if outputPath == "-" {
		_, err = io.Copy(stdout, response.Body)
		return err
	}
	destination, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return fmt.Errorf("create destination: %w", err)
	}
	defer func() {
		if closeErr := destination.Close(); returnErr == nil {
			returnErr = closeErr
		}
	}()
	if _, err := io.Copy(destination, response.Body); err != nil {
		return fmt.Errorf("write destination: %w", err)
	}
	return nil
}
