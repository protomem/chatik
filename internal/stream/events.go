package stream

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/protomem/chatik/internal/agregate"
)

const (
	NewMessage    EventType = "newMessage"
	RemoveMessage EventType = "removeMessage"

	NewChannel    EventType = "newChannel"
	RemoveChannel EventType = "removeChannel"
)

type EventType string

type NewMessagePayload struct {
	Message agregate.Message `json:"message"`
}

type RemoveMessagePayload struct {
	MessageID uuid.UUID `json:"messageId"`
}

type NewChannelPayload struct {
	Channel agregate.Channel `json:"channel"`
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

func NewMessageEvent(message agregate.Message) Event[NewMessagePayload] {
	return NewEvent[NewMessagePayload](
		NewMessage,
		NewMessagePayload{
			Message: message,
		},
	)
}

func RemoveMessageEvent(messageID uuid.UUID) Event[RemoveMessagePayload] {
	return NewEvent[RemoveMessagePayload](
		RemoveMessage,
		RemoveMessagePayload{
			MessageID: messageID,
		},
	)
}

func NewChannelEvent(channel agregate.Channel) Event[NewChannelPayload] {
	return NewEvent[NewChannelPayload](
		NewChannel,
		NewChannelPayload{
			Channel: channel,
		},
	)
}

func RemoveChannelEvent(channelID uuid.UUID) Event[RemoveChannelPayload] {
	return NewEvent[RemoveChannelPayload](
		RemoveChannel,
		RemoveChannelPayload{
			ChannelID: channelID,
		},
	)
}

func SendEvent[P PayloadOrder](s *Broadcast, e Event[P]) error {
	msg, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("stream.SendEvent: marshal: %w", err)
	}

	s.SendMessage(msg)

	return nil
}
