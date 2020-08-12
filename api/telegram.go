package api

import (
	"bing-news-api/db"
	"bing-news-api/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jinzhu/gorm"
	"log"
)

func SendNewsToUser(bot *tgbotapi.BotAPI, u tgbotapi.UpdateConfig, conn *gorm.DB) {

	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		successMessage := db.FindOrCreateUser(conn, update.Message.From.UserName, update.Message.Chat.ID, update.Message.Text)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID,
			successMessage)
		bot.Send(msg)

		handleTelegramUpdate(update, bot)

	}
}

func handleTelegramUpdate(update tgbotapi.Update, bot *tgbotapi.BotAPI) {

	newsResult := GetBingNews(update.Message.Text)
	news := models.NewsToString(newsResult)
	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, news)

	bot.Send(msg)
}
