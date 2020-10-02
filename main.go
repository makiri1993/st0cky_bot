package main

import (
	"bing-news-api/api"
	. "bing-news-api/setup"
	"bing-news-api/tasks"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gopkg.in/robfig/cron.v2"
	"time"
)

func main() {

	InitDb()
	defer CloseDb()
	InitBot()

	api.RegisterRoutes()
	addCronJob()

	TelegramBot.Start()
}

func addCronJob() {
	c := cron.New()
	c.AddFunc("@every 0h30m0s", tasks.UpdateNewsForEveryUser)
	c.AddFunc("@every 0h20m0s", tasks.SendNewNewsToEveryUser)
	c.Start()

	// Added time to see output
	time.Sleep(10 * time.Second)

	//c.Stop() // Stop the scheduler (does not stop any jobs already running).
}
