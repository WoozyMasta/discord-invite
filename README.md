# Discord Invite Link Generator

![logo]

The Discord Invite Link Generator is a application designed to create and
manage invite links for a specified Discord channel. It provides a simple
HTTP server that generates unique invite links based on configurable
parameters.

You can configure maximum age, usage limits, and uniqueness of generated
invite links.

## Configuration

The application can be configured using command-line flags or
environment variables.

```txt
Usage:
  discord-invite [OPTIONS]

Invite links generator for Discord channel.

Application Options:
  -c, --channel_id= Discord channel ID [$DINVITE_CHANNEL_ID]
  -t, --bot_token=  Discord bot token [$DINVITE_BOT_TOKEN]
  -l, --listen=     Address to listen on (default: :8080) [$DINVITE_LISTEN]
      --log-level=  Log level (default: info) [$DINVITE_LOG_LEVEL]
  -a, --max_age=    Invite max age in seconds (default: 3600) [$DINVITE_MAX_AGE]
  -u, --max_uses=   Invite max uses (default: 1) [$DINVITE_MAX_USES]
  -x, --unique      Make every invite unique [$DINVITE_UNIQUE]
  -v, --version     Show version, commit, and build time.

Help Options:
  -h, --help        Show this help message
```

### Get token

Create a bot at <https://discord.com/developers/applications> add it to your
server, make sure it can create channel invites. Copy the channel ID and get
a token from the developer portal in the bot tab.

## Usage

Run the application with the required flags:

```bash
./invite-generator --channel_id your_channel_id --bot_token your_bot_token
```

You can also use environment variables:

```bash
DINVITE_CHANNEL_ID=your_channel_id DINVITE_BOT_TOKEN=your_bot_token ./invite-generator
```

Access the invite link by navigating to <http://localhost:8080> (or your
specified listen address). The server will generate a unique Discord invite
link based on your configuration and redirect you accordingly.

## Installation

You can download the latest version of the programme by following the links:

* [MacOS arm64][]
* [MacOS amd64][]
* [Linux i386][]
* [Linux amd64][]
* [Linux arm][]
* [Linux arm64][]
* [Windows i386][]
* [Windows amd64][]
* [Windows arm64][]

For Linux you can also use the command

```bash
curl -#SfLo /usr/bin/discord-invite \
  https://github.com/WoozyMasta/discord-invite/releases/latest/download/discord-invite-linux-amd64
chmod +x /usr/bin/discord-invite
discord-invite -h && discord-invite -v
```

### Container Image

The images are published to two container registries:

* [`docker pull ghcr.io/woozymasta/discord-invite:latest`][ghcr]
* [`docker pull docker.io/woozymasta/discord-invite:latest`][docker]

Quick start:

```bash
# Pull the image
docker pull ghcr.io/woozymasta/discord-invite:latest
# Run the container with pass environment variables and exposed port
docker run --name discord-invite -d \
  -p 8080:8080 -e DINVITE_BOT_TOKEN='' -e DINVITE_CHANNEL_ID='' \
  ghcr.io/woozymasta/discord-invite:latest
```

### Systemd service

To run the Discord Invite as a systemd service, use the following example
configuration. This ensures the exporter runs on system startup.

```ini
[Unit]
Description=Invite links generator for Discord channel
Documentation=https://github.com/woozymasta/discord-invite
Wants=network-online.target
After=network-online.target

[Service]
EnvironmentFile=-/env/discord-invite.env
Environment="DINVITE_LISTEN=127.0.0.1:8080"
Environment="DINVITE_CHANNEL_ID=your_channel_id"
Environment="DINVITE_BOT_TOKEN=your_bot_token"
ExecStart=/usr/bin/discord-invite
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
```

### Windows service

You can run the Discord Invite using any method that suits you, but it's
recommended to use a Windows service for better management and reliability.

To register the service, assuming the application and configuration are
already downloaded and set up in the `C:\discord-invite` directory,
use the following commands:

```powershell
sc.exe create discord-invite `
  binPath= "C:\discord-invite\discord-invite.exe -c your_channel_id -t your_bot_token" `
  DisplayName= "Discord Invite Link Generator" `
  start= auto

sc.exe start discord-invite
sc.exe query discord-invite
```

Uninstall service

```powershell
sc.exe stop discord-invite
sc.exe query discord-invite
sc.exe delete discord-invite
```

## Alternative option

If you use a web server nginx with lua support, for example such as
OpenResty, you can use a lua script that implements the same logic of
issuing invites directly in the web server. More details in
[lua/README.md](lua/README.md)

## Support me 💖

If you enjoy my projects and want to support further development,
feel free to donate! Every contribution helps to keep the work going.
Thank you!

### Crypto Donations

<!-- cSpell:disable -->
* **BTC**: `1Jb6vZAMVLQ9wwkyZfx2XgL5cjPfJ8UU3c`
* **USDT (TRC20)**: `TN99xawQTZKraRyvPAwMT4UfoS57hdH8Kz`
* **TON**: `UQBB5D7cL5EW3rHM_44rur9RDMz_fvg222R4dFiCAzBO_ptH`
<!-- cSpell:enable -->

Your support is greatly appreciated!

<!-- Links -->
[logo]: winres/icon64.png
[ghcr]: https://github.com/WoozyMasta/discord-invite/pkgs/container/discord-invite
[docker]: https://hub.docker.com/r/woozymasta/discord-invite
[MacOS arm64]: https://github.com/WoozyMasta/discord-invite/releases/latest/download/discord-invite-darwin-arm64 "MacOS arm64 file"
[MacOS amd64]: https://github.com/WoozyMasta/discord-invite/releases/latest/download/discord-invite-darwin-amd64 "MacOS amd64 file"
[Linux i386]: https://github.com/WoozyMasta/discord-invite/releases/latest/download/discord-invite-linux-386 "Linux i386 file"
[Linux amd64]: https://github.com/WoozyMasta/discord-invite/releases/latest/download/discord-invite-linux-amd64 "Linux amd64 file"
[Linux arm]: https://github.com/WoozyMasta/discord-invite/releases/latest/download/discord-invite-linux-arm "Linux arm file"
[Linux arm64]: https://github.com/WoozyMasta/discord-invite/releases/latest/download/discord-invite-linux-arm64 "Linux arm64 file"
[Windows i386]: https://github.com/WoozyMasta/discord-invite/releases/latest/download/discord-invite-windows-386.exe "Windows i386 file"
[Windows amd64]: https://github.com/WoozyMasta/discord-invite/releases/latest/download/discord-invite-windows-amd64.exe "Windows amd64 file"
[Windows arm64]: https://github.com/WoozyMasta/discord-invite/releases/latest/download/discord-invite-windows-arm64.exe "Windows arm64 file"
