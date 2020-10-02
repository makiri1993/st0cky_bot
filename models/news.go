package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"strings"
)

// This struct formats the answer provided by the Bing News Search API.
type BingNews struct {
	ReadLink     string `json:"readLink"`
	QueryContext struct {
		OriginalQuery string `json:"originalQuery"`
		AdultIntent   bool   `json:"adultIntent"`
	} `json:"queryContext"`
	TotalEstimatedMatches int `json:"totalEstimatedMatches"`
	Sort                  []struct {
		Name       string `json:"name"`
		ID         string `json:"id"`
		IsSelected bool   `json:"isSelected"`
		URL        string `json:"url"`
	} `json:"sort"`
	Value []struct {
		Name  string `json:"name"`
		URL   string `json:"url"`
		Image struct {
			Thumbnail struct {
				ContentURL string `json:"thumbnail"`
				Width      int    `json:"width"`
				Height     int    `json:"height"`
			} `json:"thumbnail"`
		} `json:"image"`
		Description string `json:"description"`
		Provider    []struct {
			Type string `json:"_type"`
			Name string `json:"name"`
		} `json:"provider"`
		DatePublished string `json:"datePublished"`
	} `json:"value"`
}

type News struct {
	gorm.Model
	Title      string
	Url        string
	Source     string
	Date       string
	Sent       bool
	SearchTerm string
	UserID     int64
}

type NewsArray []News

func prettifyDate(date string) string {
	return strings.Replace(date[0:16], "T", " ", 1)
}

func (bingNews BingNews) ToString() string {
	var news []string
	for _, result := range bingNews.Value {
		news = append(news, result.Name, "\n", result.URL, "\n", prettifyDate(result.DatePublished), "\n\n")
	}
	return strings.Join(news, "")
}

func (bingNews BingNews) ToNewsStructs(userID int64, searchTerm string) NewsArray {
	var news []News
	for _, result := range bingNews.Value {
		source := result.Provider[0]
		news = append(news, News{Title: result.Name, Url: result.URL, Source: source.Name,
			Date: result.DatePublished, Sent: false, UserID: userID, SearchTerm: searchTerm})
	}

	return news
}

func (news News) ToString() string {
	return fmt.Sprintf("%v\n%v\n%v\n\n", news.Title, news.Url, prettifyDate(news.Date))
}

func (news NewsArray) ToString() string {
	var newsString []string
	for _, result := range news {
		newsString = append(newsString, result.ToString())
	}
	return strings.Join(newsString, "")
}
