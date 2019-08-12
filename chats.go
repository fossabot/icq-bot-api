package icqbotapi

import (
	"net/http"
	"net/url"

	"github.com/mailru/easyjson"
)

//easyjson:json
type Admin struct {
	UserID    string
	IsCreator bool
}

type ChatID string

func (r ChatID) validate() error {
	if r == "" {
		return validationErr
	}

	return nil
}

func (c ChatID) setQuery(q url.Values) {
	q.Set("chatId", string(c))
}

//easyjson:json
type GetAdminsResponse struct {
	StatusResponse
	Admins []Admin
}

func (b *Bot) GetChatAdmins(r ChatID) (*GetAdminsResponse, error) {
	req, err := http.NewRequest(http.MethodGet, b.apiBaseURL+"/chats/getAdmins", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	b.setToken(q)
	r.setQuery(q)
	req.URL.RawQuery = q.Encode()

	httpResp, err := b.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer httpResp.Body.Close()

	resp := &GetAdminsResponse{
		Admins: make([]Admin, 0),
	}

	err = easyjson.UnmarshalFromReader(httpResp.Body, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

//easyjson:json
type ChatInfoResponse struct {
	StatusResponse
	InviteLink url.URL `json:"inviteLink"`
	IsPublic   bool    `json:"public"`
	Title      string  `json:"title"`
	Group      string  `json:"group"`
}

func (b *Bot) GetChatInfo(r ChatID) (*ChatInfoResponse, error) {
	req, err := http.NewRequest(http.MethodGet, b.apiBaseURL+"/chats/getAdmins", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	b.setToken(q)
	r.setQuery(q)
	req.URL.RawQuery = q.Encode()

	httpResp, err := b.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer httpResp.Body.Close()

	resp := &ChatInfoResponse{}
	err = easyjson.UnmarshalFromReader(httpResp.Body, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

type ChatAction string

const (
	ChatActionTyping  ChatAction = "typing"
	ChatActionLooking ChatAction = "looking"
)

//easyjson:json
type ChatActionsRequest struct {
	ChatID  ChatID
	Actions []ChatAction
}

func (r *ChatActionsRequest) validate() error {
	return r.ChatID.validate()
}

func (r *ChatActionsRequest) setQuery(q url.Values) {
	r.ChatID.setQuery(q)

	for _, action := range r.Actions {
		q.Add("actions", string(action))
	}
}

func (b *Bot) SendChatActions(r ChatActionsRequest) (*StatusResponse, error) {
	req, err := http.NewRequest(http.MethodGet, b.apiBaseURL+"/chats/sendActions", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	b.setToken(q)
	r.setQuery(q)
	req.URL.RawQuery = q.Encode()

	httpResp, err := b.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer httpResp.Body.Close()

	resp := &StatusResponse{}
	err = easyjson.UnmarshalFromReader(httpResp.Body, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
