package db

import (
	. "bing-news-api/models"
	. "bing-news-api/setup"
	"github.com/jinzhu/gorm"
)

func FindOrCreateUser(userName string, telegramID int64) string {
	var user User

	if err := DbConn.Where("telegram_chat_id = ?", telegramID).First(&user).Error; err != nil {
		// error handling...
		if gorm.IsRecordNotFoundError(err) {
			DbConn.Create(&User{
				Name:           userName,
				TelegramChatID: telegramID,
			})

			return "Created user: " + userName + ". Please repeat your command."
		}
	}
	return ""
}
