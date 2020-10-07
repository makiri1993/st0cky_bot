package models

import (
	"strings"

	"github.com/jinzhu/gorm"
)

type Keyword struct {
	gorm.Model
	UserID     int64
	Searchterm string
}

type Keywords []Keyword

func (k Keywords) ToString() string {
	var keywordsString []string
	for _, keyword := range k {
		keywordsString = append(keywordsString, keyword.Searchterm)
	}
	return strings.Join(keywordsString, " ")
}
