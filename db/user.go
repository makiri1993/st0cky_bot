package db

import (
	"bing-news-api/models"
	"github.com/jinzhu/gorm"
)

func FindOrCreateUser(db *gorm.DB, userName string, telegramID int64, message string) string {
	var user models.User
	keywords := []models.Keyword{{Searchterm: message}}

	db.Where(&user).Assign(models.User{
		Name:           userName,
		TelegramChatID: telegramID,
		Keywords:       keywords}).FirstOrCreate(&user)

	return "Created user: " + user.Name
}
