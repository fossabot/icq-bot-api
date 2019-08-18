package icqbotapi

import (
	"context"
	"log"
	"net/http"
	"testing"
	"time"

	"icqbotapi/event"
)

func TestBot_GetSelf(t *testing.T) {
	const token = "001.1104030426.1757333006:757143498"
	bot := New(token, http.DefaultClient, APITypeICQ)

	r, err := bot.GetSelf(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if !r.Ok {
		t.Fatal("unexpected response status")
	}
}

func ExampleNew() {
	const token = "001.1104030426.1757333006:757143498"
	bot := New(token, http.DefaultClient, APITypeICQ)
	_ = bot
}

func ExampleBot_GetSelf() {
	const token = "001.1104030426.1757333006:757143498"
	bot := New(token, http.DefaultClient, APITypeICQ)

	data, _ := bot.GetSelf(context.Background())

	log.Printf("%#v", data)
}

func ExampleBot_PollEvents() {
	const token = "001.1104030426.1757333006:757143498"
	bot := New(token, http.DefaultClient, APITypeICQ)

	for event := range bot.PollEvents(context.Background()) {
		log.Printf("%#v", event)
	}
}

func ExampleBot_HandleEvents() {
	const token = "001.1104030426.1757333006:757143498"
	bot := New(token, http.DefaultClient, APITypeICQ)

	bot.SetNewMessageHandler(func(e event.NewMessagePayload) {
		log.Printf("%#v", e)
	})

	bot.SetErrorHandler(func(err error) {
		log.Print(err)
	})

	ctx, cancel := context.WithCancel(context.Background())
	bot.HandleEvents(ctx)
	time.Sleep(time.Second * 3)
	cancel()
}
