package bot

import (
	"log"
	"strings"

	"telebotmusicos/app/locales"
	"telebotmusicos/app/player"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	Player *player.Player
}

func (h *Handler) HandleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	log.Printf("Log: Received message from [%s]: %s", update.Message.From.UserName, update.Message.Text)

	if update.Message.IsCommand() {
		h.handleCommand(bot, update.Message)
	}
}

func (h *Handler) handleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	lang := message.From.LanguageCode

	switch message.Command() {
	case "start":
		msgText := locales.Get(lang, "start_message")
		msg := tgbotapi.NewMessage(message.Chat.ID, msgText)
		bot.Send(msg)

	case "play":
		query := strings.TrimSpace(message.CommandArguments())
		if query == "" {
			msgText := locales.Get(lang, "play_usage")
			msg := tgbotapi.NewMessage(message.Chat.ID, msgText)
			bot.Send(msg)
			return
		}
		h.Player.Play(message.Chat.ID, query)

	default:
		// Handle other commands later
	}
}