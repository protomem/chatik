package stream

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/protomem/chatik/internal/domain/model"
)

const (
	NewMessageEvn    EventType = "newMessage"
	RemoveMessageEvn EventType = "removeMessage"

	NewChannelEvn    EventType = "newChannel"
	RemoveChannelEvn EventType = "removeChannel"
)

type EventType string

type NewMessagePayload struct {
	Message model.Message `json:"message"`
}

type RemoveMessagePayload struct {
	MessageID uuid.UUID `json:"messageId"`
}

type NewChannelPayload struct {
	Channel model.Channel `json:"channel"`
}

type RemoveChannelPayload struct {
	ChannelID uuid.UUID `json:"channelId"`
}

type PayloadOrder interface {
	NewMessagePayload | RemoveMessagePayload | NewChannelPayload | RemoveChannelPayload
}

type Event[P PayloadOrder] struct {
	Type    EventType `json:"type"`
	Payload P         `json:"payload"`
}

func NewEvent[P PayloadOrder](eType EventType, payload P) Event[P] {
	return Event[P]{
		Type:    eType,
		Payload: payload,
	}
}

func NewMessageEvent(message model.Message) Event[NewMessagePayload] {
	return NewEvent[NewMessagePayload](
		NewMessageEvn,
		NewMessagePayload{
			Message: message,
		},
	)
}

func RemoveMessageEvent(messageID uuid.UUID) Event[RemoveMessagePayload] {
	return NewEvent[RemoveMessagePayload](
		RemoveMessageEvn,
		RemoveMessagePayload{
			MessageID: messageID,
		},
	)
}

func NewChannelEvent(channel model.Channel) Event[NewChannelPayload] {
	return NewEvent[NewChannelPayload](
		NewChannelEvn,
		NewChannelPayload{
			Channel: channel,
		},
	)
}

func RemoveChannelEvent(channelID uuid.UUID) Event[RemoveChannelPayload] {
	return NewEvent[RemoveChannelPayload](
		RemoveChannelEvn,
		RemoveChannelPayload{
			ChannelID: channelID,
		},
	)
}

func SendEvent[P PayloadOrder](s *Stream, e Event[P]) error {
	var (
		err error

		op = "stream.SendEvent"
	)

	msg, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("%s: marshal: %w", op, err)
	}

	err = s.Publish(NewMessage(msg))
	if err != nil {
		return fmt.Errorf("%s: publish: %w", op, err)
	}

	return nil
}
