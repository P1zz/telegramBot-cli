[![Telegram][TELEGRAM_badge]][TELEGRAM_url] [![Bash][BASH_badge]][BASH_URL] [![GO][GO_badge]][GO_url] [![License CC0][LICENSE_badge]][LICENSE_url]

# Brief
A CLI tool for interacting on telegram as a bot written in golang

- [Brief](#brief)
- [Overview](#overview)
- [Build](#build)
- [Usage](#usage)
- [Template config](#template-config)
- [Roadmap](#roadmap)
- [License](#license)

# Overview
This tool is designed to allow easy integration of telegram BOT functionality from shell scripts. <br>
Built With:
* Written entirely in [GO][GO_url]
* [Telegram library](https://github.com/go-telegram/bot) for interacting with the API.
* [Cobra framework](https://github.com/spf13/cobra) to manage CLI inputs.
* [Viper framework](https://github.com/spf13/viper) to manage config.

# Build
```go
go build telegramBot-cli.go
```

# Usage
Hint: All of the commands has ```--help``` parameter.
```bash
telegramBot-cli [ send | receive | edit | delete ]  parameters...
```
# Template config
The default path is ./ and the name is config.toml but you can specify both later.
```toml
token = "YourTelegramBotToken"
chatId = 1234567890

[send]
message = "Message"
markdownV2 = false
fileHasSpoiler = false
getTheMessageId = false
filePath = "path/to/image.jpg"
fileTimeout	= 10 #Timeout in seconds for sending a file
fileIsImage	= false
fileIsVideo	= false
replyChatID	= 0
replyMessageID = 0

[receive]
messageCounter = 0 #Numer of messages to receive, leave blank or set 0 for continuous receiving 
sync = false #Sync old messages sended while the bot was not running
printChatId = true
printMessageId = true
#Mutually exclusive
printTimestampUnix = true
printTimestampHuman = false
#URLs to get the actual received files
printPhotoUrl = false
printFileUrl = false
printAudioUrl = false

[edit]
oldMessageId = 123
newMessage = "NewMessage!"

[delete]
messageId = 123
```

# Roadmap
- [x] Send
    - [x] Message
        - [x] With markdown V2
        - [x] Print the ID of the message
        - [x] With spoiler
        - [ ] Without ringtone sound
        - [ ] Private message (Not forwardable)
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
    - [ ] User defined separator
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
    - [x]  With message Id
- [x] Edit
    - [x] Text
- [ ] Generate a template config

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
