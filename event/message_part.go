package event

import (
	"encoding/json"
)

type MessagePartType string

const (
	PartTypeSticker MessagePartType = "sticker"
	PartTypeMention MessagePartType = "mention"
	PartTypeVoice   MessagePartType = "voice"
	PartTypeFile    MessagePartType = "file"
	PartTypeForward MessagePartType = "forward"
	PartTypeReply   MessagePartType = "reply"
)

//easyjson:json
type MessagePart struct {
	Type    MessagePartType `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

//easyjson:json
type MessagePartSticker struct {
	FileID string `json:"fileId"`
}

//easyjson:json
type MessagePartMention struct {
	User
}

//easyjson:json
type MessagePartVoice struct {
	FileID string `json:"fileId"`
}

type FileType string

const (
	FileTypeImage FileType = "image"
	FileTypeAudio FileType = "audio"
	FileTypeVideo FileType = "video"
)

//easyjson:json
type MessagePartMessage struct {
	From      User   `json:"from"`
	MessageID string `json:"msgId"`
	Text      string `json:"text"`
	Timestamp uint64 `json:"timestamp"`
}

//easyjson:json
type MessagePartFile struct {
	File
}

//easyjson:json
type MessagePartForward struct {
	Message Message `json:"message"`
}

//easyjson:json
type MessagePartReply struct {
	Message Message `json:"message"`
}
