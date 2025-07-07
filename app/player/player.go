package player

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"telebotmusicos/app/userbot"
)

type Player struct {
	bot     *tgbotapi.BotAPI
	userbot *userbot.UserbotClient
}

func New(bot *tgbotapi.BotAPI, userbot *userbot.UserbotClient) *Player {
	return &Player{
		bot:     bot,
		userbot: userbot,
	}
}

func (p *Player) Play(chatID int64, query string) {
	log.Printf("Log: Received play request in chat %d for query: %s", chatID, query)
	
	// Logic to find user in voice chat, get audio, and stream will be added here later.
	// For now, let's just send a confirmation message.

	msg := tgbotapi.NewMessage(chatID, "🎶 Received play request for: "+query)
	p.bot.Send(msg)
}