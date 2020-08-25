package setup

import (
	"bing-news-api/models"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jinzhu/gorm"
	"log"
)

const (
	Endpoint = "https://search-news.cognitiveservices.azure.com/bing/v7.0/news/search"
	Token    = "17dfa447c0734e158d858e6f45526ff3"
	botToken = "1245803898:AAFZtTRZYdCk5rdCo27hHc9Wrlbmr5E46DY"
	db       = "host=127.0.0.1 port=5432 user=postgres dbname=postgres password=password sslmode=disable"
)

func InitBot() (*tgbotapi.BotAPI, tgbotapi.UpdateConfig) {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	return bot, u
}

func InitDb() *gorm.DB {
	db, err := gorm.Open("postgres", db)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database Error:%s", err))
	}

	// Migrate the schema
	db.AutoMigrate(&models.News{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Keyword{})
	return db
}

func CloseDb(db *gorm.DB) {
	db.Close()
}
