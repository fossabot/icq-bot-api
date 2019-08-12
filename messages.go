package icqbotapi

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"

	"github.com/mailru/easyjson"
)

var validationErr = errors.New("validation error")

type SendTextRequest struct {
	ChatID           string
	Text             string
	ReplyMessageID   uint64
	ForwardChatID    string
	ForwardMessageID uint64
}

func (r *SendTextRequest) validate() error {
	if r.ChatID == "" {
		return validationErr
	}

	// id цитируемого сообщения не может быть передано одновременно с forwardChatId и forwardMsgId
	if r.ReplyMessageID != 0 &&
		(r.ForwardChatID != "" || r.ForwardMessageID != 0) {
		return validationErr
	}

	// id чата, из которого будет переслано сообщение передается только с forwardMsgId.
	// не может быть передано с replyMsgId.
	if r.ForwardChatID != "" &&
		(r.ForwardMessageID == 0 || r.ReplyMessageID != 0) {
		return validationErr
	}

	// id пересылаемого сообщения передается только с forwardChatId
	// не может быть передано с replyMsgId
	if r.ForwardMessageID != 0 &&
		(r.ForwardChatID == "" || r.ReplyMessageID != 0) {
		return validationErr
	}

	return nil
}

func (r *SendTextRequest) setQuery(q url.Values) {
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
type StatusResponse struct {
	Ok          bool   `json:"ok"`
	Description string `json:"description"`
}

//easyjson:json
type MessageIDResponse struct {
	StatusResponse
	MsgID       string `json:"msgId"`
}

func (b *Bot) SendText(r *SendTextRequest) (*MessageIDResponse, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, b.apiBaseURL+"/messages/sendText", nil)
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

	resp := &MessageIDResponse{}
	err = easyjson.UnmarshalFromReader(httpResp.Body, resp)
	if err != nil {
		return nil, err
	}

	return resp, err
}

type FileRequest struct {
	SendTextRequest
	Caption string
}

func (r *FileRequest) setQuery(q url.Values) {
	r.SendTextRequest.setQuery(q)
	q.Set("caption", r.Caption)
}

type SendFileRequest struct {
	FileRequest
	FileID  string
	IsVoice bool
}

func (r *SendFileRequest) validate() error {
	if err := r.FileRequest.validate(); err != nil {
		return err
	}

	if r.FileID == "" {
		return validationErr
	}

	return nil
}

func (r *SendFileRequest) setQuery(q url.Values) {
	r.FileRequest.setQuery(q)
	q.Set("fileId", r.FileID)
}

type SendNewFileRequest struct {
	FileRequest
	File     io.Reader
	Filename string
	IsVoice  bool
}

func (b *Bot) SendFile(r *SendFileRequest) (*MessageIDResponse, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}

	method := "/messages/sendFile"
	if r.IsVoice {
		method = "/messages/sendVoice"
	}

	req, err := http.NewRequest(http.MethodGet, b.apiBaseURL+method, nil)
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

	resp := &MessageIDResponse{}
	err = easyjson.UnmarshalFromReader(httpResp.Body, resp)
	if err != nil {
		return nil, err
	}

	return resp, err
}

//easyjson:json
type UploadFileResponse struct {
	MessageIDResponse
	FileID string `json:"fileId"`
}

func (b *Bot) SendNewFile(r *SendNewFileRequest) (*UploadFileResponse, error) {
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

	if err := mw.Close(); err != nil {
		return nil, err
	}

	method := "/messages/sendFile"
	if r.IsVoice {
		method = "/messages/sendVoice"
	}

	req, err := http.NewRequest(http.MethodPost, b.apiBaseURL+method, buf)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	b.setToken(q)
	r.setQuery(q)
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Content-Type", mw.FormDataContentType())

	httpResp, err := b.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer httpResp.Body.Close()

	resp := &UploadFileResponse{}
	err = easyjson.UnmarshalFromReader(httpResp.Body, resp)
	if err != nil {
		return nil, err
	}

	return resp, err
}

//easyjson:json
type EditMessageRequest struct {
	ChatID    string `json:"chatId"`
	MessageID string `json:"msgId"`
	Text      string `json:"text"`
}

func (r *EditMessageRequest) validate() error {
	if r.ChatID == "" ||
		r.MessageID == "" ||
		r.Text == "" {
		return validationErr
	}

	return nil
}

func (r *EditMessageRequest) setQuery(q url.Values) {
	q.Set("chatId", r.ChatID)
	q.Set("msgId", r.MessageID)
	q.Set("text", r.Text)
}

func (b *Bot) EditMessage(r *EditMessageRequest) (*MessageIDResponse, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, b.apiBaseURL+"/messages/editText", nil)
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

	resp := &MessageIDResponse{}
	err = easyjson.UnmarshalFromReader(httpResp.Body, resp)
	if err != nil {
		return nil, err
	}

	return resp, err
}

//easyjson:json
type DeleteMessageRequest struct {
	ChatID string `json:"chatId"`
	MessageID string `json:"messageId"`
}

func (r *DeleteMessageRequest) validate() error {
	if r.ChatID == "" ||
		r.MessageID == "" {
		return validationErr
	}

	return nil
}

func (r *DeleteMessageRequest) setQuery(q url.Values) {
	q.Set("chatId", r.ChatID)
	q.Set("msgId", r.MessageID)
}

func (b *Bot) DeleteMessage(r *DeleteMessageRequest) (*StatusResponse, error)  {
	if err := r.validate(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, b.apiBaseURL + "/messages/deleteMessages", nil)
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

	return resp, err
}
