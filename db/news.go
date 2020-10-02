package db

import (
	. "bing-news-api/models"
	. "bing-news-api/setup"
	"github.com/jinzhu/gorm"
	"log"
)

func FindOrCreateNews(news []News, userID int64, searchTerm string) NewsArray {
	var entries NewsArray
	for _, entry := range news {
		var newsEntry News
		if err := DbConn.Where(&News{Title: entry.Title, Url: entry.Url,
			SearchTerm: searchTerm, UserID: userID}).First(&newsEntry).Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				DbConn.Create(&entry)
				entries = append(entries, entry)
			}
		} else {
			DbConn.Model(&newsEntry).Where(entry).First(&newsEntry)
			if !newsEntry.Sent {
				entries = append(entries, newsEntry)
			}
		}
	}
	return entries
}

func UpdateNews(news []News) {
	for _, entry := range news {
		DbConn.Model(&entry).Update("sent", true)
	}
}

func FindAllNewsPerUserID(userID int64) NewsArray {
	var entries NewsArray
	DbConn.Where("user_id = ? and sent = ?", userID, false).Find(&entries)

	for _, news := range entries {
		log.Println(news.Sent)
	}
	return entries
}
