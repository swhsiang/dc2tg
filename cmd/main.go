package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	godotenv "github.com/joho/godotenv"

	dcapi "github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/session"
)

// change it to woody's channel id
const TG_SUBSCRIBED_CHANNEL_IDS = -1001665500012

// change it to 168 group chat id
const TG_TARGET_CHANNEL_ID_ENV = "TG_TARGET_CHANNEL_ID"
const TG_APP_TOKEN_ENV = "TG_APP_TOKEN"
const DC_BOT_TOKEN_ENV = "DC_BOT_TOKEN"
const DC_USER_EMOJI_LIST_ENV = "DC_USER_EMOJI_LIST"
const TRIGGERED_EMOJI_ENV = "TRIGGERED_EMOJI"

const dc2tgMsgTemplate = `
	[%s]

	%s

	by	%s
	`

var (
	dcClient            *dcapi.Client
	tgClient            *tgbotapi.BotAPI
	dcUserEmojiList     []string
	TRIGGERED_EMOJI     string
	TG_LISTENER_CHAT_ID int64
)

func initClient(token string) {
	dcClient = dcapi.NewClient("Bot " + token)
	if dcClient == nil {
		log.Fatalln("Unable to create a Discord client instance.")
	}

	var err error
	tgClient, err = tgbotapi.NewBotAPI(os.Getenv(TG_APP_TOKEN_ENV))
	if err != nil {
		log.Fatalln("Unable to create a Telegram client instance.")
	}

	tgClient.Debug = true

	log.Printf("[Telegram] Authorized on account %s", tgClient.Self.UserName)

	temp := os.Getenv(DC_USER_EMOJI_LIST_ENV)
	dcUserEmojiList = strings.Split(temp, ",")
	if len(dcUserEmojiList) == 0 {
		log.Fatal("No dc user is granted permission to transfer message to telegram group. check env file.")
	}

	TRIGGERED_EMOJI = os.Getenv(TRIGGERED_EMOJI_ENV)
	if TRIGGERED_EMOJI == "" {
		log.Fatal("No emoji selected. Check env file")
	}

	TG_LISTENER_CHAT_ID, err = strconv.ParseInt(os.Getenv(TG_TARGET_CHANNEL_ID_ENV), 0, 64)
	if err != nil || TG_LISTENER_CHAT_ID == 0 || TG_LISTENER_CHAT_ID > 0 {
		log.Fatal("Incorrect TG group ID. Check env file")
	}
}

func craftTGMessage(dcChannel, dcMessage, dcUserName string) string {
	return fmt.Sprintf(dc2tgMsgTemplate, dcChannel, dcMessage, dcUserName)
}

func dc(token string) {

	s := session.New("Bot " + token)
	s.AddHandler(func(c *gateway.MessageReactionAddEvent) {

		msg, err := dcClient.Message(c.ChannelID, c.MessageID)
		if err != nil {
			log.Println("Failed to get message id ", c.MessageID, " from channel ", c.ChannelID, " err:", err.Error())
			return
		}

		channelInstance, err := dcClient.Channel(c.ChannelID)
		if err != nil {
			log.Println("Failed to get channel ", c.ChannelID, " err:", err.Error())
			return
		}

		log.Println(c.Member.User.Username, "sent", c.Emoji, "to message ", msg.Content)

		if c.Emoji.String() == TRIGGERED_EMOJI {
			for allowedUser := range dcUserEmojiList {
				if dcUserEmojiList[allowedUser] == c.Member.User.Username {

					msgInstance := craftTGMessage(channelInstance.Name, msg.Content, msg.Author.Username)

					log.Println("Sending message ", msgInstance, " to ", TG_LISTENER_CHAT_ID)

					tgMsg := tgbotapi.NewMessage(TG_LISTENER_CHAT_ID, msgInstance)
					tgClient.Send(tgMsg)
					break
				}
			}
		}
	})

	// Add the needed Gateway intents.
	s.AddIntents(gateway.IntentGuildMessageReactions)

	if err := s.Open(context.Background()); err != nil {
		log.Fatalln("Failed to connect:", err)
	}
	defer s.Close()

	u, err := s.Me()
	if err != nil {
		log.Fatalln("Failed to get myself:", err)
	}

	log.Println("Started as", u.Username)

	select {}
}

var quit = make(chan struct{})

func main() {
	const envFileName = "dev.env"
	err := godotenv.Load(envFileName)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var token = os.Getenv(DC_BOT_TOKEN_ENV)
	if token == "" {
		log.Fatalln("No $BOT_TOKEN given.")
	}

	initClient(token)

	go dc(token)

	<-quit
}
