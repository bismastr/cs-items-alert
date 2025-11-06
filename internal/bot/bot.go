package bot

import (
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	session *discordgo.Session
}

func NewBot() *Bot {
	session, err := discordgo.New(fmt.Sprintf("Bot %v", os.Getenv("DISCORD_BOT_TOKEN")))
	if err != nil {
		log.Printf("Error discordgo session %v", err)
	}

	session.Identify.Intents = discordgo.IntentsAll
	err = session.Open()
	if err != nil {
		log.Printf("Cannot create discordgo session %v", err)
	}

	return &Bot{
		session: session,
	}
}

func (b *Bot) SendMessageToChannel(channelId, content string) {
	_, err := b.session.ChannelMessageSend(channelId, content)
	if err != nil {
		log.Printf("Error sending message %v", err)
	}
}
