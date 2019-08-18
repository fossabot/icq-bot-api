package event

import (
	"encoding/json"
)

// Kind represents event kind
type Kind string

const (
	KindNewMessage      Kind = "newMessage"
	KindEditedMessage   Kind = "editedMessage"
	KindDeletedMessage  Kind = "deletedMessage"
	KindPinnedMessage   Kind = "pinnedMessage"
	KindUnpinnedMessage Kind = "unpinnedMessage"
	KindNewChatMember   Kind = "newChatMembers"
	KindLeftChatMembers Kind = "leftChatMembers"
)

//easyjson:json
type Event struct {
	EventID int             `json:"eventId"`
	Type    Kind            `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type ChatKind string

const (
	ChatKindChannel ChatKind = "channel"
	ChatKindGroup   ChatKind = "group"
)

//easyjson:json
type Chat struct {
	ChatID string   `json:"chatId"`
	Type   ChatKind `json:"type"`
	Title  string   `json:"title"`
}

//easyjson:json
type User struct {
	UserID    string `json:"userId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

//easyjson:json
type NewMessagePayload struct {
	MessageID string        `json:"msgId"`
	Chat      Chat          `json:"chat"`
	From      User          `json:"from"`
	Timestamp uint64        `json:"timestamp"`
	Text      string        `json:"text"`
	Parts     []MessagePart `json:"parts"`
}

//easyjson:json
type MessageEditPayload struct {
	MessageID string `json:"msgId"`
	Chat      Chat   `json:"chat"`
	From      User   `json:"from"`
	Timestamp uint64 `json:"timestamp"`
	Text      string `json:"text"`
	EditedAt  uint64 `json:"editedTimestamp"`
}

//easyjson:json
type MessageDeletePayload struct {
	MessageID string `json:"msgId"`
	Chat      Chat   `json:"chat"`
	Timestamp uint64 `json:"timestamp"`
}

//easyjson:json
type MessagePinPayload struct {
	MessageID string `json:"msgId"`
	Chat      Chat   `json:"chat"`
	From      User   `json:"from"`
	Timestamp uint64 `json:"timestamp"`
	Text      string `json:"text"`
}

//easyjson:json
type MessageUnpinPayload struct {
	MessageID string `json:"msgId"`
	Chat      Chat   `json:"chat"`
	Timestamp uint64 `json:"timestamp"`
}

//easyjson:json
type NewChatMembersPayload struct {
	MessageID string `json:"msgId"`
	Chat      Chat   `json:"chat"`
	Timestamp uint64 `json:"timestamp"`
}

//easyjson:json
type LeftChatMembersPayload struct {
	Chat        Chat   `json:"chat"`
	LeftMembers []User `json:"leftMembers"`
	RemovedBy   User   `json:"removedBy"`
}

type Message struct {
	MessageID string `json:"msgId"`
	From      User   `json:"from"`
	Text      string `json:"text"`
	Timestamp uint64 `json:"timestamp"`
}

//easyjson:json
type ReplyPayload struct {
	Message Message `json:"message"`
}

//easyjson:json
type ForwardPayload struct {
	Message Message `json:"message"`
}

//easyjson:json
type File struct {
	FileID  string   `json:"fileId"`
	Type    FileType `json:"type"`
	Caption string   `json:"caption"`
}

//easyjson:json
type FilePayload struct {
	File
}

//easyjson:json
type VoicePayload struct {
	FileID string `json:"fileId"`
}

//easyjson:json
type StickerPayload struct {
	FileID string `json:"fileId"`
}

//easyjson:json
type MentionPayload struct {
	User
}
