package modules

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/nokrPOC/internal/config"
	"io"
	"log"
)

type Handler interface {
	PullMsgs(w io.Writer, projectID, subID string) error
}

type PubSubHandler struct {
	config *config.HermesServiceConfig
}

func NewPubSubService(config *config.HermesServiceConfig) *PubSubHandler{
	return &PubSubHandler{config}
}

func (p *PubSubHandler) PullMsgs() error {
	projectId := p.config.ProjectID
	subID := p.config.SubID
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectId)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}
	defer client.Close()

	sub := client.Subscription(subID)

	// Create a channel to handle messages to as they come in.
	cm := make(chan *pubsub.Message)
	defer close(cm)

	// Handle individual messages in a goroutine.
	go func() {
		for msg := range cm {
			log.Printf("msg: %+v", string(msg.Data))
			// This is where we need to unpack and get the data to craft the email
			msg.Ack()
		}
	}()

	// Receive blocks until the context is cancelled or an error occurs.
	err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		cm <- msg
	})
	if err != nil {
		return fmt.Errorf("Receive: %v", err)
	}

	return nil
}