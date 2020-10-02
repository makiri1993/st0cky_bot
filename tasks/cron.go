package tasks

import (
	. "bing-news-api/api"
	"bing-news-api/db"
	. "bing-news-api/setup"
	"gopkg.in/tucnak/telebot.v2"
	"log"
)

func UpdateNewsForEveryUser() {
	users := db.GetAllUsers()

	for _, user := range users {
		for _, keyword := range user.Keywords {
			newsResult := GetBingNews(keyword.Searchterm)
			db.FindOrCreateNews(newsResult.ToNewsStructs(user.ID, keyword.Searchterm),
				user.ID, keyword.Searchterm)
			log.Printf("Fetched news for keyword: %s\n", keyword.Searchterm)
		}
		log.Printf("Fetched news for the keywords from user: %s\n", user.Name)
	}
}

func SendNewNewsToEveryUser() {
	users := db.GetAllUsers()
	for _, user := range users {
		news := db.FindAllNewsPerUserID(user.ID)
		TelegramBot.Send(&telebot.User{ID: int(user.ID)}, news.ToString())
		db.UpdateNews(news)
		log.Printf("Sent news for user: %s\n", user.Name)

	}
}
