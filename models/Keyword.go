package models

import "github.com/jinzhu/gorm"

type Keyword struct {
	gorm.Model
	UserID     int64
	Searchterm string
}
