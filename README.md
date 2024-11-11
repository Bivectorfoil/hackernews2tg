# Hack news to telegram channel

Every day telegram bot pickup the first 5 Hack news topic and push them to a private telegram channel

## Pre condistions

- A VPS which can be access telegram bot api server, Ubuntu is recommend
- Go toolchains
- A telegram bot
- A telegram channel (private or public, ATTENTIONS: different channel has different id format, eg: private: -1001234567890, publich: yourchannel_id, for private channle, "-100" is prefix)


## Usage

1. create telegram bot with @BotFather bot (if not have one)

2. got the bot token

3. git clone the repo

4. create .env file with format: TOKEN=<your_bot_token>, CHANNEL_ID=<your_channel_id>

5. run build.sh (or build.ps1 on Windows) file to build the binary file

6. previous step will generate a binary file in dist folder, contains the binary file named `newsboy`

7. run `./dist/newsboy` (or `./dist/newsboy.exe` on Windows) to start the bot

## Deploy

1. run build.sh (or build.ps1 on Windows) file to build the binary file

2. write a service for the `newboy` program, enabled it with systemctl

## TODO

- [ ] linux platform build script

- [ ] linux servce file

- [ ] github ci/cd



## LICENSE

MIT
