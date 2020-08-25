package models

import (
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
	Url    string
	Source string
	Date   string
	Sent   bool
}

func NewsToString(ans *BingNews) string {
	var news []string
	for _, result := range ans.Value {
		news = append(news, result.Name, "\n", result.URL, "\n\n")
	}
	return strings.Join(news, "")
}
