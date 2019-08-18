package icqbotapi

import (
	"bytes"
	"context"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"

	"github.com/mailru/easyjson/opt"

	"github.com/mailru/easyjson"
)

var errValidation = errors.New("validation error")

// SendSendTextRequest represents plain text interaction request.
type SendTextRequest struct {
	ChatID           string
	Text             string
	ReplyMessageID   uint64
	ForwardChatID    string
	ForwardMessageID uint64
}

func (r *SendTextRequest) validate() error {
	if r.ChatID == "" {
		return errValidation
	}

	// id цитируемого сообщения не может быть передано одновременно с forwardChatId и forwardMsgId.
	if r.ReplyMessageID != 0 &&
		(r.ForwardChatID != "" || r.ForwardMessageID != 0) {
		return errValidation
	}

	// id чата, из которого будет переслано сообщение передается только с forwardMsgId,
	// не может быть передано с replyMsgId.
	if r.ForwardChatID != "" &&
		(r.ForwardMessageID == 0 || r.ReplyMessageID != 0) {
		return errValidation
	}

	// id пересылаемого сообщения передается только с forwardChatId,
	// не может быть передано с replyMsgId.
	if r.ForwardMessageID != 0 &&
		(r.ForwardChatID == "" || r.ReplyMessageID != 0) {
		return errValidation
	}

	return nil
}

func (r *SendTextRequest) contributeToQuery(q url.Values) {
	q.Set("chatId", r.ChatID)
	q.Set("text", r.Text)

	if r.ReplyMessageID != 0 {
		q.Set("replyMsgId", strconv.FormatUint(r.ReplyMessageID, 10))
	}

	if r.ForwardChatID != "" {
		q.Set("forwardChatId", r.ForwardChatID)
	}

	if r.ForwardMessageID != 0 {
		q.Set("forwardMsgId", strconv.FormatUint(r.ForwardMessageID, 10))
	}
}

//easyjson:json
// StatusResponse represents common response status data.
type StatusResponse struct {
	Ok          bool       `json:"ok"`
	Description opt.String `json:"description"`
}

//easyjson:json
// StatusMessageIDResponse represents response status data for requests which deal with messages.
type StatusMessageIDResponse struct {
	StatusResponse
	MessageID string `json:"msgId"`
}

//nolint:dupl
// SendText performs plain text message function.
func (b *Bot) SendText(ctx context.Context, r *SendTextRequest) (*StatusMessageIDResponse, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, b.apiBaseURL+"/messages/sendText", nil)
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

	resp := &StatusMessageIDResponse{}
	err = easyjson.UnmarshalFromReader(httpResp.Body, resp)
	if err != nil {
		return nil, err
	}

	return resp, err
}

type fileRequest struct {
	SendTextRequest
	Caption string
	IsVoice bool
}

func (r *fileRequest) contributeToQuery(q url.Values) {
	r.SendTextRequest.contributeToQuery(q)
	q.Set("caption", r.Caption)
}

// SendFileRequest presents request with plain text with attached existing file/voice.
type SendFileRequest struct {
	fileRequest
	FileID string
}

func (r *SendFileRequest) validate() error {
	if err := r.fileRequest.validate(); err != nil {
		return err
	}

	if r.FileID == "" {
		return errValidation
	}

	return nil
}

func (r *SendFileRequest) contributeToQuery(q url.Values) {
	r.fileRequest.contributeToQuery(q)
	q.Set("fileId", r.FileID)
}

// SendNewFileRequest presents request with plain text with attached new file/voice.
type SendNewFileRequest struct {
	fileRequest
	File     io.Reader
	Filename string
}

// SendFile provides the function of sending text messages with already downloaded file attachments.
func (b *Bot) SendFile(ctx context.Context, r *SendFileRequest) (*StatusMessageIDResponse, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}

	m := "/messages/sendFile"
	if r.IsVoice {
		m = "/messages/sendVoice"
	}

	req, err := http.NewRequest(http.MethodGet, b.apiBaseURL+m, nil)
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

	resp := &StatusMessageIDResponse{}
	err = easyjson.UnmarshalFromReader(httpResp.Body, resp)
	if err != nil {
		return nil, err
	}

	return resp, err
}

//easyjson:json
// SendNewFileResponse contains information about sent message with attached files.
type SendNewFileResponse struct {
	StatusMessageIDResponse
	FileID string `json:"fileId"`
}

// SendNewFile provides the function of sending text messages with file attachments.
func (b *Bot) SendNewFile(ctx context.Context, r *SendNewFileRequest) (*SendNewFileResponse, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	mw := multipart.NewWriter(buf)

	part, err := mw.CreateFormFile("file", r.Filename)
	if err != nil {
		return nil, err
	}

	if _, err = io.Copy(part, r.File); err != nil {
		return nil, err
	}

	//nolint:govet
	if err := mw.Close(); err != nil {
		return nil, err
	}

	m := "/messages/sendFile"
	if r.IsVoice {
		m = "/messages/sendVoice"
	}

	req, err := http.NewRequest(http.MethodPost, b.apiBaseURL+m, buf)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	r.contributeToQuery(q)
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Content-Type", mw.FormDataContentType())

	httpResp, err := b.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	defer httpResp.Body.Close()

	resp := &SendNewFileResponse{}
	err = easyjson.UnmarshalFromReader(httpResp.Body, resp)
	if err != nil {
		return nil, err
	}

	return resp, err
}

//easyjson:json
// EditMessageRequest represents data for editing a messages.
type EditMessageRequest struct {
	ChatID    string `json:"chatId"`
	MessageID string `json:"msgId"`
	Text      string `json:"text"`
}

func (r *EditMessageRequest) validate() error {
	if r.ChatID == "" ||
		r.MessageID == "" ||
		r.Text == "" {
		return errValidation
	}

	return nil
}

func (r *EditMessageRequest) contributeToQuery(q url.Values) {
	q.Set("chatId", r.ChatID)
	q.Set("msgId", r.MessageID)
	q.Set("text", r.Text)
}

//nolint:dupl
// EditMessage provides the function of editing messages.
func (b *Bot) EditMessage(ctx context.Context, r *EditMessageRequest) (*StatusMessageIDResponse, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, b.apiBaseURL+"/messages/editText", nil)
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

	resp := &StatusMessageIDResponse{}
	err = easyjson.UnmarshalFromReader(httpResp.Body, resp)
	if err != nil {
		return nil, err
	}

	return resp, err
}

//easyjson:json
// DeleteMessageRequest represents data for deleting a messages.
type DeleteMessageRequest struct {
	ChatID    string `json:"chatId"`
	MessageID string `json:"msgId"`
}

func (r *DeleteMessageRequest) validate() error {
	if r.ChatID == "" ||
		r.MessageID == "" {
		return errValidation
	}

	return nil
}

func (r *DeleteMessageRequest) contributeToQuery(q url.Values) {
	q.Set("chatId", r.ChatID)
	q.Set("msgId", r.MessageID)
}

//nolint:dupl
// EditMessage provides the function of deleting messages.
func (b *Bot) DeleteMessage(ctx context.Context, r *DeleteMessageRequest) (*StatusResponse, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, b.apiBaseURL+"/messages/deleteMessages", nil)
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

	resp := &StatusResponse{}
	err = easyjson.UnmarshalFromReader(httpResp.Body, resp)
	if err != nil {
		return nil, err
	}

	return resp, err
}
