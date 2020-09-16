package db

import (
	. "bing-news-api/models"
	. "bing-news-api/setup"
)

func FindOrCreateUser(userName string, telegramID int64) string {
	var user User

	DbConn.Where(&user).Assign(User{
		Name:           userName,
		TelegramChatID: telegramID,
	}).FirstOrCreate(&user)

	return "Created user: " + user.Name
}
