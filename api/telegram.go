package api

import (
	"bing-news-api/db"
	"bing-news-api/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jinzhu/gorm"
	"log"
)

const (
	KeyboardNews = "Get news"
	KeyboardUser = "Get user data"
	CommandStart = "start"
)

var numericKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(KeyboardNews),
		tgbotapi.NewKeyboardButton(KeyboardUser),
	),
)

func SendNewsToUser(bot *tgbotapi.BotAPI, u tgbotapi.UpdateConfig, conn *gorm.DB) {

	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		successMessage := db.FindOrCreateUser(conn, update.Message.From.UserName, update.Message.Chat.ID)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID,
			successMessage)
		bot.Send(msg)

		if handleCommands(bot, update) {
			return
		}

		handleKeyboardTasks(bot, update, msg)
	}
}

func handleKeyboardTasks(bot *tgbotapi.BotAPI, update tgbotapi.Update, msg tgbotapi.MessageConfig) {
	switch update.Message.Text {
	case KeyboardNews:
		handleTelegramUpdate(update, bot)

	default:
		msg.Text = "I don't know that task"
	}
	bot.Send(msg)
}

func handleCommands(bot *tgbotapi.BotAPI, update tgbotapi.Update) bool {
	if update.Message.IsCommand() {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		switch update.Message.Command() {
		case CommandStart:
			msg.ReplyMarkup = numericKeyboard
			msg.Text = "What can I do for you?"
		default:
			msg.Text = "I don't know that command"
		}
		bot.Send(msg)
		return true
	}
	return false
}

func handleTelegramUpdate(update tgbotapi.Update, bot *tgbotapi.BotAPI) {

	newsResult := GetBingNews(update.Message.Text)
	news := models.NewsToString(newsResult)
	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, news)

	bot.Send(msg)
}
