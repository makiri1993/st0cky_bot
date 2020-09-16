package api

import (
	"bing-news-api/db"
	"bing-news-api/models"
	. "bing-news-api/setup"
	. "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

const (
	KeyboardNews        = "Get news"
	KeyboardAddKeyboard = "Add search keyword"
	KeyboardUser        = "Get user data"
	CommandStart        = "start"
)

var numericKeyboard = NewReplyKeyboard(
	NewKeyboardButtonRow(
		NewKeyboardButton(KeyboardNews),
		NewKeyboardButton(KeyboardAddKeyboard),
		NewKeyboardButton(KeyboardUser),
	),
)

func SendNewsToUser() {

	updates, _ := Bot.GetUpdatesChan(BotUpdateConfig)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		successMessage := db.FindOrCreateUser(update.Message.From.UserName, update.Message.Chat.ID)
		msg := NewMessage(update.Message.Chat.ID,
			successMessage)
		Bot.Send(msg)

		if handleCommands(update) {
			return
		}

		handleKeyboardTasks(update, msg)
	}
}

func handleKeyboardTasks(update Update, msg MessageConfig) {
	switch update.Message.Text {
	case KeyboardNews:
		handleTelegramUpdate(update)
	case KeyboardAddKeyboard:
		handleTelegramUpdate(update)
	default:
		msg.Text = "I don't know that task"
	}
	Bot.Send(msg)
}

func handleCommands(update Update) bool {
	if update.Message.IsCommand() {
		msg := NewMessage(update.Message.Chat.ID, "")

		switch update.Message.Command() {
		case CommandStart:
			msg.ReplyMarkup = numericKeyboard
			msg.Text = "What can I do for you?"
		default:
			msg.Text = "I don't know that command"
		}
		Bot.Send(msg)
		return true
	}
	return false
}

func handleTelegramUpdate(update Update) {

	newsResult := GetBingNews(update.Message.Text)
	news := models.NewsToString(newsResult)
	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

	msg := NewMessage(update.Message.Chat.ID, news)

	Bot.Send(msg)
}
