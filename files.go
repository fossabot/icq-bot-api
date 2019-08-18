package icqbotapi

import (
	"context"
	"net/http"
	"net/url"

	"github.com/mailru/easyjson"
)

// FileID represents file identifier.
type FileID string

func (r FileID) validate() error {
	if r == "" {
		return errValidation
	}

	return nil
}

func (r FileID) contributeToQuery(q url.Values) {
	q.Set("fileId", string(r))
}

//easyjson:json
// FileInfoResponse represents info about uploaded file.
type FileInfoResponse struct {
	StatusResponse
	Type     string `json:"type"`
	Size     uint64 `json:"size"`
	Filename string `json:"filename"`
	URL      string `json:"url"`
}

//nolint:dupl
// GetFileInfo provides the file information function.
func (b *Bot) GetFileInfo(ctx context.Context, r FileID) (*FileInfoResponse, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, b.apiBaseURL+"/files/getInfo", nil)
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

	resp := &FileInfoResponse{}
	err = easyjson.UnmarshalFromReader(httpResp.Body, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
