# telegramBot-cli

[![CI](https://github.com/P1zz/telegramBot-cli/actions/workflows/ci.yml/badge.svg)](https://github.com/P1zz/telegramBot-cli/actions/workflows/ci.yml)
[![Go](https://img.shields.io/badge/Go-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![License: CC0](https://img.shields.io/badge/license-CC0--1.0-blue)](LICENSE)

`telegramBot-cli` exposes common Telegram Bot API operations as shell-friendly commands. It can send, receive, edit, and delete messages without requiring a dedicated bot application.

## Features

- Send text, images, and videos, including captions, spoilers, replies, and MarkdownV2.
- Receive text and media metadata continuously, for a fixed duration, or for a fixed message count.
- Produce escaped text output or machine-readable JSON Lines.
- Edit and delete messages.
- Read options from flags, environment variables, or a TOML config file.
- Shut down cleanly on Ctrl+C and termination signals.

## Requirements

- Go 1.26 or newer to build from source.
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

Edit and delete a message:

```bash
telegram-cli edit --chat-id 123456789 --message-id 42 --message 'Corrected text'
telegram-cli delete --chat-id 123456789 --message-id 42
```

Run `telegram-cli --help` or `telegram-cli <command> --help` for the complete flag reference.

## Configuration

The default config path is `./config.toml`. Select another file with `--config` or `-C`. A missing default file is allowed.

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

## Receive output

Text output retains the shell-friendly field format. Pipes, backslashes, and line breaks inside values are escaped:

```text
|DATE:2026-07-22T14:30:00+02:00|CHAT_ID:123456789|MESSAGE_ID:42|TEXT:hello\|world|
```

For scripts, `--output jsonl` is recommended:

```json
{"chat_id":123456789,"message_id":42,"text":"hello | world"}
```

Media options return Telegram's `file_path`, not a complete download URL. Telegram download URLs contain the bot token and are therefore unsafe to print to logs.

By default, `receive` ignores queued messages older than the command startup time. Use `--sync` to include them. Use `--chat-id 0` or omit the flag to receive from every chat.

## Development

```bash
make fmt-check
make test
make vet
```

Create cross-platform release binaries with checksums:

```bash
make release VERSION=1.0.0
```

The command layer lives in `cmd/`, while reusable validation and output formatting live in `internal/`. Unit tests cover validation, receive filtering and media selection, and both output formats. CI runs formatting, race-enabled tests, vet, and a build.

## Roadmap

- Send audio, voice messages, and generic files.
- Add notification and content-protection options.
- Support a configurable text separator while preserving JSON Lines as the stable automation format.
- Generate a documented config template directly from the CLI.

## License

This project is released under [CC0 1.0](LICENSE).
