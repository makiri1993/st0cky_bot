package main

import (
	"bing-news-api/api"
	. "bing-news-api/setup"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {

	InitDb()
	defer CloseDb()
	InitBot()

	api.SendNewsToUser()
}
