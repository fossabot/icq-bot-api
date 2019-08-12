package icqbotapi

import (
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

var bot = Bot{
	token,
	apiBaseURL,
	&http.Client{
		Timeout: time.Minute * 2,
	},
	time.Minute,
}

func TestBot_SendText(t *testing.T) {
	req := &SendTextRequest{
		ChatID: "p.radkov@corp.mail.ru",
		Text:   "kek",
	}

	resp, err := bot.SendText(req)
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("%+v", resp)
}

func TestBot_SendFile(t *testing.T) {
	req := &SendFileRequest{
		FileRequest: FileRequest{
			SendTextRequest: SendTextRequest{
				ChatID: "p.radkov@corp.mail.ru",
				Text:   "is's pepe",
			},
			Caption: "it's pepe",
		},
		FileID: "05j5L69UrfAdj8tZCGyi8H5d5160d61af",
	}

	go func() {
		for range time.Tick(time.Second * 10) {
			bot.SendChatActions(ChatActionsRequest{"p.radkov@corp.mail.ru", []ChatAction{}})
		}
	}()

	resp, err := bot.SendFile(req)
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("%#v", resp)

	x, cancel := bot.PollEvents()
	<-x
	cancel()
}

func TestBot_SendNewFile(t *testing.T) {
	f, err := os.Open("/Users/p.radkov/pepe.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	req := &SendNewFileRequest{
		FileRequest: FileRequest{
			SendTextRequest: SendTextRequest{
				ChatID: "p.radkov@corp.mail.ru",
				Text:   "is's pepe",
			},
			Caption: "it's pepe",
		},
		File:     f,
		Filename: "pepe.jpg",
	}

	resp, err := bot.SendNewFile(req)
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("%#v", resp)
}

func TestBot_EditMessage(t *testing.T) {
	req := &EditMessageRequest{
		ChatID:    "p.radkov@corp.mail.ru",
		MessageID: "6724288965706252425",
		Text:      "keklol",
	}

	resp, err := bot.EditMessage(req)
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("%#v", resp)
}

func TestBot_DeleteMessage(t *testing.T) {
	req := &DeleteMessageRequest{
		ChatID:    "p.radkov@corp.mail.ru",
		MessageID: "6724275801631490259",
	}

	resp, err := bot.DeleteMessage(req)
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("%#v", resp)
}
