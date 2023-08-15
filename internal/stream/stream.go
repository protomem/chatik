package stream

import (
	"bytes"
	"context"
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/protomem/chatik/pkg/logging"
)

type Subscriber func() <-chan *Message

type Stream struct {
	logger logging.Logger
	pubsub *gochannel.GoChannel
}

func New(logger logging.Logger) *Stream {
	return &Stream{
		logger: logger.With("component", "stream"),
		pubsub: gochannel.NewGoChannel(
			gochannel.Config{
				OutputChannelBuffer: 0,
			},
			watermill.NopLogger{},
		),
	}
}

func (s *Stream) Subscribe(ctx context.Context) (Subscriber, error) {
	msgs, err := s.pubsub.Subscribe(ctx, "*")
	if err != nil {
		return nil, fmt.Errorf("stream.Subscribe: %w", err)
	}

	return func() <-chan *Message {
		res := make(chan *Message)

		go func() {
			for msg := range msgs {
				res <- NewMessage(msg.Payload)
				msg.Ack()
			}

			close(res)
		}()

		return res
	}, nil
}

func (s *Stream) Publish(msg *Message) error {
	err := s.pubsub.Publish("*", msg.m)
	if err != nil {
		return fmt.Errorf("stream.Publish: %w", err)
	}

	return nil
}

func (s *Stream) Close() error {
	err := s.pubsub.Close()
	if err != nil {
		return fmt.Errorf("stream.Close: %w", err)
	}

	return nil
}

type Message struct {
	m *message.Message
}

func NewMessage(payload []byte) *Message {
	return &Message{
		m: message.NewMessage(watermill.NewUUID(), bytes.Clone(payload)),
	}
}

func (m *Message) Payload() []byte {
	return m.m.Payload
}
