package main

import (
	"bing-news-api/api"
	"bing-news-api/setup"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {

	db := setup.InitDb()
	defer setup.CloseDb(db)

	bot, u := setup.InitBot()
	api.SendNewsToUser(bot, u, db)
}
