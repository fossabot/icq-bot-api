package icqbotapi

import (
	"log"
	"net/http"
	"testing"
)

const token = "001.1104030426.1757333006:757143498"

func TestBot_Get(t *testing.T) {
	bot := Bot{
		token,
		apiBaseURL,
		http.DefaultClient,
	}

	r, err := bot.Get()
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("%+v", r)
}