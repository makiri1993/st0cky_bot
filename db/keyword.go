package db

import (
	. "bing-news-api/models"
	. "bing-news-api/setup"
)

func FindOrCreateKeyword(searchTerm string, userID int64) string {
	var keyword Keyword

	DbConn.Where(&keyword).Assign(Keyword{
		Searchterm: searchTerm,
		UserID:     userID,
	}).FirstOrCreate(&keyword)

	return "Created keyword: " + keyword.Searchterm
}
