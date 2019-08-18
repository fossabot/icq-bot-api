package icqbotapi

import (
	"context"
	"net/http"
	"time"

	"github.com/mailru/easyjson"

	"icqbotapi/event"
)

const (
	icqAPIBaseURL   = "https://api.icq.net/bot/v1"
	agentAPIBaseURL = "https://agent.mail.ru/bot/v1"
	tokenQueryParam = "token"
)

type apiType int

const (
	APITypeICQ apiType = iota
	APITypeAgent
)

type newMessageHandlerFunc func(e event.NewMessagePayload)
type editMessageHandlerFunc func(e event.MessageEditPayload)
type deleteMessageHandlerFunc func(e event.MessageDeletePayload)
type pinMessageHandlerFunc func(e event.MessagePinPayload)
type unpinMessageHandlerFunc func(e event.MessageUnpinPayload)
type newChatMembersHandlerFunc func(e event.NewChatMembersPayload)
type leftChatMembersHandlerFunc func(e event.LeftChatMembersPayload)
type errorHandlerFunc func(err error)

type botState int

const (
	botStatePolling = iota
	botStateHandling
	botStateStopped
)

type botHandlers struct {
	newMessageHandler    newMessageHandlerFunc
	editMessageHandler   editMessageHandlerFunc
	deleteMessageHandler deleteMessageHandlerFunc

	pinMessageHandler   pinMessageHandlerFunc
	unpinMessageHandler unpinMessageHandlerFunc

	newChatMemberHandler   newChatMembersHandlerFunc
	leftChatMembersHandler leftChatMembersHandlerFunc

	errorHandler errorHandlerFunc
}

type Bot struct {
	state botState

	token        string
	apiBaseURL   string
	client       *http.Client
	pollDuration time.Duration
	handlers     botHandlers
}

// New creates new instance of Bot
func New(token string, client *http.Client, t apiType) *Bot {
	apiBaseURL := icqAPIBaseURL
	if t == APITypeAgent {
		apiBaseURL = agentAPIBaseURL
	}

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
	} `json:"photo"`
	Ok bool `json:"ok"`
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

// SetErrorHandler sets the processing errors handler.
func (b *Bot) SetErrorHandler(fn errorHandlerFunc) {
	b.handlers.errorHandler = fn
}

func (b *Bot) unmarshalEvent(p []byte, v easyjson.Unmarshaler) bool {
	if err := easyjson.Unmarshal(p, v); err != nil {
		b.handleError(err)
		return false
	}

	return true
}

func (b *Bot) handleNewMessage(r event.Event) {
	if b.handlers.newMessageHandler != nil {
		e := event.NewMessagePayload{}
		if !b.unmarshalEvent(r.Payload, &e) {
			return
		}

		b.handlers.newMessageHandler(e)
	}
}

func (b *Bot) handleEditMessage(r event.Event) {
	if b.handlers.editMessageHandler != nil {
		e := event.MessageEditPayload{}
		if !b.unmarshalEvent(r.Payload, &e) {
			return
		}

		b.handlers.editMessageHandler(e)
	}
}

func (b *Bot) handleDeleteMessage(r event.Event) {
	if b.handlers.deleteMessageHandler != nil {
		e := event.MessageEditPayload{}
		if !b.unmarshalEvent(r.Payload, &e) {
			return
		}

		b.handlers.editMessageHandler(e)
	}
}

func (b *Bot) handlePinMessage(r event.Event) {
	if b.handlers.pinMessageHandler != nil {
		e := event.MessagePinPayload{}
		if !b.unmarshalEvent(r.Payload, &e) {
			return
		}

		b.handlers.pinMessageHandler(e)
	}
}

func (b *Bot) handleUnpinMessage(r event.Event) {
	if b.handlers.unpinMessageHandler != nil {
		e := event.MessageUnpinPayload{}
		if !b.unmarshalEvent(r.Payload, &e) {
			return
		}

		b.handlers.unpinMessageHandler(e)
	}
}

func (b *Bot) handleNewChatMember(r event.Event) {
	if b.handlers.newChatMemberHandler != nil {
		e := event.NewChatMembersPayload{}
		if !b.unmarshalEvent(r.Payload, &e) {
			return
		}

		b.handlers.newChatMemberHandler(e)
	}
}

func (b *Bot) handleLeftChatMember(r event.Event) {
	if b.handlers.leftChatMembersHandler != nil {
		e := event.LeftChatMembersPayload{}
		if !b.unmarshalEvent(r.Payload, &e) {
			return
		}

		b.handlers.leftChatMembersHandler(e)
	}
}

func (b *Bot) handleError(err error) {
	if b.handlers.errorHandler != nil {
		b.handlers.errorHandler(err)
	}
}
