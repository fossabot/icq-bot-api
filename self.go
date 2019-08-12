package icqbotapi

import (
	"net/http"

	"github.com/mailru/easyjson"
)

//easyjson:json
type GetResponse struct {
	UserID    string `json:"userId"`
	Nick      string `json:"nick"`
	FirstName string `json:"firstName"`
	About     string `json:"about"`
	Photo     []struct {
		URL string `json:"url"`
	}
	Ok bool
}

// Get returns information about Bot.
// It can be used to validate API token
func (b *Bot) Get() (*GetResponse, error) {
	req, err := http.NewRequest(http.MethodGet, b.apiBaseURL+"/self/get", nil)
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	b.setToken(query)
	req.URL.RawQuery = query.Encode()

	httpResp, err := b.client.Do(req)
	if err != nil {
		return nil, err
	}

	resp := &GetResponse{}
	err = easyjson.UnmarshalFromReader(httpResp.Body, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
