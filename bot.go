package icqbotapi

import (
	"net/http"
	"net/url"
	"time"
)

const apiBaseURL = "https://api.icq.net/bot/v1"
const tokenQueryParam = "token"

type Bot struct {
	token        string
	apiBaseURL   string
	client       *http.Client
	pollDuration time.Duration
}

func New(token string, client *http.Client) *Bot {
	if client == nil {
		client = http.DefaultClient
	}

	return &Bot{
		token,
		apiBaseURL,
		client,
		time.Minute,
	}
}

func (b *Bot) setToken(q url.Values) {
	q.Add(tokenQueryParam, b.token)
}
