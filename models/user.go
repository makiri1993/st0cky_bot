package models

import (
	"fmt"
	"time"
)

type User struct {
	ID               int64 `gorm:"primaryKey"`
	Name             string
	AutomaticSending bool
	Keywords         []Keyword
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time `sql:"index"`
}

func (user User) ToString() string {
	return fmt.Sprintf("Name: %v \nWith keywords: %v", user.Name, user.Keywords)
}
