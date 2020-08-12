package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Name           string
	TelegramChatID int64
	Keywords       []Keyword `gorm:"foreignkey:UserID"`
}

type Keyword struct {
	gorm.Model
	UserID     int
	Searchterm string
}

func (user User) ToString() string {
	return fmt.Sprintf("Name: %v \nWith keywords: %v", user.Name, user.Keywords)
}
