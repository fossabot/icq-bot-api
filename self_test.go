package icqbotapi

import (
	"net/http"
	"testing"
	"time"
)

const token = "001.1104030426.1757333006:757143498"

func TestBot_Get(t *testing.T) {
	bot := Bot{
		token,
		apiBaseURL,
		http.DefaultClient,
		time.Minute,
	}

	//r, err := bot.Get()
	//if err != nil {
	//	t.Fatal(err)
	//}

	//log.Printf("%+v", r)

	x, cancel := bot.PollEvents()
	<-x
	cancel()
}
