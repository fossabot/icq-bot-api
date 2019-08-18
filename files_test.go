package icqbotapi

import (
	"context"
	"log"
	"net/http"
)

func ExampleBot_GetFileInfo() {
	const token = "001.1104030426.1757333006:757143498"
	bot := New(token, http.DefaultClient, APITypeICQ)

	data, _ := bot.GetFileInfo(context.Background(), "05j5Lk9Oka1VBpeMLS7Qv35d5189ac1af")

	log.Printf("%#v", data)
}
