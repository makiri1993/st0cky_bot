package api

import (
	"bing-news-api/db"
	. "bing-news-api/setup"
	. "gopkg.in/tucnak/telebot.v2"
	"log"
	"strconv"
)

var (
	helloCommand             = Command{Text: "/hello", Description: "Say hello to St0cky. Be nice man."}
	addKeywordCommand        = Command{Text: "/add_keyword", Description: "Add a new keyword for your search queries. Example: /add_keyword <keyword here>."}
	removeKeywordCommand     = Command{Text: "/remove_keyword", Description: "Remove an existing keyword from your search queries."}
	getNewsCommand           = Command{Text: "/get_news", Description: "Get news by providing a keyword. There are no checks so you will just get the newest information."}
	toggleLiveUpdatesCommand = Command{Text: "/toggle_live_updates", Description: "Toggle the mechanism for the automatic sending of new news."}
	helpCommand              = Command{Text: "/help", Description: "Just some help."}
)

func RegisterRoutes() {
	err := TelegramBot.SetCommands([]Command{helloCommand, addKeywordCommand, removeKeywordCommand, getNewsCommand, toggleLiveUpdatesCommand, helpCommand})

	if err != nil {
		panic(err)
	}

	TelegramBot.Handle(helloCommand.Text, helloHandler)
	TelegramBot.Handle(addKeywordCommand.Text, addKeywordHandler)
	TelegramBot.Handle(removeKeywordCommand.Text, removeKeywordHandler)
	TelegramBot.Handle(getNewsCommand.Text, getNewsHandler)
	//TelegramBot.Handle(toggleLiveUpdatesCommand.Text, toggleLiveUpdatesHandler)
	TelegramBot.Handle(helpCommand.Text, helpHandler)

}

func checkIfUserExists(m *Message) {
	isCreated, createdMessage := db.FindOrCreateUser(m.Sender.Username, m.Chat.ID)
	if isCreated == true {
		safeTelegramSend(m.Sender, createdMessage)
	}
}

func helloHandler(m *Message) {
	checkIfUserExists(m)

	safeTelegramSend(m.Sender, "Hey man! I'm looking forward to help you.")
}

func addKeywordHandler(m *Message) {
	checkIfUserExists(m)

	if m.Payload == "" {
		safeTelegramSend(m.Sender, "Can you repeat the keyword? I didn't get it.")
		return
	}

	message := db.FindOrCreateKeyword(m.Payload, m.Chat.ID)
	safeTelegramSend(m.Sender, message)
}

func removeKeywordHandler(m *Message) {
	checkIfUserExists(m)

	keywords := db.GetAllKeywordsPerUser(m.Chat.ID)

	var (
		selector   = &ReplyMarkup{}
		inlineRows []Btn
	)
	for _, keyword := range keywords {
		stringifiedID := strconv.Itoa(int(keyword.ID))
		btn := selector.Data(keyword.Searchterm, stringifiedID, stringifiedID)
		inlineRows = append(inlineRows, selector.Row(btn)...)

		TelegramBot.Handle(&btn, func(c *Callback) {

			keywordID, err := strconv.ParseUint(c.Data, 10, 64)

			if err != nil {
				log.Panic(err)
			}

			db.DeleteKeyword(keywordID)

			err = TelegramBot.Respond(c, &CallbackResponse{
				Text: "Keyword deleted!",
			})

			if err != nil {
				log.Panicf("Problem with callback for %v", m.Chat.ID)
			}
		})
	}

	selector.Inline(inlineRows)
	safeTelegramSend(m.Sender, "Please select the keyword you want to remove.", selector)

}

func helpHandler(m *Message) {
	checkIfUserExists(m)

	safeTelegramSend(m.Sender, "Here will be the help")
}

func getNewsHandler(m *Message) {
	checkIfUserExists(m)

	if m.Payload == "" {
		safeTelegramSend(m.Sender, "Can you repeat the search term? I didn't get it.")

		return
	}

	newsResult := GetBingNews(m.Payload)
	news := db.FindOrCreateNews(newsResult.ToNewsStructs(m.Chat.ID, m.Payload), m.Chat.ID, m.Payload)

	safeTelegramSend(m.Sender, news.ToString())

	db.UpdateNews(news)
}

func safeTelegramSend(to *User, what string, options ...interface{}) {
	if len(options) == 0 {
		message, err := TelegramBot.Send(to, what)
		if err != nil {
			log.Panicf("Error with message: %v \n\n %s", message, err)
		}
		return
	}

	message, err := TelegramBot.Send(to, what, options...)
	if err != nil {
		log.Panicf("Error with message: %v \n\n %s", message, err)
	}
}

//func toggleLiveUpdatesHandler(m *Message) {
//
//}
