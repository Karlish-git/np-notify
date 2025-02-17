package models

import (
	"np-notify/pkg/utils"
	"time"
)

type MessageCommentWrapper struct {
	MessageKey string           `json:"message_key"`
	Messages   []MessageComment `json:"messages"`
}

type MessageComment struct {
	Key     string `json:"key"`
	Payload struct {
		SenderUid int    `json:"senderUid"`
		Body      string `json:"body"`
	} `json:"payload"`
	Created time.Time `json:"created"`
}

func (mcw *MessageCommentWrapper) Parse(c []interface{}) {
	mcw.MessageKey = c[0].(string)

	msgs, err := utils.Unmarshalinterface[[]MessageComment](c[1])
	if err != nil {
		utils.LogError("Failed to parse message comments: %v", err)
	}
	mcw.Messages = msgs
}
