# telegramBot-cli

[![CI](https://github.com/P1zz/telegramBot-cli/actions/workflows/ci.yml/badge.svg)](https://github.com/P1zz/telegramBot-cli/actions/workflows/ci.yml)
[![Go](https://img.shields.io/badge/Go-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![License: CC0](https://img.shields.io/badge/license-CC0--1.0-blue)](LICENSE)

`telegramBot-cli` exposes common Telegram Bot API operations as shell-friendly commands. It can send, receive, edit, and delete messages without requiring a dedicated bot application.

## Features

- Send text, images, videos, documents, audio, and voice messages, including captions, spoilers, replies, silent delivery, content protection, and MarkdownV2.
- Receive text and media metadata continuously, for a fixed duration, or for a fixed message count.
- Produce escaped text output or machine-readable JSON Lines.
- Edit and delete messages.
- Read messages from stdin and emit structured JSON results.
- Download files without exposing the bot token in logs or process output.
- Read options from flags, environment variables, a token file, or a TOML config file.
- Shut down cleanly on Ctrl+C and termination signals.

## Requirements

- Go 1.27 or newer to build from source (the project toolchain is Go 1.27.2).
- A Telegram bot token created with [BotFather](https://core.telegram.org/bots/features#botfather).
- A chat ID where the bot is allowed to perform the requested operation.

Treat the bot token like a password. Do not commit `config.toml`, print the token in logs, or pass it directly on a shared machine where command lines can be inspected. The environment variable is the safest default.

## Installation

Build the current version:

```bash
git clone https://github.com/P1zz/telegramBot-cli.git
cd telegramBot-cli
go build -o telegram-cli ./cmd/telegram-cli
```

Or use the Makefile:

```bash
make test
make build
./builds/telegram-cli version
```

## Quick start

Export the token so it is not repeated in every command:

```bash
export TELEGRAM_CLI_TOKEN='123456:replace-with-your-token'
```

Send a message:

```bash
telegram-cli send --chat-id 123456789 --message 'Deployment completed'
```

Send an image and print the resulting message ID:

```bash
telegram-cli send \
  --chat-id 123456789 \
  --image \
  --file ./diagram.png \
  --message 'Current architecture' \
  --print-message-id
```

Receive the next message as JSON Lines:

```bash
telegram-cli receive --counter 1 --output jsonl
```

Pipe message text from another command:

```bash
printf 'Build completed' | telegram-cli send --chat-id 123456789 --stdin --output json
```

Download a received `file_id` without printing a token-bearing URL:

```bash
telegram-cli download --file-id ABC123 --output ./download.bin
```

Edit and delete a message:

```bash
telegram-cli edit --chat-id 123456789 --message-id 42 --message 'Corrected text'
telegram-cli delete --chat-id 123456789 --message-id 42
```

Run `telegram-cli --help` or `telegram-cli <command> --help` for the complete flag reference.

## Configuration

Configuration is discovered first in the current directory and then in the operating system's user config directory, such as `~/.config/telegram-cli/config.toml`. Select an exact file with `--config` or `-C`. A missing discovered file is allowed, while a missing explicitly selected file is an error.

Generate a safe, documented user config without overwriting an existing one:

```bash
telegram-cli config init
telegram-cli config init --output ./config.toml
```

Use `--force` only when an existing destination should be replaced.

```toml
# Prefer TELEGRAM_CLI_TOKEN instead of storing the token here.
chatId = 123456789

[send]
message = "Hello from telegramBot-cli"
markdownV2 = false
fileHasSpoiler = false
printMessageID = false
filePath = ""
fileTimeout = 60
fileIsImage = false
fileIsVideo = false
replyChatID = 0
replyMessageID = 0

[receive]
counter = 0
timeout = 0
sync = false
output = "text"
printChatId = true
printMessageId = true
printTimestampUnix = false
printTimestampHuman = true
printPhotoPath = false
printFilePath = false
printAudioPath = false

[edit]
oldMessageId = 0
newMessage = ""

[delete]
messageId = 0
```

Configuration precedence is:

1. Command-line flags.
2. Environment variables.
3. TOML configuration.
4. Built-in defaults.

Environment variables start with `TELEGRAM_CLI_`. Dots become underscores; for example, `receive.timeout` maps to `TELEGRAM_CLI_RECEIVE_TIMEOUT`.

For service environments, `--token-file /run/secrets/telegram-token` reads a token from a file. The file must not be accessible by group or other users. `--token-file -` reads it from stdin, but cannot be combined with `send --stdin`. A config file containing a token is subject to the same permission check.

## Receive output

Text output retains the shell-friendly field format. Pipes, backslashes, and line breaks inside values are escaped:

```text
|DATE:2026-07-22T14:30:00+02:00|EVENT:message|CHAT_ID:123456789|MESSAGE_ID:42|TEXT:hello\|world|
```

For scripts, `--output jsonl` is recommended:

```json
{"event":"message","chat_id":123456789,"message_id":42,"text":"hello | world"}
```

Media options return Telegram's `file_path`, not a complete download URL. Telegram download URLs contain the bot token and are therefore unsafe to print to logs.

By default, `receive` ignores queued messages older than the command startup time. Use `--sync` to include them. Use `--chat-id 0` or omit the flag to receive from every chat.

The `event` field distinguishes `message`, `edited_message`, `channel_post`, and `edited_channel_post`. Non-sync polling starts from Telegram's latest update offset so an old queue is discarded efficiently.

## Development

```bash
make fmt-check
make test
make vet
make vuln-check
```

`make vuln-check` runs Go's `govulncheck` with tests and the `e2e` build tag
included, and shows detailed traces. It checks for known vulnerabilities across
the module dependencies. It requires Go 1.27 or newer and downloads the pinned
scanner version automatically when needed.

Run the opt-in Telegram lifecycle test only with a dedicated test bot and chat. It creates, edits, and deletes a temporary message:

```bash
export TELEGRAM_CLI_TOKEN='test-bot-token'
export TELEGRAM_CLI_CHAT_ID='123456789'
make test-e2e
```

Create cross-platform release binaries with checksums:

```bash
make release VERSION=1.0.0
```

The command layer lives in `cmd/`, while reusable validation, Telegram interfaces, and output formatting live in `internal/`. Tests cover API parameter construction, validation, downloading, receive filtering and media selection, and both output formats. CI runs formatting, race-enabled tests with a coverage floor, vet, GolangCI-Lint, `govulncheck`, Gitleaks, and a build. Dependabot tracks Go modules and Actions.

## Shell completion

Cobra provides completion scripts for Bash, Zsh, Fish, and PowerShell:

```bash
# Bash
telegram-cli completion bash > ~/.local/share/bash-completion/completions/telegram-cli

# Zsh
telegram-cli completion zsh > "${fpath[1]}/_telegram-cli"

# Fish
telegram-cli completion fish > ~/.config/fish/completions/telegram-cli.fish

# PowerShell
telegram-cli completion powershell > telegram-cli.ps1
```

## Compatibility and releases

Former camel-case flags such as `--chatId`, `--messageText`, and `--printPhotoUrl` remain available as deprecated aliases for migration. New scripts should use the kebab-case names shown by `--help`.

Pushing a signed or annotated `v*` tag runs the release workflow, builds all supported platforms, packages the binary with documentation and example configuration, publishes checksums, and creates GitHub release notes.

## Roadmap

- Support a configurable text separator while preserving JSON Lines as the stable automation format.

## License

This project is released under [CC0 1.0](LICENSE).
