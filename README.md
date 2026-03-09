[![Telegram][TELEGRAM_badge]][TELEGRAM_url] [![Bash][BASH_badge]][BASH_URL] [![GO][GO_badge]][GO_url] [![License CC0][LICENSE_badge]][LICENSE_url]
# Brief
A CLI tool for interacting on telegram as a bot written in golang

- [Brief](#brief)
- [Overview](#overview)
    - [Roadmap](#roadmap)
          - [If you like this repo star and share it!](#if-you-like-this-repo-star-and-share-it)
- [Build](#build)
- [Usage](#usage)
          - [Hint: All of the functions below has ```--help``` parameter.](#hint-all-of-the-functions-below-has---help-parameter)
  - [Send](#send)
  - [Receive](#receive)
  - [Delete](#delete)
  - [Edit](#edit)
- [License](#license)

# Overview
This tool is designed to allow easy integration of telegram BOT functionality from shell scripts. <br>
Built With:
* Written entirely in [GO][GO_url]
* [Telegram library](https://github.com/go-telegram/bot) for interacting with the API.
* [Cobra framework](https://github.com/spf13/cobra) to manage CLI inputs.
### Roadmap
- [x] Send
    - [x] Message
        - [x] With markdown V2
        - [x] Print the ID of the message
    - [x] Image
        - [x] With spoiler
    - [x] Video
        - [x] With spoiler
    - [ ] Audio
    - [ ] Voice
    - [ ] File
    - [x] Reply
        - [x] Same chat
        - [x] Other chat
- [x] Receive
    - [x] Text
    - [x] Audio/Voice
    - [x] Photo
    - [x] Files
    - [x] Receive n messages
    - [ ] Receive messages within n time
    - [x] Receive messages continuously
    - [x] Sync messages received while offline
    - [x] Discard received while offline
    - [x] Receive from a specific chat
    - [x] Receive from all chat
    - [x] Get the message id
    - [x] Get the message time
        - [x] Unix
        - [x] Human Readable
- [x] Delete
- [x] Edit
    - [x] Text
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
  -H, --fileHasSpoiler       The file is send with hidden preview
  -p, --filePath string      Path of the image/video to send
  -T, --fileTimeout int      Timeout in seconds for sending a file (default 60)
  -h, --help                 help for send
  -2, --markDownV2           Message text is parsed in markdown v2
  -m, --messageText string   Message text to send
  -i, --pathIsImage          The path is an image to send
  -v, --pathIsVideo          The path is a video to send
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
|DATA|CHAT_ID|MESSAGE_ID|URL|FILE_NAME|FILE_CAPTION|TEXT|

Usage:
  telegram-cli receive [flags]

Flags:
  -c, --chatId int            ID of the chat, leave blank or set 0 if you want to listen all chats
  -h, --help                  help for receive
  -n, --messageCounter int    Numer of messages to receive, leave blank or set 0 for continuous receiving
  -A, --printAudioUrl         Print the audio url
  -C, --printChatId           Print the chat ID
  -F, --printFileUrl          Print the file url
  -M, --printMessageId        Print the message ID of each message
  -P, --printPhotoUrl         Print the photo url
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
  -i, --messageId int   ID of the message you wan't to delete
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
## Edit
Parameters
```
Edit a text message

Usage:
  telegram-cli edit [flags]

Flags:
  -c, --chatId int           ID of the chat, leave blank or set 0 if you want to listen all chats
  -h, --help                 help for edit
  -i, --messageId int        ID of the message you wan't to edit
  -m, --messageText string   Text of the new message
  -t, --token string         Token from bot fathers
```
Launch
```go
go run telegramBot-cli.go edit {parameters}
```
or
```shell
telegramBot-cli edit {parameters}
```
# License
telegramBot-cli repo is under CC0 1.0.

[GO_badge]: https://img.shields.io/badge/Go-00ADD8?logo=Go&logoColor=white&style=for-the-badge
[GO_url]: https://go.dev 

[LICENSE_badge]: https://img.shields.io/npm/l/cc-md?color=blue&style=for-the-badge
[LICENSE_url]: https://creativecommons.org/public-domain/cc0/

[BASH_badge]: https://img.shields.io/badge/Bash-4EAA25?style=for-the-badge&logo=gnubash&logoColor=white
[BASH_URL]: https://wikipedia.org/wiki/Bash

[TELEGRAM_badge]: https://img.shields.io/badge/Telegram-2CA5E0?style=for-the-badge&logo=telegram&logoColor=white
[TELEGRAM_URL]: https://core.telegram.org/
