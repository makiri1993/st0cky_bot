package main

import (
	"bing-news-api/api"
	. "bing-news-api/setup"
	"bing-news-api/tasks"
	"log"
	"time"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gopkg.in/robfig/cron.v2"
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
	tasks.UpdateNewsForEveryUser()
	tasks.SendNewNewsToEveryUser()
	addCronFunc(c, "@every 0h60m0s", tasks.UpdateNewsForEveryUser)
	addCronFunc(c, "@every 0h70m0s", tasks.SendNewNewsToEveryUser)
	c.Start()

	// Added time to see output
	time.Sleep(10 * time.Second)

	//c.Stop() // Stop the scheduler (does not stop any jobs already running).
}

func addCronFunc(c *cron.Cron, interval string, handler func()) {
	_, err := c.AddFunc(interval, handler)

	if err != nil {
		log.Panic("Init cron function failed.", err)
	}
}
