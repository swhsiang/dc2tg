# dc2tg
A bot that copies Discord messages and post them to Telegram channels

## Features

* Transfer messages from Discord channels to Telegram channels
* Transfer messages from a Telegram channel to another Telegram channel

## Quick Start

* Remember to update `dev.env` file first.
* Create dev.env file by following the below format
```
TG_APP_TOKEN="TOKEN"
DC_BOT_TOKEN="TOKEN"
DC_USER_EMOJI_LIST="LIST OF DC USERNAME (NO NUMBER)"
# DEFAULT EMOJI
TRIGGERED_EMOJI="<:🚀:>"
# search how to get Telegram group id
TG_TARGET_CHANNEL_ID=-1000000000
```