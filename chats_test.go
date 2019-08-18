package icqbotapi

import (
	"context"
	"log"
	"net/http"
	"testing"
	"time"
)

func ExampleBot_GetChatAdmins() {
	const token = "001.1104030426.1757333006:757143498"
	bot := New(token, http.DefaultClient, APITypeICQ)

	data, _ := bot.GetChatAdmins(context.Background(), "chat1")

	log.Printf("%#v", data)
}

func TestName(t *testing.T) {
	const token = "001.1104030426.1757333006:757143498"
	bot := New(token, http.DefaultClient, APITypeICQ)
	events := make(chan ChatActionsRequest)

	bot.SendChatActions(context.Background(), events)

	chats := []ChatID{
		"chat1",
		"chat2",
		"chat3",
	}

	for i := 0; i < 50; i++ {
		events <- ChatActionsRequest{
			ChatID: chats[i%len(chats)],
			Actions: []ChatAction{
				ChatActionTyping,
			},
		}

		time.Sleep(time.Millisecond * 500)
	}

	close(events)
}

func ExampleBot_SendChatActions() {
	const token = "001.1104030426.1757333006:757143498"
	bot := New(token, http.DefaultClient, APITypeICQ)
	events := make(chan ChatActionsRequest)

	bot.SendChatActions(context.Background(), events)

	chats := []ChatID{
		"chat1",
		"chat2",
		"chat3",
	}

	for i := 0; i < 50; i++ {
		events <- ChatActionsRequest{
			ChatID: chats[i%len(chats)],
			Actions: []ChatAction{
				ChatActionTyping,
			},
		}

		time.Sleep(time.Millisecond * 100)
	}

	close(events)
}

func ExampleBot_GetChatInfo() {
	const token = "001.1104030426.1757333006:757143498"
	bot := New(token, http.DefaultClient, APITypeICQ)
	chatInfo, _ := bot.GetChatInfo(context.Background(), "chat1")

	log.Printf("%#v", chatInfo)
}
