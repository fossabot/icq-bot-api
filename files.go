package icqbotapi

import (
	"net/http"
	"net/url"

	"github.com/mailru/easyjson"
)

type FileID string

func (r FileID) validate() error {
	if r == "" {
		return validationErr
	}

	return nil
}

func (r FileID) setQuery(q url.Values) {
	q.Set("fileId", string(r))
}

//easyjson:json
type FileInfoResponse struct {
	StatusResponse
	Type     string `json:"type"`
	Size     string `json:"size"`
	Filename string `json:"filename"`
	URL      string `json:"url"`
}

func (b *Bot) GetFileInfo(r FileID) (*FileInfoResponse, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, b.apiBaseURL+"/files/getInfo", nil)
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

	resp := &FileInfoResponse{}
	err = easyjson.UnmarshalFromReader(httpResp.Body, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
