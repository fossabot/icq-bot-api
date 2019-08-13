package icqbotapi

import (
	"context"
	"net/http"
	"time"

	"github.com/mailru/easyjson"

	"icqbotapi/event"
)

const apiBaseURL = "https://api.icq.net/bot/v1"
const tokenQueryParam = "token"

type newMessageHandlerFunc func(e event.EventNewMessagePayload)
type editMessageHandlerFunc func(e event.EventTypeEditedPayload)
type deleteMessageHandlerFunc func(e event.EventTypeDeletedPayload)
type pinMessageHandlerFunc func(e event.EventTypePinnedPayload)
type unpinMessageHandlerFunc func(e event.EventTypeUnpinnedPayload)
type newChatMembersHandlerFunc func(e event.EventTypeNewChatMembersPayload)
type leftChatMembersHandlerFunc func(e event.EventTypeLeftChatMembersPayload)
type changeChatInfoHandlerFunc func(e event.EventTypeLeftChatMembersPayload)

type botState int

const (
	botStatePolling = iota
	botStateHandling
	botStateStopped
)

type Bot struct {
	state botState

	token        string
	apiBaseURL   string
	client       *http.Client
	pollDuration time.Duration

	handlers struct {
		newMessageHandler      newMessageHandlerFunc
		editMessageHandler     editMessageHandlerFunc
		deleteMessageHandler   deleteMessageHandlerFunc
		pinMessageHandler      pinMessageHandlerFunc
		unpinMessageHandler    unpinMessageHandlerFunc
		newChatMemberHandler   newChatMembersHandlerFunc
		leftChatMembersHandler leftChatMembersHandlerFunc
		changeChatInfoHandler  changeChatInfoHandlerFunc
	}
}

// New creates new instance of Bot
func New(token string, client *http.Client) *Bot {
	return &Bot{
		state:        botStateStopped,
		token:        token,
		apiBaseURL:   apiBaseURL,
		client:       client,
		pollDuration: time.Minute,
	}
}

//easyjson:json
// GetSelfResponse represents info about self.
type GetSelfResponse struct {
	UserID    string `json:"userId"`
	Nick      string `json:"nick"`
	FirstName string `json:"firstName"`
	About     string `json:"about"`
	Photo     []struct {
		URL string `json:"url"`
	}
	Ok bool
}

// GetSelf returns information about Bot.
// It can be used to validate API token.
func (b *Bot) GetSelf(ctx context.Context) (*GetSelfResponse, error) {
	r, err := http.NewRequest(http.MethodGet, b.apiBaseURL+"/self/get", nil)
	if err != nil {
		return nil, err
	}

	q := r.URL.Query()
	r.URL.RawQuery = q.Encode()

	httpResp, err := b.doRequest(ctx, r)
	if err != nil {
		return nil, err
	}

	resp := &GetSelfResponse{}
	err = easyjson.UnmarshalFromReader(httpResp.Body, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (b *Bot) doRequest(ctx context.Context, r *http.Request) (*http.Response, error) {
	r = r.WithContext(ctx)
	q := r.URL.Query()
	q.Add(tokenQueryParam, b.token)
	r.URL.RawQuery = q.Encode()

	return b.client.Do(r)
}

// SetNewMessageHandler sets the handler to events about new message.
func (b *Bot) SetNewMessageHandler(fn newMessageHandlerFunc) {
	b.handlers.newMessageHandler = fn
}

// SetEditMessageHandler sets the handler to events about edit message.
func (b *Bot) SetEditMessageHandler(fn editMessageHandlerFunc) {
	b.handlers.editMessageHandler = fn
}

// sets the handler to events about delete message.
func (b *Bot) SetDeleteMessageHandler(fn deleteMessageHandlerFunc) {
	b.handlers.deleteMessageHandler = fn
}

// sets the handler to events about pin message.
func (b *Bot) SetPinMessageHandler(fn pinMessageHandlerFunc) {
	b.handlers.pinMessageHandler = fn
}

// SetUnpinMessageHandler sets the handler to events about unpin message.
func (b *Bot) SetUnpinMessageHandler(fn unpinMessageHandlerFunc) {
	b.handlers.unpinMessageHandler = fn
}

// SetNewChatMemberHandler sets the handler to events about new chat member.
func (b *Bot) SetNewChatMemberHandler(fn newChatMembersHandlerFunc) {
	b.handlers.newChatMemberHandler = fn
}

// SetLeftChatMemberHandler sets the handler to events about left chat member.
func (b *Bot) SetLeftChatMemberHandler(fn leftChatMembersHandlerFunc) {
	b.handlers.leftChatMembersHandler = fn
}

// SetChangeChatInfoHandler sets the handler to events about change chat properties.
func (b *Bot) SetChangeChatInfoHandler(fn changeChatInfoHandlerFunc) {
	b.handlers.changeChatInfoHandler = fn
}
