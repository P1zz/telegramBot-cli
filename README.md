[![Telegram][TELEGRAM_badge]][TELEGRAM_url] [![Bash][BASH_badge]][BASH_URL] [![GO][GO_badge]][GO_url] [![License CC0][LICENSE_badge]][LICENSE_url]
# Brief
A CLI tool for interacting on telegram as a bot written in golang

- [Overview](#overview)
- [Build](#build)
- [Usage](#usage)
  - [Send](#send)
  - [Receive](#receive)
  - [Delete](#delete)
- [License](#license)

# Overview
This tool is designed to allow easy integration of telegram BOT functionality from shell scripts. <br>
Built With:
* Written entirely in [GO][GO_url]
* [Telegram library](https://github.com/go-telegram/bot) for interacting with the API.
* [Cobra framework](https://github.com/spf13/cobra) to manage CLI inputs.
###### If you like this repo star and share it!
# Build
```go
go build telegramBot-cli.go
```
# Usage
###### Hint: All of the functions below has ```--help``` parameter.
## Send
Parameters
```
Send a message in a chat as bot with text or an image

Usage:
  telegram-cli send [flags]

Flags:
  -c, --chatId int           Your chat ID
  -h, --help                 help for send
  -i, --imagePath string     Path of the image to send
  -m, --messageText string   Message text to send
  -M, --printMessageId       Print message id of your message
  -x, --replyChatId int      Chat id you want to reply
  -y, --replyMessageId int   Message id you want to reply
  -t, --token string         Token from bot fathers
```
Launch
```go
go run telegramBot-cli.go send {parameters}
```
or
```shell
telegramBot-cli send {parameters}
```
## Receive
Parameters
```
Receive a message as bot with the pattern below
DATA | CHAT_ID | MESSAGE_ID | MESSAGE

Usage:
  telegram-cli receive [flags]

Flags:
  -c, --chatId int            ID of the chat, leave blank or set 0 if you want to listen all chats
  -h, --help                  help for receive
  -n, --messageCounter int    Numer of messages to receive, leave blank or set 0 for continuous receiving
  -C, --printChatId           Print the chat ID
  -M, --printMessageId        Print the message ID of each message
  -H, --printTimestampHuman   Print the datetime human readable
  -U, --printTimestampUnix    Print the datetime UNIX
  -s, --sync                  Sync old messages sended while the bot was not running
  -t, --token string          Token from bot fathers
```
Launch
```go
go run telegramBot-cli.go receive {parameters}
```
or
```shell
telegramBot-cli receive {parameters}
```
## Delete
Parameters
```
Delete a message

Usage:
  telegram-cli delete [flags]

Flags:
  -c, --chatId int      ID of the chat, leave blank or set 0 if you want to listen all chats
  -h, --help            help for delete
  -m, --messageId int   ID of the message tou wan't to delete
  -t, --token string    Token from bot fathers
```
Launch
```go
go run telegramBot-cli.go delete {parameters}
```
or
```shell
telegramBot-cli delete {parameters}
```
# License
telegramBot-cli repo is under CC0 1.0.

[GO_badge]: https://img.shields.io/badge/Go-00ADD8?logo=Go&logoColor=white&style=for-the-badge
[GO_url]: https://jquery.com 

[LICENSE_badge]: https://img.shields.io/npm/l/cc-md?color=blue&style=for-the-badge
[LICENSE_url]: https://creativecommons.org/public-domain/cc0/

[BASH_badge]: https://img.shields.io/badge/Bash-4EAA25?style=for-the-badge&logo=gnubash&logoColor=white
[BASH_URL]: https://wikipedia.org/wiki/Bash

[TELEGRAM_badge]: https://img.shields.io/badge/Telegram-2CA5E0?style=for-the-badge&logo=telegram&logoColor=white
[TELEGRAM_URL]: https://core.telegram.org/
