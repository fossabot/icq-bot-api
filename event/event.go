package event

import (
	"encoding/json"
)

type EventType string

const (
	EventTypeNewMessage      EventType = "newMessage"
	EventTypeEditedMessage   EventType = "editedMessage"
	EventTypeDeletedMessage  EventType = "deletedMessage"
	EventTypePinnedMessage   EventType = "pinnedMessage"
	EventTypeUnpinnedMessage EventType = "unpinnedMessage"
	EventTypeNewChatMembers  EventType = "newChatMembers"
	EventTypeLeftChatMembers EventType = "leftChatMembers"
)

//easyjson:json
type Event struct {
	EventID int             `json:"eventId"`
	Type    EventType       `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type ChatType string

const (
	ChatTypeChannel = "channel"
)

//easyjson:json
type Chat struct {
	ChatID string   `json:"chatId"`
	Type   ChatType `json:"type"`
	Title  string   `json:"title"`
}

//easyjson:json
type User struct {
	UserID    string `json:"userId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

//easyjson:json
type EventNewMessagePayload struct {
	MessageID string        `json:"messageId"`
	Chat      Chat          `json:"chat"`
	From      User          `json:"from"`
	Timestamp uint64        `json:"timestamp"`
	Text      string        `json:"text"`
	Parts     []MessagePart `json:"parts"`
}

//easyjson:json
type EventTypeEditedPayload struct {
	MessageID string `json:"messageId"`
	Chat      Chat   `json:"chat"`
	From      User   `json:"from"`
	Timestamp uint64 `json:"timestamp"`
	Text      string `json:"text"`
	EditedAt  uint64 `json:"editedTimestamp"`
}

//easyjson:json
type EventTypeDeletedPayload struct {
	MessageID string `json:"messageId"`
	Chat      Chat   `json:"chat"`
	Timestamp uint64 `json:"timestamp"`
}

//easyjson:json
type EventTypePinnedPayload struct {
	MessageID string `json:"messageId"`
	Chat      Chat   `json:"chat"`
	From      User   `json:"from"`
	Timestamp uint64 `json:"timestamp"`
	Text      string `json:"text"`
}

//easyjson:json
type EventTypeUnpinnedPayload struct {
	MessageID string `json:"messageId"`
	Chat      Chat   `json:"chat"`
	Timestamp uint64 `json:"timestamp"`
}

//easyjson:json
type EventTypeNewChatMembersPayload struct {
	MessageID string `json:"messageId"`
	Chat      Chat   `json:"chat"`
	Timestamp uint64 `json:"timestamp"`
}

//easyjson:json
type EventTypeLeftChatMembersPayload struct {
	Chat        Chat   `json:"chat"`
	LeftMembers []User `json:"leftMembers"`
	RemovedBy   User   `json:"removedBy"`
}

type Message struct {
	MessageID string `json:"messageId"`
	From      User   `json:"from"`
	Text      string `json:"text"`
	Timestamp uint64 `json:"timestamp"`
}

//easyjson:json
type EventTypeReplyPayload struct {
	Message Message `json:"message"`
}

//easyjson:json
type EventTypeForwardPayload struct {
	Message Message `json:"message"`
}

//easyjson:json
type File struct {
	FileID  string   `json:"fileId"`
	Type    FileType `json:"type"`
	Caption string   `json:"caption"`
}

//easyjson:json
type EventTypeFilePayload struct {
	File
}

//easyjson:json
type EventTypeVoicePayload struct {
	FileID string `json:"fileId"`
}

//easyjson:json
type EventTypeStickerPayload struct {
	FileID string `json:"fileId"`
}

//easyjson:json
type EventTypeMentionPayload struct {
	User
}
