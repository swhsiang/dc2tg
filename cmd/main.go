package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	godotenv "github.com/joho/godotenv"
)

// change it to woody's channel id
const TG_SUBSCRIBED_CHANNEL_IDS = -1001665500012

// change it to 168 group chat id
const TG_LISTENER_CHAT_ID = -1001586751727

func telegram() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TG_APP_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		// NOTE forward messages from one channel to another group chat
		if update.ChannelPost != nil {
			msg := tgbotapi.NewForward(TG_LISTENER_CHAT_ID, update.ChannelPost.Chat.ID, update.ChannelPost.MessageID)
			bot.Send(msg)
		}
	}
}

func main() {
	const envFileName = "dev.env"
	err := godotenv.Load(envFileName)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	telegram()
}
