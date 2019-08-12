package icqbotapi

import (
	"log"
	"net/http"
	"testing"
	"time"
)

func TestBot_GetChatAdmins(t *testing.T) {
	bot := Bot{
		token,
		apiBaseURL,
		http.DefaultClient,
	}

	chatID := ChatID("p.radkov@corp.mail.ru")
	data, err := bot.GetChatAdmins(chatID)
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("%#v", data)
}

func TestBot_SendChatActions(t *testing.T) {
	bot := Bot{
		token,
		apiBaseURL,
		http.DefaultClient,
	}

	chatID := ChatID("p.radkov@corp.mail.ru")
	for i := 0; i < 5; i++ {
		data, err := bot.SendChatActions(ChatActionsRequest{ChatID: chatID, Actions: []string{"typing"}})
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("%#v", data)
		time.Sleep(time.Second * 2)
	}

}
