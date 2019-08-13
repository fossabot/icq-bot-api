package icqbotapi

import (
	"context"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

var bot = New(
	token,
	&http.Client{
		Timeout: time.Minute * 2,
	})

func TestBot_SendText(t *testing.T) {
	req := &SendTextRequest{
		ChatID: "p.radkov@corp.mail.ru",
		Text:   "kek",
	}

	resp, err := bot.SendText(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("%+v", resp)
}

func TestBot_SendFile(t *testing.T) {
	req := &SendFileRequest{
		fileRequest: fileRequest{
			SendTextRequest: SendTextRequest{
				ChatID: "p.radkov@corp.mail.ru",
				Text:   "is's pepe",
			},
			Caption: "it's pepe",
		},
		FileID: "05j5L69UrfAdj8tZCGyi8H5d5160d61af",
	}

	resp, err := bot.SendFile(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("%#v", resp)
}

func TestBot_SendNewFile(t *testing.T) {
	f, err := os.Open("/Users/p.radkov/pepe.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	req := &SendNewFileRequest{
		fileRequest: fileRequest{
			SendTextRequest: SendTextRequest{
				ChatID: "p.radkov@corp.mail.ru",
				Text:   "is's pepe",
			},
			Caption: "it's pepe",
		},
		File:     f,
		Filename: "pepe.jpg",
	}

	resp, err := bot.SendNewFile(context.Background(), req)
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

	resp, err := bot.EditMessage(context.Background(), req)
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

	resp, err := bot.DeleteMessage(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("%#v", resp)
}
