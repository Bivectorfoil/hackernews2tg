# Hack news to telegram channel

Every day telegram bot pickup the first 5 Hack news topic and push them to a private telegram channel

## Pre condistions

- A VPS which can be access telegram bot api server, Ubuntu is recommend
- Go toolchains
- A telegram bot
- A telegram channel (private or public, ATTENTIONS: different channel has different id format, eg: private: -1001234567890, publich: yourchannel_id, for private channle, "-100" is prefix)

## Development

1. create telegram bot with @BotFather bot (if not have one)

2. got the bot token

3. git clone the repo

4. create .env file with format: TOKEN=<your_bot_token>, CHANNEL_ID=<your_channel_id>

5. run build.sh (or build.ps1 on Windows) file to build the binary file

6. previous step will generate a binary file in dist folder, contains the binary file named `newsboy`

7. run `./dist/newsboy` (or `./dist/newsboy.exe` on Windows) to start the bot

## Installation (On Linux)

1. Clone the repository:
   ```bash
   git clone https://github.com/your-username/newsboy.git
   cd newsboy
   ```

2. Build the project:
   ```bash
   go build -o dist/newsboy
   ```

3. Create and edit configuration file:
   ```bash
   cp .env.example .env
   vim .env
   ```

4. Run installation script:
   ```bash
   sudo deploy/scripts/install.sh
   ```

### Common Service Management Commands

```bash
# Stop service
sudo systemctl stop newsboy

# Restart service
sudo systemctl restart newsboy

# View logs
sudo journalctl -u newsboy
```

## TODO

- [x] linux platform build script
- [x] linux service file
- [x] support any permission user to run linux install script
- [ ] github ci/cd



## LICENSE

MIT
