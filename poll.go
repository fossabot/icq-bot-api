package icqbotapi

import (
	"context"
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

	go b.poll(ctx, events)
	go func() {
		for ev := range events {
			switch ev.Type {
			case event.KindNewMessage:
				b.handleNewMessage(ev)
			case event.KindEditedMessage:
				b.handleEditMessage(ev)
			case event.KindDeletedMessage:
				b.handleDeleteMessage(ev)
			case event.KindPinnedMessage:
				b.handlePinMessage(ev)
			case event.KindUnpinnedMessage:
				b.handleUnpinMessage(ev)
			case event.KindNewChatMember:
				b.handleNewChatMember(ev)
			case event.KindLeftChatMembers:
				b.handleLeftChatMember(ev)
			default:
				log.Panicf("unexpected event type: %v", ev.Type)
			}
		}
	}()
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

		maxEventID := 0
		for _, ev := range resp.Events {
			if ev.EventID > maxEventID {
				maxEventID = ev.EventID
			}

			events <- ev
		}

		lastEventID = maxEventID
	}
}
