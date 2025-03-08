package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	session *discordgo.Session
}

func NewBot(session *discordgo.Session) *Bot {
	session.Identify.Intents = discordgo.IntentsAll
	err := session.Open()
	if err != nil {
		log.Printf("Cannot create discordgo session %v", err)
	}

	return &Bot{
		session: session,
	}
}
