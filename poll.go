package icqbotapi

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/mailru/easyjson"

	"icqbotapi/event"
)

//easyjson:json
type PollResponse struct {
	Events []event.Event `json:"events"`
}

func (b *Bot) PollEvents(ctx context.Context) <-chan event.Event {
	if b.state != botStateStopped {
		return nil
	}

	b.state = botStatePolling
	events := make(chan event.Event)

	go b.poll(ctx, events)

	return events
}

func (b *Bot) HandleEvents(ctx context.Context) {
	if b.state != botStateStopped {
		return
	}

	b.state = botStateHandling
	events := make(chan event.Event)

	unmarshal := func(r json.RawMessage, u easyjson.Unmarshaler) error {
		err := easyjson.Unmarshal(r, u)
		if err != nil {
			log.Print(err)
		}

		return err
	}

	go b.poll(ctx, events)

	for ev := range events {
		switch ev.Type {
		case event.EventTypeNewMessage:
			payload := event.EventNewMessagePayload{}
			err := unmarshal(ev.Payload, &payload)
			if err != nil {
				continue
			}

			b.handlers.newMessageHandler(payload)
		case event.EventTypeEditedMessage:
			payload := event.EventTypeEditedPayload{}
			err := unmarshal(ev.Payload, &payload)
			if err != nil {
				continue
			}

			b.handlers.editMessageHandler(payload)
		case event.EventTypeDeletedMessage:
			payload := event.EventTypeDeletedPayload{}
			err := unmarshal(ev.Payload, &payload)
			if err != nil {
				continue
			}

			b.handlers.deleteMessageHandler(payload)
		case event.EventTypePinnedMessage:
			payload := event.EventTypePinnedPayload{}
			err := unmarshal(ev.Payload, &payload)
			if err != nil {
				continue
			}

			b.handlers.pinMessageHandler(payload)
		case event.EventTypeUnpinnedMessage:
			payload := event.EventTypeUnpinnedPayload{}
			err := unmarshal(ev.Payload, &payload)
			if err != nil {
				continue
			}

			b.handlers.unpinMessageHandler(payload)
		case event.EventTypeNewChatMembers:
			payload := event.EventTypeNewChatMembersPayload{}
			err := unmarshal(ev.Payload, &payload)
			if err != nil {
				continue
			}

			b.handlers.newChatMemberHandler(payload)
		case event.EventTypeLeftChatMembers:
			payload := event.EventTypeLeftChatMembersPayload{}
			err := unmarshal(ev.Payload, &payload)
			if err != nil {
				continue
			}

			b.handlers.leftChatMembersHandler(payload)
		default:
			log.Panicf("unexpected event type: %v", ev.Type)
		}
	}
}

func (b *Bot) poll(ctx context.Context, events chan<- event.Event) {
	lastEventID := 0

	for {
		select {
		case <-ctx.Done():
			close(events)
			return
		default:
		}

		req, err := http.NewRequest(http.MethodGet, b.apiBaseURL+"/events/get", nil)
		if err != nil {
			log.Printf("polling timeout exceeded: %v", err)
			continue
		}

		q := req.URL.Query()
		q.Set("lastEventId", strconv.Itoa(lastEventID))
		q.Set("pollTime", strconv.FormatInt(60, 10))
		req.URL.RawQuery = q.Encode()

		httpResp, err := b.doRequest(ctx, req)
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

			events <- ev
		}

		lastEventID = maxEventId
	}
}
