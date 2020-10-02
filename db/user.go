package db

import (
	. "bing-news-api/models"
	. "bing-news-api/setup"
	"github.com/jinzhu/gorm"
	"log"
)

func FindOrCreateUser(userName string, telegramID int64) string {
	var user User

	if err := DbConn.Where("id = ?", telegramID).First(&user).Error; err != nil {
		// error handling...
		if gorm.IsRecordNotFoundError(err) {
			DbConn.Create(&User{
				Name:             userName,
				ID:               telegramID,
				AutomaticSending: true,
			})

			return "Created user: " + userName + ". Please repeat your command."
		}
	}
	return ""
}

func GetAllUsers() []User {
	var users []User
	// Get all records
	result := DbConn.Preload("Keywords").Find(&users)
	// SELECT * FROM users;
	if result.Error != nil {
		log.Panic(result.Error)
	}
	return users
}

func UpdateAutomaticNewsSendingOfUser() {}
