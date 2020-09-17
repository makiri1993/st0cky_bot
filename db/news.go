package db

import (
	. "bing-news-api/models"
	. "bing-news-api/setup"
	"github.com/jinzhu/gorm"
)

func FindOrCreateNews(news []News) []News {
	var entries []News
	for _, entry := range news {
		var newsEntry News
		if err := DbConn.Where(entry).First(&newsEntry).Error; err != nil {
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
