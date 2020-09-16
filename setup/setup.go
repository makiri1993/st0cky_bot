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

var DbConn *gorm.DB
var Bot *tgbotapi.BotAPI
var BotUpdateConfig tgbotapi.UpdateConfig

func InitBot() {
	var err error
	Bot, err = tgbotapi.NewBotAPI(botToken)

	if err != nil {
		log.Panic(err)
	}

	Bot.Debug = true

	log.Printf("Authorized on account %s", Bot.Self.UserName)

	BotUpdateConfig = tgbotapi.NewUpdate(0)
	BotUpdateConfig.Timeout = 60
}

func InitDb() {
	var err error
	DbConn, err = gorm.Open("postgres", db)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database Error:%s", err))
	}

	// Migrate the schema
	DbConn.AutoMigrate(&models.News{})
	DbConn.AutoMigrate(&models.User{})
	DbConn.AutoMigrate(&models.Keyword{})
}

func CloseDb() {
	DbConn.Close()
}
