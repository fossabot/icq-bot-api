package icqbotapi

import (
	"context"
	"net/http"
	"net/url"

	"github.com/mailru/easyjson"
)

//easyjson:json
// Admin represents chat administrator.
type Admin struct {
	UserID    string
	IsCreator bool
}

// ChatID represents chat identifier.
type ChatID string

func (c ChatID) validate() error {
	if c == "" {
		return errValidation
	}

	return nil
}

func (c ChatID) contributeToQuery(q url.Values) {
	q.Set("chatId", string(c))
}

//easyjson:json
// GetAdminsResponse represents information about chat administration.
type GetAdminsResponse struct {
	StatusResponse
	Admins []Admin
}

// GetChatAdmins provides a function to obtain information about the chat administration.
func (b *Bot) GetChatAdmins(ctx context.Context, r ChatID) (*GetAdminsResponse, error) {
	req, err := http.NewRequest(http.MethodGet, b.apiBaseURL+"/chats/getAdmins", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	r.contributeToQuery(q)
	req.URL.RawQuery = q.Encode()

	httpResp, err := b.doRequest(ctx, req)
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
// ChatInfoResponse represents chat properties.
type ChatInfoResponse struct {
	StatusResponse
	InviteLink url.URL `json:"inviteLink"`
	IsPublic   bool    `json:"public"`
	Title      string  `json:"title"`
	Group      string  `json:"group"`
}

//nolint:dupl
// GetChatInfo returns information about chat.
func (b *Bot) GetChatInfo(ctx context.Context, r ChatID) (*ChatInfoResponse, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, b.apiBaseURL+"/chats/getInfo", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	r.contributeToQuery(q)
	req.URL.RawQuery = q.Encode()

	httpResp, err := b.doRequest(ctx, req)
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

// ChatAction represents actions which can be in chats.
type ChatAction string

const (
	ChatActionTyping  ChatAction = "typing"
	ChatActionLooking ChatAction = "looking"
)

//easyjson:json
// ChatActionsRequest represents container to sent actions to chat.
type ChatActionsRequest struct {
	ChatID  ChatID
	Actions []ChatAction
}

func (r *ChatActionsRequest) validate() error {
	return r.ChatID.validate()
}

func (r *ChatActionsRequest) contributeToQuery(q url.Values) {
	r.ChatID.contributeToQuery(q)

	for _, action := range r.Actions {
		q.Add("actions", string(action))
	}
}

// SendChatActions provides a function to sent actions to chat.
func (b *Bot) SendChatActions(ctx context.Context, reqs <-chan ChatActionsRequest) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case r, ok := <-reqs:
				if !ok {
					return
				}

				if err := r.validate(); err != nil {
					b.handleError(err)
					continue
				}

				req, err := http.NewRequest(http.MethodGet, b.apiBaseURL+"/chats/sendActions", nil)
				if err != nil {
					b.handleError(err)
					continue
				}

				req = req.WithContext(ctx)

				q := req.URL.Query()
				r.contributeToQuery(q)
				req.URL.RawQuery = q.Encode()

				httpResp, err := b.doRequest(ctx, req)
				if err != nil {
					b.handleError(err)
					continue
				}

				httpResp.Body.Close()
			}
		}
	}()
}
