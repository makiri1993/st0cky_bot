package db

import (
	. "bing-news-api/models"
	. "bing-news-api/setup"
	"github.com/jinzhu/gorm"
)

func FindOrCreateKeyword(searchTerm string, userID int64) string {
	var keyword Keyword

	if err := DbConn.Where("user_id = ? and searchterm = ?", userID, searchTerm).First(&keyword).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			DbConn.Create(&Keyword{
				Searchterm: searchTerm,
				UserID:     userID,
			})

			return "Created keyword: '" + searchTerm + "'."
		}
	}
	return "No worries! I'm already looking for this keyword."
}

func GetAllKeywordsPerUser(userID int64) Keywords {
	var entries Keywords
	DbConn.Model(&Keyword{}).Where("user_id = ?", userID).Find(&entries)
	return entries
}

func DeleteKeyword(keywordID uint64) {
	DbConn.Delete(&Keyword{}, keywordID)
}
