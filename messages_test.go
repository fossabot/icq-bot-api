package icqbotapi

import (
	"context"
	"log"
	"net/http"
	"os"
)

func ExampleBot_SendText() {
	const token = "001.1104030426.1757333006:757143498"
	bot := New(token, http.DefaultClient, APITypeICQ)

	req := &SendTextRequest{
		ChatID: "chat1",
		Text:   "kek",
	}

	resp, _ := bot.SendText(context.Background(), req)

	log.Printf("%+v", resp)
}

func ExampleBot_SendFile() {
	const token = "001.1104030426.1757333006:757143498"
	bot := New(token, http.DefaultClient, APITypeICQ)

	req := &SendFileRequest{
		fileRequest: fileRequest{
			SendTextRequest: SendTextRequest{
				ChatID: "chat1",
				Text:   "is's pepe",
			},
			Caption: "it's pepe",
		},
		FileID: "05j5L69UrfAdj8tZCGyi8H5d5160d61af",
	}

	resp, _ := bot.SendFile(context.Background(), req)

	log.Printf("%#v", resp)
}

func ExampleBot_SendNewFile() {
	const token = "001.1104030426.1757333006:757143498"
	bot := New(token, http.DefaultClient, APITypeICQ)

	f, err := os.Open("./pepe.jpg")
	if err != nil {
		panic(err)
	}

	defer f.Close()

	req := &SendNewFileRequest{
		fileRequest: fileRequest{
			SendTextRequest: SendTextRequest{
				ChatID: "chat1",
				Text:   "is's pepe",
			},
			Caption: "it's pepe",
		},
		File:     f,
		Filename: "pepe.jpg",
	}

	resp, _ := bot.SendNewFile(context.Background(), req)

	log.Printf("%#v", resp)
}

func ExampleBot_EditMessage() {
	const token = "001.1104030426.1757333006:757143498"
	bot := New(token, http.DefaultClient, APITypeICQ)

	req := &EditMessageRequest{
		ChatID:    "chat1",
		MessageID: "6724288965706252425",
		Text:      "kek",
	}

	resp, _ := bot.EditMessage(context.Background(), req)

	log.Printf("%#v", resp)
}

func ExampleBot_DeleteMessage() {
	const token = "001.1104030426.1757333006:757143498"
	bot := New(token, http.DefaultClient, APITypeICQ)

	req := &DeleteMessageRequest{
		ChatID:    "chat1",
		MessageID: "6724275801631490259",
	}

	resp, _ := bot.DeleteMessage(context.Background(), req)

	log.Printf("%#v", resp)
}
