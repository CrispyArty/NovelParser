package mailer

import (
	"github.com/crispyarty/novelparser/internal/mailer/common"
	"github.com/crispyarty/novelparser/internal/mailer/gmailapi"
)

func Validate() {
	if err := gmailapi.ValidateToken(); err != nil {
		gmailapi.CreateNewToken()
	}
}

func SendBook(bookpath string) {
	message := &common.Message{
		From: "crispykindle@gmai.com",
		To:   "crispykindle_4228f4@kindle.com",
		// To:       "artemti0@gmail.com",
		Subject:  "Hello from NovelParserGo!",
		Text:     "This is a book email sent from a NovelParserGo program.",
		Bookpath: bookpath,
	}

	gmailapi.SendBook(message)
}
