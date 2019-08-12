package icqbotapi

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/mailru/easyjson"

	"icqbotapi/event"
)

//easyjson:json
type PollResponse struct {
	Events []event.Event `json:"events"`
}

func (b *Bot) PollEvents() (<-chan event.Event, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	lastEventID := 0

	events := make(chan event.Event)

	go func() {
		for {
			log.Print("req ", strconv.FormatInt(60, 10))
			req, err := http.NewRequest(http.MethodGet, b.apiBaseURL+"/events/get", nil)
			if err != nil {
				log.Printf("polling timeout exceeded: %v", err)
				continue
			}

			q := req.URL.Query()
			b.setToken(q)
			q.Set("lastEventId", strconv.Itoa(lastEventID))
			q.Set("pollTime", strconv.FormatInt(60, 10))
			req.URL.RawQuery = q.Encode()

			httpResp, err := b.client.Do(req)
			if err != nil {
				log.Printf("poll request error exceeded: %v", err)
				continue
			}

			resp := &PollResponse{
				make([]event.Event, 0),
			}

			err = easyjson.UnmarshalFromReader(httpResp.Body, resp)
			httpResp.Body.Close()
			if err != nil {
				panic(err)
			}

			maxEventId := 0
			for _, ev := range resp.Events {
				if ev.EventID > maxEventId {
					maxEventId = ev.EventID
				}

				switch ev.Type {
				case event.EventTypeNewMessage:
					payload := &event.EventNewMessagePayload{}
					err = easyjson.UnmarshalFromReader(strings.NewReader(string(ev.Payload)), payload)
					if err != nil {
						panic(err)
					}

					spew.Dump(payload)

				case event.EventTypeEditedMessage:
				case event.EventTypeDeletedMessage:
				case event.EventTypePinnedMessage:
				case event.EventTypeUnpinnedMessage:
				case event.EventTypeNewChatMembers:
				case event.EventTypeLeftChatMembers:
				case event.EventTypeChangedChatInfo:
				default:
					log.Panicf("unexpected event type: %v", ev.Type)
				}

			}

			lastEventID = maxEventId
		}
	}()

	_ = ctx

	return events, cancel
}
