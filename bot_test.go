package icqbotapi

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"

	"icqbotapi/event"
)

const token = "001.1104030426.1757333006:757143498"

//func TestBot_Get(t *testing.T) {
//	bot := Bot{
//		token,
//		apiBaseURL,
//		http.DefaultClient,
//		time.Minute,
//	}
//
//	//r, err := bot.Get()
//	//if err != nil {
//	//	t.Fatal(err)
//	//}
//
//	//log.Printf("%+v", r)
//
//	x, cancel := bot.PollEvents()
//	<-x
//	cancel()
//}

func TestBot_PollEvents(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	bot := New(token, http.DefaultClient)
	bot.SetNewMessageHandler(func(e event.EventNewMessagePayload) {
		spew.Dump(e)
	})

	go bot.HandleEvents(ctx)

	time.Sleep(time.Second * 20)
	cancel()
}
