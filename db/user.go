package db

import (
	"bing-news-api/models"
	"github.com/jinzhu/gorm"
)

func FindOrCreateUser(db *gorm.DB, userName string, telegramID int64) string {
	var user models.User

	db.Where(&user).Assign(models.User{
		Name:           userName,
		TelegramChatID: telegramID,
	}).FirstOrCreate(&user)

	return "Created user: " + user.Name
}
