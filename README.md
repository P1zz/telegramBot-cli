[![Telegram][TELEGRAM_badge]][TELEGRAM_url] [![Bash][BASH_badge]][BASH_URL] [![GO][GO_badge]][GO_url] [![License CC0][LICENSE_badge]][LICENSE_url]
<br><br>

# Brief
A CLI tool for interacting on telegram as a bot written in golang

- [Brief](#brief)
- [Overview](#overview)
- [Build](#build)
- [Usage](#usage)
- [Roadmap](#roadmap)
- [License](#license)

# Overview
This tool is designed to allow easy integration of telegram BOT functionality from shell scripts. <br>
Built With:
* Written entirely in [GO][GO_url]
* [Telegram library](https://github.com/go-telegram/bot) for interacting with the API.
* [Cobra framework](https://github.com/spf13/cobra) to manage CLI inputs.
* [Viper framework](https://github.com/spf13/viper) to manage config.
<br><br>

# Build
```go
go build telegramBot-cli.go
```
<br><br>

# Usage
Hint: All of the commands has ```--help``` parameter.
```go
telegramBot-cli { send | receive | edit | delete }  parameters
```
<br><br>

# Roadmap
- [x] Send
    - [x] Message
        - [x] With markdown V2
        - [x] Print the ID of the message
        - [x] With spoiler
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
<br><br>

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
