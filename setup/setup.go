package setup

import (
	"bing-news-api/models"
	"fmt"
	"github.com/jinzhu/gorm"
	. "gopkg.in/tucnak/telebot.v2"
	"log"
	"os"
	"time"
)

const (
	Endpoint = "https://search-news.cognitiveservices.azure.com/bing/v7.0/news/search"
	Token    = "17dfa447c0734e158d858e6f45526ff3"
	botToken = "1245803898:AAFZtTRZYdCk5rdCo27hHc9Wrlbmr5E46DY"
)

var DbConn *gorm.DB
var TelegramBot *Bot

func InitBot() {
	var err error

	TelegramBot, err = NewBot(Settings{
		// You can also set custom API URL.
		// If field is empty it equals to "https://api.telegram.org".
		Token:  botToken,
		Poller: &LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", TelegramBot.Me.Username)

}

func InitDb() {
	var err error
	DbConn, err = gorm.Open("postgres", os.Getenv("db"))
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database Error:%s", err))
	}

	//DbConn.LogMode(true)
	// Migrate the schema
	DbConn.AutoMigrate(&models.News{})
	DbConn.AutoMigrate(&models.User{})
	DbConn.AutoMigrate(&models.Keyword{})
}

func CloseDb() {
	DbConn.Close()
}
