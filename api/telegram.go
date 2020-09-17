package api

import (
	"bing-news-api/db"
	. "bing-news-api/setup"
	. "gopkg.in/tucnak/telebot.v2"
	"strings"
)

var (
	HelloCommand      = Command{Text: "/hello", Description: "Say hello to St0cky. Be nice man."}
	AddKeywordCommand = Command{Text: "/add_keyword", Description: "Add a new keyword for your search queries. Example: /add_keyword <keyword here>."}
	GetNewsCommand    = Command{Text: "/get_news", Description: "Get news by providing a keyword. There are no checks so you will just get the newest information."}
	HelpCommand       = Command{Text: "/help", Description: "Just some help."}
)

func RegisterRoutes() {
	TelegramBot.SetCommands([]Command{HelloCommand, AddKeywordCommand, GetNewsCommand, HelpCommand})

	TelegramBot.Handle(HelloCommand.Text, helloHandler)
	TelegramBot.Handle(AddKeywordCommand.Text, addKeywordHandler)
	TelegramBot.Handle(GetNewsCommand.Text, getNewsHandler)
	TelegramBot.Handle(HelpCommand.Text, helpHandler)

	TelegramBot.Start()
}

func checkIfUserExists(m *Message) {
	isCreated := db.FindOrCreateUser(m.Sender.Username, m.Chat.ID)
	if isCreated != "" {
		TelegramBot.Send(m.Sender, isCreated)

	}
}

func helloHandler(m *Message) {
	checkIfUserExists(m)
	TelegramBot.Send(m.Sender, "Hey man! I'm looking forward to help you.")
}

func addKeywordHandler(m *Message) {
	checkIfUserExists(m)

	if m.Payload == "" {
		TelegramBot.Send(m.Sender, "Can you repeat the keyword? I didn't get it.")
		return
	}

	message := db.FindOrCreateKeyword(m.Payload, m.Chat.ID)
	TelegramBot.Send(m.Sender, message)
}

func helpHandler(m *Message) {
	checkIfUserExists(m)

	TelegramBot.Send(m.Sender, "Here will be the help")
}

func getNewsHandler(m *Message) {
	checkIfUserExists(m)

	if m.Payload == "" {
		TelegramBot.Send(m.Sender, "Can you repeat the search term? I didn't get it.")
		return
	}

	newsResult := GetBingNews(m.Payload)
	news := db.FindOrCreateNews(newsResult.ToNewsStructs())

	var newsString []string
	for _, result := range news {
		newsString = append(newsString, result.ToString())
	}

	TelegramBot.Send(m.Sender, strings.Join(newsString, ""))

	db.UpdateNews(news)
}
